package ordermanager

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
}

func receiver() {
	for {
		select {
		case swOrder := <-channels.SWOrderTOM:
			//Placeholder
			fmt.Println(swOrder)
		case hwOrder := <-channels.OrderFHM:
			orderRegHW(hwOrder)
		case costReq := <-channels.CostRequestTOM:
			//Placeholder
			fmt.Println(costReq)
		case orderComplete := <-channels.OrderCompleteTOM:
			//Placeholder
			fmt.Println(orderComplete)
		case orderRecAck := <-channels.OrderRecvAckTOM:
			//Placeholder
			fmt.Println(orderRecAck)
		}
	}
}

func orderRegHW(order datatypes.SWOrder) {
	//Make cost request
	var request datatypes.CostRequest
	request.Direction = order.Dir
	request.Floor = order.Floor
	request.Costsignature = "abc"
	//request.Costsignature = createCostSignature(order)

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
			if costAns.Costsignature != request.Costsignature {
				channels.CostAnswerTOM <- costAns
				break
			}
			fmt.Printf("Correct signature: %#v\n", costAns)
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

	return

}

func createCostSignature(order datatypes.SWOrder) string {
	timeSinceStart := time.Since(start)
	tStr := strconv.FormatInt(timeSinceStart.Nanoseconds()/1e6, 10)
	floorInt := int(order.Floor)
	dirInt := int(order.Dir)
	orderStr := "F:" + strconv.Itoa(floorInt) + "D:" + strconv.Itoa(dirInt)
	return orderStr + ":" + tStr
}

func TestOrderRegHW() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Hit enter to send dummyHWOrder")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	var dummyOrder datatypes.SWOrder
	dummyOrder.Floor = 2
	dummyOrder.Dir = 1

	var dummyCostAns1 datatypes.CostAnswer
	dummyCostAns1.CostValue = 5
	dummyCostAns1.SourceID = "Dummy1"
	dummyCostAns1.Costsignature = "123"

	var dummyCostAns2 datatypes.CostAnswer
	dummyCostAns2.CostValue = 7
	dummyCostAns2.SourceID = "Dummy2"
	dummyCostAns2.Costsignature = "abc"

	channels.OrderFHM <- dummyOrder

	channels.CostAnswerTOM <- dummyCostAns1
	time.Sleep(2 * time.Millisecond)
	channels.CostAnswerTOM <- dummyCostAns2

}
