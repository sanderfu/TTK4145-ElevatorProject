package ordermanager

import (
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

////////////////////////////////////////////////////////////////////////////////
// Private variables
////////////////////////////////////////////////////////////////////////////////

var costRequestTimeoutMS time.Duration
var orderRecvAckWaitMS time.Duration
var maxCostValue int
var backupTakeoverTimeoutS time.Duration
var start time.Time

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////
// Private functions
////////////////////////////////////////////////////////////////////////////////

func orderRegistrationHW() {
	for {
		order := <-channels.OrderFhmTom
		//Make cost request
		var request = datatypes.CostRequest{
			OrderType: order.OrderType,
			Floor:     order.Floor,
		}
		//Broadcast cost request
		channels.CostRequestFomTnm <- request
		//Wait for answers
		done := time.After(costRequestTimeoutMS * time.Millisecond)
		primaryCost := maxCostValue + 1
		backupCost := maxCostValue + 1
	costWaitloop:
		for {
			select {
			case <-done:
				break costWaitloop
			case costAns := <-channels.CostAnswerFnmTom:
				// Assign primary and backup elevator
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
		if backupCost == maxCostValue+1 {
			order.BackupID = order.PrimaryID
		}
		channels.SWOrderFomTnm <- order
		//Wait for OrderRecAck from primary and backup
		done2 := time.After(orderRecvAckWaitMS * time.Millisecond)
		ackCounter := 0
	ackWaitloop:
		for {
			select {
			case <-done2:
				//Timer reached end, the order transmit is assumed to have failed and order is put back into the channel
				channels.OrderFhmTom <- order
				break ackWaitloop
			case orderRecvAck := <-channels.OrderRecvAckFnmTom:
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
					channels.OrderRegisteredFomTnm <- orderReg
					break ackWaitloop
				}
			}
		}
	}
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
		case order := <-channels.SWOrderPrimaryFnmTom:
			queueOrder := generateQueueOrder(order)
			channels.PrimaryQueueAppend <- queueOrder
		case order := <-channels.SWOrderBackupFnmTom:
			queueOrder := generateQueueOrder(order)
			channels.BackupQueueAppend <- queueOrder
		}
	}
}

func orderCompleteListener() {
	for {
		select {
		case orderComplete := <-channels.OrderCompleteFnmTom:
			if orderComplete.OrderType == datatypes.OrderInside && orderComplete.SourceID != orderComplete.ArrivalID {
				//The order is not for this elevator, discard it
				continue
			}
			var queueOrder = datatypes.QueueOrder{
				OrderType: orderComplete.OrderType,
				Floor:     orderComplete.Floor,
			}
			channels.PrimaryQueueRemove <- queueOrder
			channels.BackupQueueRemove <- queueOrder
			channels.ClearLightsFomThm <- orderComplete
		case orderComplete := <-channels.OrderCompleteFfsmTom:
			channels.OrderCompleteFomTnm <- orderComplete
		}
	}
}

func orderRegisteredListener() {
	for {
		orderReg := <-channels.OrderRegisteredFnmTom
		if orderReg.OrderType == datatypes.OrderInside && orderReg.SourceID != orderReg.ArrivalID {
			// Not our inside order, discards
		} else {
			channels.SetLightsFomThm <- orderReg
		}
	}
}
