package networkmanager

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TTK4145/Network-go/network/bcast"
	"github.com/TTK4145/Network-go/network/localip"
	"github.com/TTK4145/Network-go/network/peers"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	packetduplicates = 10
)

var recentSignatures []string

//SWOrderTX for transmitting to network via driver
var SWOrderTX chan datatypes.SWOrder = make(chan datatypes.SWOrder)

//SWOrderRX for recieveing from network via driver
var SWOrderRX chan datatypes.SWOrder = make(chan datatypes.SWOrder)

//CostRequestTX for transmitting to network via driver
var CostRequestTX chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostRequestRX for recieveing from network via driver
var CostRequestRX chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostAnswerTX for transmitting to network via driver
var CostAnswerTX chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//CostAnswerRX for recieveing from network via driver
var CostAnswerRX chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//OrderRecvAckTX ...
var OrderRecvAckTX chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderRecvAckRX ...
var OrderRecvAckRX chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderCompleteTX ...
var OrderCompleteTX chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

//OrderCompleteRX ...
var OrderCompleteRX chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

//SWOrderTOM channel from Network Manager to Order Manager
var SWOrderTOM chan datatypes.SWOrder = make(chan datatypes.SWOrder)

//SWOrderFOM channel from Order Manager to Network Manager
var SWOrderFOM chan datatypes.SWOrder = make(chan datatypes.SWOrder)

//CostRequestTOM ...
var CostRequestTOM chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostRequestFOM ...
var CostRequestFOM chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostAnswerTOM ...
var CostAnswerTOM chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//CostAnswerFOM ...
var CostAnswerFOM chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//OrderRecvAckTOM ...
var OrderRecvAckTOM chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderRecvAckFOM ...
var OrderRecvAckFOM chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderCompleteTOM ...
var OrderCompleteTOM chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

//OrderCompleteFOM ...
var OrderCompleteFOM chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

//NetworkManager to start networkmanager routine.
func NetworkManager() {
	//Create an empty recentSignatures array
	recentSignatures = make([]string, 0)

	//Defining NetworkManager id based on IP and process ID
	var id string

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Receiver(16569, SWOrderRX, CostRequestRX, CostAnswerRX, OrderRecvAckRX, OrderCompleteRX)
	transmitter(16569)
}

func createSignature(structType int) string {
	t := time.Now()
	timeStr := t.Format("20060102150405")
	senderIPStr, _ := localip.LocalIP()
	return senderIPStr + "@" + timeStr + ":" + strconv.Itoa(structType)
}

func checkDuplicate(signature string) bool {
	for i := 0; i < len(recentSignatures); i++ {
		if recentSignatures[i] == signature {
			return true
		}
	}
	recentSignatures = append(recentSignatures, signature)
	return false
}

//TestSignatures tests that the signature system works as intended
func TestSignatures() {
	sign1 := createSignature(0)
	time.Sleep(1 * time.Second)
	sign2 := createSignature(2)
	if !checkDuplicate(sign1) {
		fmt.Println(sign1, " was not already in list")
	}
	if !checkDuplicate(sign2) {
		fmt.Println(sign2, " was not already in list")
	}
	if checkDuplicate(sign1) {
		fmt.Println(sign1, "was already in the list")
	}
}

//transmitter Function for applying packet redundancy before transmitting over network.
func transmitter(port int) {
	go bcast.Transmitter(16569, SWOrderTX, CostRequestTX, CostAnswerTX, OrderRecvAckTX, OrderCompleteTX)
	for {
		select {
		case order := <-SWOrderFOM:
			order.Signature = createSignature(0)
			for i := 0; i < packetduplicates; i++ {
				SWOrderTX <- order
			}
		case costReq := <-CostRequestFOM:
			costReq.Signature = createSignature(1)
			for i := 0; i < packetduplicates; i++ {
				CostRequestTX <- costReq
			}
		case costAns := <-CostAnswerFOM:
			costAns.Signature = createSignature(2)
			for i := 0; i < packetduplicates; i++ {
				CostAnswerTX <- costAns
			}
		case orderRecvAck := <-OrderRecvAckFOM:
			orderRecvAck.Signature = createSignature(3)
			for i := 0; i < packetduplicates; i++ {
				OrderRecvAckTX <- orderRecvAck
			}
		case orderComplete := <-OrderCompleteFOM:
			orderComplete.Signature = createSignature(4)
			for i := 0; i < packetduplicates; i++ {
				OrderCompleteTX <- orderComplete
			}
		}
	}
}

/*
func receiver(port int) {
	go bcast.Receiver(16569, SWOrderRX, CostRequestRX, CostAnswerRX, OrderRecvAckRX, OrderCompleteRX)
	for {
		select {
		case order := <-SWOrderRX:

		}
	}
}
*/

//TestSending Function to test basic order transmission over network
func TestSending() {
	for {
		var testOrdre datatypes.SWOrder
		testOrdre.PrimaryID = "12345"
		testOrdre.BackupID = "67890"
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.SECOND
		SWOrderTX <- testOrdre
		time.Sleep(1 * time.Second)
	}
}

//TestRecieving Function to test basic order transmission over network
func TestRecieving() {
	for {
		select {
		case order := <-SWOrderRX:
			fmt.Printf("Received: %#v\n", order)

		}
	}
}

/*
//TestSendingBuffered Function to test need for buffered RX channels
func TestSendingBuffered() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hit enter to start sending 10 SW orders")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	for i := 0; i < 10; i++ {
		fmt.Printf("Sending hardware order %v\n", i)
		var testOrdre datatypes.SWOrder
		testOrdre.Primary_id = "12345"
		testOrdre.Backup_id = "67890"
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.SECOND
		SWOrderTX <- testOrdre
		time.Sleep(1 * time.Second)
	}
}

/7testRecieveingBuffered Function to test need for buffered RX channels
func TestRecievingBuffered() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hit enter to start processing orders on SWOrderRX channel")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	for {
		select {
		case order := <-SWOrderRX:
			fmt.Printf("Received: %#v\n", order)
		}
	}
}
*/
