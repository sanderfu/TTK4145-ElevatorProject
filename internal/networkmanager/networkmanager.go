package networkmanager

import (
	"fmt"
	"os"
	"time"

	"github.com/TTK4145/Network-go/network/bcast"
	"github.com/TTK4145/Network-go/network/localip"
	"github.com/TTK4145/Network-go/network/peers"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
}

// We make channels for sending and receiving our custom data types
//SW Order:
var SWOrderTX chan datatypes.SW_Order = make(chan datatypes.SW_Order)
var SWOrderRX chan datatypes.SW_Order = make(chan datatypes.SW_Order)

//Cost request:
var CostRequestTX chan datatypes.Cost_request = make(chan datatypes.Cost_request)
var CostRequestRX chan datatypes.Cost_request = make(chan datatypes.Cost_request)

//Cost answer:
var CostAnswerTX chan datatypes.Cost_answer = make(chan datatypes.Cost_answer)
var CostAnswerRX chan datatypes.Cost_answer = make(chan datatypes.Cost_answer)

//Order_recv_ack:
var OrderRecvAckTX chan datatypes.Order_recv_ack = make(chan datatypes.Order_recv_ack)
var OrderRecvAckRX chan datatypes.Order_recv_ack = make(chan datatypes.Order_recv_ack)

//Order complete:
var OrderCompleteTX chan datatypes.Order_complete = make(chan datatypes.Order_complete)
var OrderCompleteRX chan datatypes.Order_complete = make(chan datatypes.Order_complete)

func NetworkManager() {
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

	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, helloTx, SWOrderTX, CostRequestTX, CostAnswerTX, OrderRecvAckTX, OrderCompleteTX)
	go bcast.Receiver(16569, helloRx, SWOrderRX, CostRequestRX, CostAnswerRX, OrderRecvAckRX, OrderCompleteRX)

	/*
		// The example message. We just send one of these every second.
		go func() {
			helloMsg := HelloMsg{"Hello from " + id, 0}
			for {
				helloMsg.Iter++
				helloTx <- helloMsg
				time.Sleep(1 * time.Second)
			}
		}()

		fmt.Println("Started")
		for {
			select {
			case p := <-peerUpdateCh:
				fmt.Printf("Peer update:\n")
				fmt.Printf("  Peers:    %q\n", p.Peers)
				fmt.Printf("  New:      %q\n", p.New)
				fmt.Printf("  Lost:     %q\n", p.Lost)

			case a := <-helloRx:
				fmt.Printf("Received: %#v\n", a)
			}
		}
	*/
}

func TestSending() {
	for {
		var testOrdre datatypes.SW_Order
		testOrdre.Primary_id = "12345"
		testOrdre.Backup_id = "67890"
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.SECOND
		SWOrderTX <- testOrdre
		time.Sleep(1 * time.Second)
	}
}

func TestRecieving() {
	for {
		select {
		case order := <-SWOrderRX:
			fmt.Printf("Received: %#v\n", order)

		}
	}
}
