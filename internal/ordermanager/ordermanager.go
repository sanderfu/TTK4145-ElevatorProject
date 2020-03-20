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

var primaryQueue []datatypes.QueueOrder
var backupQueue []datatypes.QueueOrder

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

	go costRequestListener()
	go orderRegHW()
	go orderRegSW()
	go queueModifier()
	go orderCompleteWatch()
	go backupWatch()
	go orderRegisteredWatch()
}

////////////////////////////////////////////////////////////////////////////////
// Cost
////////////////////////////////////////////////////////////////////////////////
// TODO: Find bug here
func genCostAns(costReq datatypes.CostRequest) datatypes.CostAnswer {
	var costAns datatypes.CostAnswer
	costAns.DestinationID = costReq.SourceID
	if costReq.OrderType == datatypes.INSIDE && costReq.SourceID != costReq.DestinationID {
		costAns.CostValue = maxCost + 1
	} else {
		costAns.CostValue = 2*len(primaryQueue) + 1*len(backupQueue)
	}
	return costAns
}

func costRequestListener() {
	var costReq datatypes.CostRequest
	var costAns datatypes.CostAnswer
	for {
		costReq = <-channels.CostRequestTOM
		//Change next 2 lines with genCostAns function
		costAns.CostValue = 2*len(primaryQueue) + 1*len(backupQueue)
		costAns.DestinationID = costReq.SourceID
		channels.CostAnswerFOM <- costAns
	}
}

func orderRegHW() {
	for {
		order := <-channels.OrderFHM

		//Make cost request
		// TODO: Intialize struct nicer, see fsm or hwman or something
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

func removeFromQueue(order datatypes.QueueOrder, primary bool) {
	var queue []datatypes.QueueOrder
	switch primary {
	case true:
		queue = primaryQueue
	case false:
		queue = backupQueue
	}

	j := 0
	for _, element := range queue {
		if order.Floor != element.Floor {
			queue[j] = element
			j++
		}
	}

	queue = queue[:j]
	switch primary {
	case true:
		primaryQueue = queue
	case false:
		backupQueue = queue
	}
}

func orderInQueue(order datatypes.QueueOrder, primary bool) bool {
	switch primary {
	case true:
		for _, elem := range primaryQueue {
			if elem.Floor == order.Floor && elem.OrderType == order.OrderType {
				return true
			}
		}
	case false:
		for _, elem := range backupQueue {
			if elem.Floor == order.Floor && elem.OrderType == order.OrderType {
				return true
			}
		}
	}
	return false
}

func queueModifier() {
	for {
		select {
		case primaryOrder := <-primaryAppend:
			primary := true
			if !orderInQueue(primaryOrder, primary) {
				primaryQueue = append(primaryQueue, primaryOrder)
				logger.WriteLog(primaryQueue, primary, "/logs/")
				generateOrderRecvAck(primaryOrder)
			}
		case primaryOrder := <-primaryRemove:
			primary := true
			removeFromQueue(primaryOrder, primary)
			logger.WriteLog(primaryQueue, primary, "/logs/")
		case backupOrder := <-backupAppend:
			primary := false
			backupQueue = append(backupQueue, backupOrder)
			logger.WriteLog(backupQueue, primary, "/logs/")
			generateOrderRecvAck(backupOrder)
		case backupOrder := <-backupRemove:
			primary := false
			removeFromQueue(backupOrder, primary)
			logger.WriteLog(backupQueue, primary, "/logs/")
		}
	}
}

func orderCompleteWatch() {
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

func backupWatch() {
	for {
		for _, elem := range backupQueue {
			if time.Since(elem.RegistrationTime) > backupWaitS*time.Second {
				backupRemove <- elem
				primaryAppend <- elem
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func orderRegisteredWatch() {
	for {
		orderReg := <-channels.OrderRegisteredTOM

		channels.OrderRegisteredTHM <- orderReg
	}
}

//Worker functions
func GetFirstOrderInQueue() datatypes.QueueOrder {
	return primaryQueue[0]
}

func QueueEmpty() bool {
	if len(primaryQueue) == 0 {
		return true
	} else {
		return false
	}
}

func OrderToTakeAtFloor(floor int, ordertype int) bool {

	for _, order := range primaryQueue {
		if order.Floor == floor && (order.OrderType == ordertype || order.OrderType == datatypes.INSIDE) {
			return true
		}
	}
	return false
}
