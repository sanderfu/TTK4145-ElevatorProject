package networkmanager

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TTK4145/Network-go/network/bcast"
	"github.com/TTK4145/Network-go/network/localip"
	. "github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	packetduplicates    = 10
	maxuniquesignatures = 25
	removeinclean       = int(maxuniquesignatures / 5)
)

var recentSignatures []string

var ip string

var start time.Time

var mode datatypes.NWMMode

//For some reason needs 1 in these struct{} channels to make it work
var killTransmitter = make(chan struct{}, 1)
var killReceiver = make(chan struct{}, 1)

var initTransmitter = make(chan struct{}, 1)
var initReceiver = make(chan struct{}, 1)

//NetworkManager to start networkmanager routine.
func NetworkManager() {
	//Start timer used for signatures
	start = time.Now()

	//Start networkWatch to detect connection loss (and switch to localhost)
	go networkWatch()

	//Initialize everything that need initializing
	recentSignatures = make([]string, 0)
	initTransmitter <- struct{}{}
	initReceiver <- struct{}{}
	InitDriverTX <- struct{}{}
	InitDriverRX <- struct{}{}
	mode = datatypes.Network
	ip, _ = localip.LocalIP()

	for {
		select {
		case <-initTransmitter:
			go transmitter(16569)
		case <-initReceiver:
			go receiver(16569)
		}
	}
}

func networkWatch() {
	for {
		time.Sleep(1000 * time.Millisecond)
		theIP, err := localip.LocalIP()
		fmt.Println("NetworkWatch checking state, the IP is", theIP)
		if err != nil {
			if mode != datatypes.Localhost {
				ip = "LOCALHOST"
				mode = datatypes.Localhost
				killTransmitter <- struct{}{}
				killReceiver <- struct{}{}
			}
		} else {
			if mode != datatypes.Network {
				ip = theIP
				mode = datatypes.Network
				killTransmitter <- struct{}{}
				killReceiver <- struct{}{}
			}
		}
	}
}

func createSignature(structType int) string {
	timeSinceStart := time.Since(start)
	t := strconv.FormatInt(timeSinceStart.Nanoseconds()/1e6, 10)
	senderIPStr := ip
	return senderIPStr + "@" + t + ":" + strconv.Itoa(structType)
}

func checkDuplicate(signature string) bool {
	for i := 0; i < len(recentSignatures); i++ {
		if recentSignatures[i] == signature {
			return true
		}
	}
	recentSignatures = append(recentSignatures, signature)
	if len(recentSignatures) > maxuniquesignatures {
		cleanArray()
	}
	return false
}

func cleanArray() {

	for i := 0; i < len(recentSignatures)-removeinclean; i++ {
		recentSignatures[i] = recentSignatures[i+removeinclean]
	}
	recentSignatures = recentSignatures[:len(recentSignatures)-removeinclean]
}

//TestSignatures tests that the signature system works as intended
func TestSignatures() {
	for i := 0; i < maxuniquesignatures*2; i++ {
		sign1 := createSignature(i)
		checkDuplicate(sign1)
		printRecentSignatures()
	}
}

func printRecentSignatures() {
	fmt.Println("")
	fmt.Println("Recentsignatures:")
	for j := 0; j < len(recentSignatures); j++ {
		fmt.Println(recentSignatures[j])
	}
}

//transmitter Function for applying packet redundancy before transmitting over network.
func transmitter(port int) {
	go bcast.Transmitter(port, mode, SWOrderTX, CostRequestTX, CostAnswerTX, OrderRecvAckTX, OrderCompleteTX)
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
		case <-killTransmitter:
			KillDriverTX <- struct{}{}
			initTransmitter <- struct{}{}
			return
		}
	}
}

func receiver(port int) {
	go bcast.Receiver(port, SWOrderRX, CostRequestRX, CostAnswerRX, OrderRecvAckRX, OrderCompleteRX)
	for {
		select {
		case order := <-SWOrderRX:
			if !checkDuplicate(order.Signature) {
				SWOrderTOM <- order
			}
		case costReq := <-CostRequestRX:
			if !checkDuplicate(costReq.Signature) {
				CostRequestTOM <- costReq
			}
		case costAns := <-CostAnswerRX:
			if !checkDuplicate(costAns.Signature) {
				CostAnswerTOM <- costAns
			}
		case orderRecvAck := <-OrderRecvAckRX:
			if !checkDuplicate(orderRecvAck.Signature) {
				OrderRecvAckTOM <- orderRecvAck
			}
		case orderComplete := <-OrderCompleteRX:
			if !checkDuplicate(orderComplete.Signature) {
				OrderCompleteTOM <- orderComplete
			}
		case <-killReceiver:
			KillDriverRX <- struct{}{}
			initReceiver <- struct{}{}
			return
		}
	}
}

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

//TestSendingRedundant Function to test transmitting redundancy measures to packet loss.
func TestSendingRedundant(ordersToSend int) {
	//Create dummy order from order manager
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Hit enter to start sending %v SW from 'Order Manager'", ordersToSend)
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	for i := 0; i < ordersToSend; i++ {
		fmt.Printf("Sending order %v\n", i)
		var testOrdre datatypes.SWOrder
		testOrdre.PrimaryID = "RedundantPackage"
		testOrdre.BackupID = strconv.Itoa(i)
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.SECOND
		SWOrderFOM <- testOrdre
		time.Sleep(1 * time.Second)
	}

	for i := 0; i < ordersToSend; i++ {
		fmt.Printf("Sending order %v\n", i)
		var testOrdre datatypes.SWOrder
		testOrdre.Signature = createSignature(0)
		testOrdre.PrimaryID = "NonRedundantPackage"
		testOrdre.BackupID = strconv.Itoa(i)
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.SECOND
		SWOrderTX <- testOrdre
		time.Sleep(1 * time.Second)
	}

}

//TestReceivingRedundant Function to test that order manager gets unique packages only
func TestReceivingRedundant(runtime time.Duration) {
	countRedundant := 0
	countNonRedundant := 0
	done := time.After(runtime * time.Second)
readLoop:
	for {
		select {
		case order := <-SWOrderTOM:
			if order.PrimaryID == "RedundantPackage" {
				countRedundant++
			} else if order.PrimaryID == "NonRedundantPackage" {
				countNonRedundant++
			}
			fmt.Printf("Received: %#v\n", order)
		case <-done:
			break readLoop
		}
	}
	fmt.Printf("Test results:\n Unique RedundantPackages received: %v\n Unique NonRedundantPackages received %v\n", countRedundant, countNonRedundant)
	printRecentSignatures()
}
