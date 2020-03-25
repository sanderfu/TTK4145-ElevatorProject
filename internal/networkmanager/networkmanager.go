package networkmanager

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TTK4145/Network-go/network/bcast"
	"github.com/TTK4145/Network-go/network/localip"
	. "github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

var packetDuplicates int
var maxUniqueSignatures int
var removePercent int

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
	//Update global variables based on configuration
	packetDuplicates = configuration.Config.NetworkPacketDuplicates
	maxUniqueSignatures = configuration.Config.MaxUniqueSignatures
	removePercent = configuration.Config.UniqueSignatureRemovalPercentage

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
	ip += ":" + strconv.Itoa(os.Getpid())

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
		if err != nil {
			if mode != datatypes.Localhost {
				ip = "LOCALHOST" + ":" + strconv.Itoa(os.Getpid())
				mode = datatypes.Localhost
				killTransmitter <- struct{}{}
				killReceiver <- struct{}{}
			}
		} else {
			if mode != datatypes.Network {
				ip = theIP + ":" + strconv.Itoa(os.Getpid())
				mode = datatypes.Network
				killTransmitter <- struct{}{}
				killReceiver <- struct{}{}
			}
		}
	}
}

func createSignature(structType int) string {
	//Delay the signing process by 1ms to guarantee unique signatures
	time.Sleep(1 * time.Millisecond)
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
	if len(recentSignatures) > maxUniqueSignatures {
		cleanArray()
	}
	recentSignatures = append(recentSignatures, signature)
	return false
}

// Remove a percentage (removePercent) from the front of recentSignatures
func cleanArray() {
	firstIndex := int(maxUniqueSignatures * removePercent / 100)

	recentSignatures = recentSignatures[firstIndex:]
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
	go bcast.Transmitter(port, mode, SWOrderTX, CostRequestTX, CostAnswerTX, OrderRecvAckTX, OrderCompleteTX, OrderRegisteredTX)
	for {
		select {
		case order := <-SWOrderFOM:
			order.Signature = createSignature(0)
			order.SourceID = ip
			for i := 0; i < packetDuplicates; i++ {
				SWOrderTX <- order
			}
		case costReq := <-CostRequestFOM:
			costReq.Signature = createSignature(1)
			costReq.SourceID = ip
			for i := 0; i < packetDuplicates; i++ {
				CostRequestTX <- costReq
			}
		case costAns := <-CostAnswerFOM:
			costAns.Signature = createSignature(2)
			costAns.SourceID = ip
			for i := 0; i < packetDuplicates; i++ {
				CostAnswerTX <- costAns
			}
		case orderRecvAck := <-OrderRecvAckFOM:
			orderRecvAck.Signature = createSignature(3)
			orderRecvAck.SourceID = ip
			for i := 0; i < packetDuplicates; i++ {
				OrderRecvAckTX <- orderRecvAck
			}
		case orderComplete := <-OrderCompleteFOM:
			orderComplete.Signature = createSignature(4)
			for i := 0; i < packetDuplicates; i++ {
				OrderCompleteTX <- orderComplete
			}
		case orderRegistered := <-OrderRegisteredFOM:
			orderRegistered.Signature = createSignature(5)
			for i := 0; i < packetDuplicates; i++ {
				OrderRegisteredTX <- orderRegistered
			}
		case <-killTransmitter:
			KillDriverTX <- struct{}{}
			initTransmitter <- struct{}{}
			return
		}
	}
}

func receiver(port int) {
	go bcast.Receiver(port, SWOrderRX, CostRequestRX, CostAnswerRX, OrderRecvAckRX, OrderCompleteRX, OrderRegisteredRX)
	for {
		select {
		case order := <-SWOrderRX:
			if !checkDuplicate(order.Signature) {
				if order.PrimaryID == ip && order.BackupID == ip {
					SWOrderFNMPrimary <- order
					SWOrderFNMBackup <- order
				} else if order.PrimaryID == ip {
					SWOrderFNMPrimary <- order
				} else if order.BackupID == ip {
					SWOrderFNMBackup <- order
				}
			}
		case costReq := <-CostRequestRX:
			if !checkDuplicate(costReq.Signature) {
				costReq.DestinationID = ip
				CostRequestFNM <- costReq
			}
		case costAns := <-CostAnswerRX:
			if costAns.DestinationID != ip {
				continue
			}
			if !checkDuplicate(costAns.Signature) {
				CostAnswerFNM <- costAns
			}
		case orderRecvAck := <-OrderRecvAckRX:
			if orderRecvAck.DestinationID != ip {
				continue
			}
			if !checkDuplicate(orderRecvAck.Signature) {
				OrderRecvAckFNM <- orderRecvAck
			}
		case orderComplete := <-OrderCompleteRX:
			if !checkDuplicate(orderComplete.Signature) {
				OrderCompleteFNM <- orderComplete
			}
		case orderRegistered := <-OrderRegisteredRX:
			if !checkDuplicate(orderRegistered.Signature) {
				OrderRegisteredFNM <- orderRegistered
			}
		case <-killReceiver:
			KillDriverRX <- struct{}{}
			initReceiver <- struct{}{}
			return
		}
	}
}
