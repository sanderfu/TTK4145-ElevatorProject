package networkmanager

import (
	"os"
	"strconv"
	"time"

	"github.com/TTK4145/Network-go/network/bcast"
	"github.com/TTK4145/Network-go/network/localip"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

var broadCastPort int
var packetDuplicates int
var maxUniqueSignatures int
var removePercent int
var recentSignatures []string
var localID string // IP address and process ID
var start time.Time
var mode datatypes.NWMMode

//NetworkManager to start networkmanager routine.
func NetworkManager() {
	//Update global variables based on configuration
	broadCastPort = configuration.Config.BroadcastPort
	packetDuplicates = configuration.Config.NetworkPacketDuplicates
	maxUniqueSignatures = configuration.Config.MaxUniqueSignatures
	removePercent = configuration.Config.UniqueSignatureRemovalPercentage

	//Start timer used for signatures
	start = time.Now()

	//Start connectionWatchdog to detect connection loss (and switch to localhost)
	go connectionWatchdog()

	//Initialize everything that need initializing
	recentSignatures = make([]string, 0)
	channels.InitTransmitter <- struct{}{}
	channels.InitReceiver <- struct{}{}
	channels.InitDriverTX <- struct{}{}
	channels.InitDriverRX <- struct{}{}
	mode = datatypes.Network
	localID, _ = localip.LocalIP()
	localID += ":" + strconv.Itoa(os.Getpid())

	for {
		select {
		case <-channels.InitTransmitter:
			go transmitter(broadCastPort)
		case <-channels.InitReceiver:
			go receiver(broadCastPort)
		}
	}
}

func connectionWatchdog() {
	for {
		time.Sleep(1000 * time.Millisecond)
		IPAddr, err := localip.LocalIP()
		if err != nil {
			if mode != datatypes.Localhost {
				localID = "LOCALHOST" + ":" + strconv.Itoa(os.Getpid())
				mode = datatypes.Localhost
				channels.KillTransmitter <- struct{}{}
				channels.KillReceiver <- struct{}{}
			}
		} else {
			if mode != datatypes.Network {
				localID = IPAddr + ":" + strconv.Itoa(os.Getpid())
				mode = datatypes.Network
				channels.KillTransmitter <- struct{}{}
				channels.KillReceiver <- struct{}{}
			}
		}
	}
}

func createUniqueSignature() string {
	//Delay the signing process by 1ms to guarantee unique signatures
	time.Sleep(1 * time.Millisecond)
	timeStamp := strconv.FormatInt(time.Since(start).Nanoseconds()/1e6, 10)
	return localID + "@" + timeStamp
}

// Return true on success, false if signature already exists in recentSignatures
func addSignature(signature string) bool {
	for i := 0; i < len(recentSignatures); i++ {
		if recentSignatures[i] == signature {
			return false
		}
	}
	if len(recentSignatures) > maxUniqueSignatures {
		// Remove a percentage (removePercent) from the from of recentSignatures
		firstIndex := int(maxUniqueSignatures * removePercent / 100)
		recentSignatures = recentSignatures[firstIndex:]
	}
	recentSignatures = append(recentSignatures, signature)
	return true
}

//transmitter Function for applying packet redundancy before transmitting over network.
func transmitter(port int) {
	go bcast.Transmitter(port, mode, channels.SWOrderTX, channels.CostRequestTX,
		channels.CostAnswerTX, channels.OrderRecvAckTX, channels.OrderCompleteTX,
		channels.OrderRegisteredTX)
	for {
		select {
		case order := <-channels.SWOrderFOM:
			order.Signature = createUniqueSignature()
			order.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.SWOrderTX <- order
			}
		case costReq := <-channels.CostRequestFOM:
			costReq.Signature = createUniqueSignature()
			costReq.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.CostRequestTX <- costReq
			}
		case costAns := <-channels.CostAnswerFOM:
			costAns.Signature = createUniqueSignature()
			costAns.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.CostAnswerTX <- costAns
			}
		case orderRecvAck := <-channels.OrderRecvAckFOM:
			orderRecvAck.Signature = createUniqueSignature()
			orderRecvAck.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.OrderRecvAckTX <- orderRecvAck
			}
		case orderComplete := <-channels.OrderCompleteFOM:
			orderComplete.Signature = createUniqueSignature()
			for i := 0; i < packetDuplicates; i++ {
				channels.OrderCompleteTX <- orderComplete
			}
		case orderRegistered := <-channels.OrderRegisteredFOM:
			orderRegistered.Signature = createUniqueSignature()
			for i := 0; i < packetDuplicates; i++ {
				channels.OrderRegisteredTX <- orderRegistered
			}
		case <-channels.KillTransmitter:
			channels.KillDriverTX <- struct{}{}
			channels.InitTransmitter <- struct{}{}
			return
		}
	}
}

func receiver(port int) {
	go bcast.Receiver(port, channels.SWOrderRX, channels.CostRequestRX,
		channels.CostAnswerRX, channels.OrderRecvAckRX,
		channels.OrderCompleteRX, channels.OrderRegisteredRX)
	for {
		select {
		case order := <-channels.SWOrderRX:
			if addSignature(order.Signature) {
				if order.PrimaryID == localID && order.BackupID == localID {
					channels.SWOrderFNMPrimary <- order
					channels.SWOrderFNMBackup <- order
				} else if order.PrimaryID == localID {
					channels.SWOrderFNMPrimary <- order
				} else if order.BackupID == localID {
					channels.SWOrderFNMBackup <- order
				}
			}
		case costReq := <-channels.CostRequestRX:
			if addSignature(costReq.Signature) {
				costReq.DestinationID = localID
				channels.CostRequestFNM <- costReq
			}
		case costAns := <-channels.CostAnswerRX:
			if costAns.DestinationID != localID {
				continue
			}
			if addSignature(costAns.Signature) {
				channels.CostAnswerFNM <- costAns
			}
		case orderRecvAck := <-channels.OrderRecvAckRX:
			if orderRecvAck.DestinationID != localID {
				continue
			}
			if addSignature(orderRecvAck.Signature) {
				channels.OrderRecvAckFNM <- orderRecvAck
			}
		case orderComplete := <-channels.OrderCompleteRX:
			if addSignature(orderComplete.Signature) {
				channels.OrderCompleteFNM <- orderComplete
			}
		case orderRegistered := <-channels.OrderRegisteredRX:
			if addSignature(orderRegistered.Signature) {
				channels.OrderRegisteredFNM <- orderRegistered
			}
		case <-channels.KillReceiver:
			channels.KillDriverRX <- struct{}{}
			channels.InitReceiver <- struct{}{}
			return
		}
	}
}
