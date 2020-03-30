package ordermanager

import (
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

var costRequestTimeoutMS time.Duration
var orderRecvAckWaitMS time.Duration
var maxCostValue int
var backupTakeoverTimeoutS time.Duration

var start time.Time

//OrderManager ...
func OrderManager() {

	//Set global values based on configuration
	costRequestTimeoutMS = time.Duration(configuration.Config.CostRequestTimeoutMS)
	orderRecvAckWaitMS = time.Duration(configuration.Config.OrderReceiveAckTimeoutMS)
	maxCostValue = configuration.Config.MaxCostValue
	backupTakeoverTimeoutS = time.Duration(configuration.Config.BackupTakeoverTimeoutS)
	lastPID := configuration.Flags.LastPID

	start = time.Now()

	//If is resuming (after crash), load queues into memory
	restoreQueues(lastPID)

	go costRequestListener()
	go orderRegistrationHW()
	go orderRegistrationSW()
	go queueModifier()
	go orderCompleteListener()
	go backupTimeoutListener()
	go orderRegisteredListener()
}

func orderRegistrationHW() {
	for {
		order := <-channels.OrderFHM

		//Make cost request
		var request = datatypes.CostRequest{
			OrderType: order.OrderType,
			Floor:     order.Floor,
		}

		//Broadcast cost request
		channels.CostRequestFOM <- request

		//Wait for answers
		done := time.After(costRequestTimeoutMS * time.Millisecond)
		primaryCost := maxCostValue + 1
		backupCost := maxCostValue + 1
	costWaitloop:
		for {
			select {
			case <-done:
				break costWaitloop
			case costAns := <-channels.CostAnswerFNM:
				if costAns.CostValue < primaryCost {
					backupCost = primaryCost
					primaryCost = costAns.CostValue
					order.PrimaryID = costAns.SourceID
				} else if costAns.CostValue < backupCost {
					backupCost = costAns.CostValue
					order.BackupID = costAns.SourceID
				}
			}
		}
		//Handle situation with no backup
		if backupCost == maxCostValue+1 {
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
			case orderRecvAck := <-channels.OrderRecvAckFNM:
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
	var orderRecvAck = datatypes.OrderRecvAck{
		OrderType:     queueOrder.OrderType,
		Floor:         queueOrder.Floor,
		DestinationID: queueOrder.SourceID,
	}
	channels.OrderRecvAckFOM <- orderRecvAck
}

func generateQueueOrder(order datatypes.Order) datatypes.QueueOrder {
	var queueOrder = datatypes.QueueOrder{
		SourceID:         order.SourceID,
		OrderType:        order.OrderType,
		Floor:            order.Floor,
		RegistrationTime: time.Now(),
	}
	return queueOrder
}

func orderRegistrationSW() {
	for {
		select {
		case order := <-channels.SWOrderFNMPrimary:
			queueOrder := generateQueueOrder(order)
			channels.PrimaryQueueAppend <- queueOrder
		case order := <-channels.SWOrderFNMBackup:
			queueOrder := generateQueueOrder(order)
			channels.BackupQueueAppend <- queueOrder
		}
	}
}

func orderCompleteListener() {
	for {
		select {
		case orderComplete := <-channels.OrderCompleteFNM:
			var queueOrder = datatypes.QueueOrder{
				OrderType: orderComplete.OrderType,
				Floor:     orderComplete.Floor,
			}
			channels.PrimaryQueueRemove <- queueOrder
			channels.BackupQueueRemove <- queueOrder
			channels.ClearLightsFOM <- orderComplete
		case orderComplete := <-channels.OrderCompleteFFSM:
			channels.OrderCompleteFOM <- orderComplete
		}
	}
}

func orderRegisteredListener() {
	for {
		orderReg := <-channels.OrderRegisteredFNM

		channels.SetLightsFOM <- orderReg
	}
}
