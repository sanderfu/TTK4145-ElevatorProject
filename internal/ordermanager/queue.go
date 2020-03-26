package ordermanager

import (
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
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
		case primaryOrder := <-channels.PrimaryQueueAppend:
			primary := true
			if !orderInQueue(&primaryQueue, primaryOrder) {
				primaryQueue = append(primaryQueue, primaryOrder)
				logger.SaveQueue(primaryQueue, primary)
				generateOrderRecvAck(primaryOrder)
			}
		case primaryOrder := <-channels.PrimaryQueueRemove:
			primary := true
			removeFromQueue(&primaryQueue, primaryOrder)
			logger.SaveQueue(primaryQueue, primary)
		case backupOrder := <-channels.BackupQueueAppend:
			primary := false
			backupQueue = append(backupQueue, backupOrder)
			logger.SaveQueue(backupQueue, primary)
			generateOrderRecvAck(backupOrder)
		case backupOrder := <-channels.BackupQueueRemove:
			primary := false
			removeFromQueue(&backupQueue, backupOrder)
			logger.SaveQueue(backupQueue, primary)
		}
	}
}

func GetFirstOrderInQueue() datatypes.QueueOrder {
	return primaryQueue[0]
}

func QueueEmpty() bool {
	return len(primaryQueue) == 0
}

func OrderToTakeAtFloor(floor int, ordertype int) bool {

	for _, order := range primaryQueue {
		if order.Floor == floor && (order.OrderType == ordertype || order.OrderType == datatypes.OrderInside) {
			return true
		}
	}
	return false
}

func backupListener() {
	for {
		for _, elem := range backupQueue {
			if time.Since(elem.RegistrationTime) > backupTakeoverTimeoutS*time.Second {
				channels.BackupQueueRemove <- elem
				channels.PrimaryQueueAppend <- elem
			}
		}
		time.Sleep(1 * time.Second)
	}
}
