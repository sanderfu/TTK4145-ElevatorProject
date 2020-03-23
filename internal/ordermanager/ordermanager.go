package ordermanager

import (
	"fmt"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/logger"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
)

var answerWaitMS time.Duration
var orderRecvAckWaitMS time.Duration
var maxCost int
var backupWaitS time.Duration

var start time.Time

var primaryAppend chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)
var primaryRemove chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)

var backupAppend chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)
var backupRemove chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)

//OrderManager ...
func OrderManager(resuming bool, lastPID string) {

	//Set global values based on configuration
	answerWaitMS = time.Duration(datatypes.Config.CostRequestTimeoutMS)
	orderRecvAckWaitMS = time.Duration(datatypes.Config.OrderReceiveAckTimeoutMS)
	maxCost = datatypes.Config.MaxCostValue
	backupWaitS = time.Duration(datatypes.Config.BackupTakeoverTimeoutS)

	start = time.Now()

	//If is resuming (after crash), load queues into memory
	if resuming {
		fmt.Println("Importing queue from crashed session")
		dir := "/" + lastPID + "/" + "logs"
		logger.ReadLogQueue(&primaryQueue, true, dir)
		logger.ReadLogQueue(&backupQueue, false, dir)
		var orderReg datatypes.OrderRegistered
		logger.WriteLog(primaryQueue, true, "/logs/")
		for i := 0; i < len(primaryQueue); i++ {
			orderReg.Floor = primaryQueue[i].Floor
			orderReg.OrderType = primaryQueue[i].OrderType
			channels.OrderRegisteredFOM <- orderReg
		}
		logger.WriteLog(backupQueue, false, "/logs/")
		fmt.Println("Resume successful")
	}

	go costReqListener()
	go orderRegHW()
	go orderRegSW()
	go queueModifier()
	go orderCompleteListener()
	go backupListener()
	go orderRegisteredListener()
}

func orderRegHW() {
	for {
		order := <-channels.OrderFHM

		//Make cost request
		var request datatypes.CostRequest
		request.OrderType = order.OrderType
		request.Floor = order.Floor

		//Broadcast cost request
		channels.CostRequestFOM <- request

		//Wait for answers
		done := time.After(answerWaitMS * time.Millisecond)
		primaryCost := maxCost + 1
		backupCost := maxCost + 1
	costWaitloop:
		for {
			select {
			case <-done:
				break costWaitloop
			case costAns := <-channels.CostAnswerTOM:
				if costAns.CostValue < primaryCost {
					backupCost = primaryCost
					primaryCost = costAns.CostValue
					order.BackupID = order.PrimaryID
					order.PrimaryID = costAns.SourceID
				} else if costAns.CostValue < backupCost {
					backupCost = costAns.CostValue
					order.BackupID = costAns.SourceID
				}
			}
		}
		//Handle situation with no backup
		if backupCost == maxCost+1 {
			order.BackupID = order.PrimaryID
		}
		channels.SWOrderFOM <- order
		//Wait for OrderRecAck from primary and backup
		done2 := time.After(orderRecvAckWaitMS * time.Millisecond)
		ackCounter := 0
	ackWaitloop:
		for {
			select {
			case <-done2:
				//Timer reached end, the order transmit is assumed to have failed and order is put back into the channel
				channels.OrderFHM <- order
				break ackWaitloop
			case orderRecvAck := <-channels.OrderRecvAckTOM:
				if orderRecvAck.SourceID == order.PrimaryID || orderRecvAck.SourceID == order.BackupID {
					//Check that ack matches order, if not throw it away as it has probably arrived to late for prev. order
					if orderRecvAck.Floor == order.Floor && orderRecvAck.OrderType == order.OrderType {
						ackCounter++
					}
				}
				if ackCounter == 2 {
					//Transmit was successful
					var orderReg = datatypes.OrderRegistered{
						Floor:     order.Floor,
						OrderType: order.OrderType,
					}
					channels.OrderRegisteredFOM <- orderReg
					break ackWaitloop
				}
			}
		}
	}
}

func generateOrderRecvAck(queueOrder datatypes.QueueOrder) {
	var orderRecvAck datatypes.OrderRecvAck
	orderRecvAck.OrderType = queueOrder.OrderType
	orderRecvAck.Floor = queueOrder.Floor
	orderRecvAck.DestinationID = queueOrder.SourceID
	channels.OrderRecvAckFOM <- orderRecvAck
}

func generateQueueOrder(order datatypes.Order) datatypes.QueueOrder {
	var queueOrder datatypes.QueueOrder
	queueOrder.SourceID = order.SourceID
	queueOrder.OrderType = order.OrderType
	queueOrder.Floor = order.Floor
	queueOrder.RegistrationTime = time.Now()
	return queueOrder
}

func orderRegSW() {
	for {
		select {
		case order := <-channels.SWOrderTOMPrimary:
			queueOrder := generateQueueOrder(order)
			primaryAppend <- queueOrder
		case order := <-channels.SWOrderTOMBackup:
			queueOrder := generateQueueOrder(order)
			backupAppend <- queueOrder
		}
	}
}

func orderCompleteListener() {
	for {
		select {
		case orderComplete := <-channels.OrderCompleteTOM:
			var queueOrder datatypes.QueueOrder
			queueOrder.OrderType = orderComplete.OrderType
			fmt.Println("Forwarding remove request to queueModifier")
			queueOrder.Floor = orderComplete.Floor
			primaryRemove <- queueOrder
			backupRemove <- queueOrder
			channels.OrderCompleteTHM <- orderComplete
			fmt.Println("The remove request has been handeled")
		case orderComplete := <-channels.OrderCompleteFFSM:
			channels.OrderCompleteFOM <- orderComplete
		}
	}
}

func orderRegisteredListener() {
	for {
		orderReg := <-channels.OrderRegisteredTOM

		channels.OrderRegisteredTHM <- orderReg
	}
}
