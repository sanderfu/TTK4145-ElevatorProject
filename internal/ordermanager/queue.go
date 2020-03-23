package ordermanager

import (
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/logger"
)

var primaryQueue []datatypes.QueueOrder
var backupQueue []datatypes.QueueOrder

func removeFromQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) {
	j := 0
	for _, element := range *queue {
		if order.Floor != element.Floor {
			(*queue)[j] = element
			j++
		}
	}

	*queue = (*queue)[:j]
}

func orderInQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) bool {
	for _, elem := range *queue {
		if elem.Floor == order.Floor && elem.OrderType == order.OrderType {
			return true
		}
	}
	return false
}

func queueModifier() {
	for {
		select {
		case primaryOrder := <-primaryAppend:
			primary := true
			if !orderInQueue(&primaryQueue, primaryOrder) {
				primaryQueue = append(primaryQueue, primaryOrder)
				logger.WriteLog(primaryQueue, primary, "/logs/")
				generateOrderRecvAck(primaryOrder)
			}
		case primaryOrder := <-primaryRemove:
			primary := true
			removeFromQueue(&primaryQueue, primaryOrder)
			logger.WriteLog(primaryQueue, primary, "/logs/")
		case backupOrder := <-backupAppend:
			primary := false
			backupQueue = append(backupQueue, backupOrder)
			logger.WriteLog(backupQueue, primary, "/logs/")
			generateOrderRecvAck(backupOrder)
		case backupOrder := <-backupRemove:
			primary := false
			removeFromQueue(&backupQueue, backupOrder)
			logger.WriteLog(backupQueue, primary, "/logs/")
		}
	}
}

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

func backupListener() {
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
