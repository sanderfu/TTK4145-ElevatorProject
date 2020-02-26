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
	answerWaitMS = 1000
	maxCost      = 10
)

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
			//Placeholder
			fmt.Println(swOrder)
		case costReq := <-channels.CostRequestTOM:
			dummyCostAns(costReq, 5)
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
		waitLoop:
			for {
				select {
				case <-done:
					break waitLoop
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

func dummyCostAns(costreq datatypes.CostRequest, costValue int) {
	var costAns datatypes.CostAnswer
	costAns.CostValue = costValue
	costAns.DestinationID = costreq.SourceID

	channels.CostAnswerFOM <- costAns
}
