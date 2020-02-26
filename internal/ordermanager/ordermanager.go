package ordermanager

import (
	"fmt"
	"os"
	"strconv"
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
			fmt.Printf("Received: %#v\n", swOrder)
			dummyOrderRecvAck(swOrder, failSendingAck)
		case costReq := <-channels.CostRequestTOM:
			fmt.Printf("Received: %#v\n", costReq)
			dummyCostAns(costReq)
			fmt.Println("DummyCostAns replied")
		case orderComplete := <-channels.OrderCompleteTOM:
			//Placeholder
			fmt.Println(orderComplete)
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
						fmt.Printf("%#v\n", costAns)
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
				fmt.Println("Done choosing primary and backup")
				fmt.Printf("%#v\n", order)
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
						fmt.Printf("%#v\n", orderRecvAck)
						if orderRecvAck.SourceID == order.PrimaryID || orderRecvAck.SourceID == order.BackupID {
							//Check that ack matches order, if not throw it away as it has probably arrived to late for prev. order
							if orderRecvAck.Floor == order.Floor && orderRecvAck.Dir == order.Dir {
								ackCounter++
							}
							//If we only wait for ack from ourselves, we append the extra ackCount here
							if order.PrimaryID == order.BackupID {
								ackCounter++
							}
						}
						if ackCounter == 2 {
							//Transmit was successful
							fmt.Println("Order transmit was successfull")
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

	fmt.Println("Choose if Ack should fail (YES=1/NO=0)")
	var fail string
	fmt.Scan(&fail)
	failSendingAck, _ = strconv.ParseBool(fail)
	fmt.Println("Should the OrderRecvAck fail: ", failSendingAck)
	for {
		fmt.Println("Want to send dummy HW order?(y/n))")
		var ans string
		fmt.Scan(&ans)
		if ans == "y" || ans == "" {
			fmt.Println("Sending HW order")
			var dummyOrder datatypes.SWOrder
			dummyOrder.Floor = 2
			dummyOrder.Dir = 1

			channels.OrderFHM <- dummyOrder
		} else {
			break
		}
	}
	fmt.Println("The configuration is done!")
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
		orderRecvAck.DestinationID = swOrder.SourceID
		channels.OrderRecvAckFOM <- orderRecvAck
	} else {
		return
	}
}
