package ordermanager

import (
	"fmt"
	"os"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/logger"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
)

const (
	answerWaitMS       = 250
	orderRecvAckWaitMS = 250
	maxCost            = 10
)

//Test variables
var failSendingAck = false
var costValue = 5

var start time.Time

var primaryQueue []datatypes.QueueOrder
var backupQueue []datatypes.QueueOrder

var primaryAppend chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)
var primaryRemove chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)

var backupAppend chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)
var backupRemove chan datatypes.QueueOrder = make(chan datatypes.QueueOrder)

//OrderManager ...
func OrderManager() {

	start = time.Now()

	fmt.Println("Hello go, this is Order Manager speaking!")

	go receiver()
	go orderRegHW()
	go orderRegSW()
	go queueModifier()
	go orderCompleteWatch()

}

func receiver() {
	for {
		select {
		case costReq := <-channels.CostRequestTOM:
			dummyCostAns(costReq)
		default:
		}
	}
}

func orderRegHW() {
	for {
		select {
		case order := <-channels.OrderFHM:
			{
				//Make cost request
				var request datatypes.CostRequest
				request.Direction = order.Dir
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
							if orderRecvAck.Floor == order.Floor && orderRecvAck.Dir == order.Dir {
								ackCounter++
							}
						}
						if ackCounter == 2 {
							//Transmit was successful
							break ackWaitloop
						}
					}
				}
			}
		}
	}
}

func ConfigureAndRunTest() {
	fmt.Println("I am process: ", os.Getpid())

	fmt.Println("Choose the cost value")
	fmt.Scan(&costValue)
	fmt.Println("The cost value of this process is: ", costValue)

	for {
		fmt.Println("Want to send dummy HW order?(y/n))")
		var ans string
		fmt.Scan(&ans)
		if ans == "y" || ans == "" {
			var dummyOrder datatypes.Order
			fmt.Println("Floor: ")
			fmt.Scan(&dummyOrder.Floor)
			fmt.Println("Dir (1/-1):")
			fmt.Scan(&dummyOrder.Dir)
			fmt.Println("Sending HW order")
			channels.OrderFHM <- dummyOrder
		}
		fmt.Println("Want to send dummy OrderComplete? (y/n))")
		fmt.Scan(&ans)
		if ans == "y" || ans == "" {
			var dummyComplete datatypes.OrderComplete
			fmt.Println("Floor: ")
			fmt.Scan(&dummyComplete.Floor)
			fmt.Println("Dir (1/-1):")
			fmt.Scan(&dummyComplete.Dir)
			fmt.Println("Sending OrderComplete")
			channels.OrderCompleteFOM <- dummyComplete
		}
	}
}

func dummyCostAns(costreq datatypes.CostRequest) {
	var costAns datatypes.CostAnswer
	costAns.CostValue = costValue
	costAns.DestinationID = costreq.SourceID
	channels.CostAnswerFOM <- costAns
}

func generateOrderRecvAck(queueOrder datatypes.QueueOrder) {
	var orderRecvAck datatypes.OrderRecvAck
	orderRecvAck.Dir = queueOrder.Dir
	orderRecvAck.Floor = queueOrder.Floor
	orderRecvAck.DestinationID = queueOrder.SourceID
	channels.OrderRecvAckFOM <- orderRecvAck
}

func generateQueueOrder(order datatypes.Order) datatypes.QueueOrder {
	var queueOrder datatypes.QueueOrder
	queueOrder.SourceID = order.SourceID
	queueOrder.Dir = order.Dir
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

func queueModifier() {
	for {
		select {
		case primaryOrder := <-primaryAppend:
			primary := true
			primaryQueue = append(primaryQueue, primaryOrder)
			logger.WriteLog(primaryQueue, primary, "/logs/")
			generateOrderRecvAck(primaryOrder)
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
			queueOrder.Dir = orderComplete.Dir
			fmt.Println("Forwarding remove request to queueModifier")
			queueOrder.Floor = orderComplete.Floor
			primaryRemove <- queueOrder
			backupRemove <- queueOrder
			fmt.Println("The remove request has been handeled")
		}
	}
}

func PeekFirstOrderInQueue(order *datatypes.QueueOrder) {

	if len(primaryQueue) > 0 {
		*order = primaryQueue[0]
		return
	}
	order = nil

	//var dummyOrder datatypes.QueueOrder
	//dummyOrder.Dir = datatypes.DOWN
	//dummyOrder.Floor = datatypes.FOURTH

	//*order = dummyOrder
}

// TODO: Make function for checking if elevator should stop at floor
