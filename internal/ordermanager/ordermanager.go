package ordermanager

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
)

const (
	answerWaitMS       = 1000
	orderRecvAckWaitMS = 1000
	maxCost            = 10
)

//Test variables
var failSendingAck = false
var costValue = 5

var start time.Time

var startRegistration = make(chan struct{}, 1)

func OrderManager() {

	start = time.Now()

	fmt.Println("Hello go, this is Order Manager speaking!")

	go receiver()
	go orderRegHW()
}

func receiver() {
	for {
		select {
		case swOrder := <-channels.SWOrderTOM:
			dummyOrderRecvAck(swOrder, failSendingAck)
			fmt.Printf("Received: %#v\n", swOrder)
		case costReq := <-channels.CostRequestTOM:
			dummyCostAns(costReq)
		case orderComplete := <-channels.OrderCompleteTOM:
			//Placeholder
			fmt.Println(orderComplete)
		case orderRecAck := <-channels.OrderRecvAckTOM:
			//Placeholder
			fmt.Println(orderRecAck)
		}
	}
}

func orderRegHW() {
	for {
		select {
		case order := <-channels.OrderFHM:
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
					fmt.Printf("%#v\n", costAns)
					if costAns.CostValue < primaryCost {
						backupCost = primaryCost
						primaryCost = costAns.CostValue
						order.PrimaryID = order.BackupID
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
			done = time.After(orderRecvAckWaitMS * time.Millisecond)
			ackCounter := 0
		ackWaitloop:
			for {
				select {
				case <-done:
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

func TestOrderRegHW() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Hit enter to send dummyHWOrder")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	var dummyOrder datatypes.SWOrder
	dummyOrder.Floor = 2
	dummyOrder.Dir = 1

	channels.OrderFHM <- dummyOrder

}

func dummyCostAns(costreq datatypes.CostRequest) {
	var costAns datatypes.CostAnswer
	costAns.CostValue = costValue
	costAns.DestinationID = costreq.SourceID

	channels.CostAnswerFOM <- costAns
}

func dummyOrderRecvAck(swOrder datatypes.SWOrder, fail bool) {
	if !fail {
		var orderRecvAck datatypes.OrderRecvAck
		orderRecvAck.Dir = swOrder.Dir
		orderRecvAck.Floor = swOrder.Floor
		channels.OrderRecvAckFOM <- orderRecvAck
	} else {
		return
	}
}
