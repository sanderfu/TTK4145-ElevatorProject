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

////////////////////////////////////////////////////////////////////////////////
// Private variables
////////////////////////////////////////////////////////////////////////////////

var broadCastPort int
var packetDuplicates int
var maxUniqueSignatures int
var removePercent int
var recentSignatures []string
var localID string // IP address and process ID
var start time.Time
var mode datatypes.NWMMode

////////////////////////////////////////////////////////////////////////////////
// Public function
////////////////////////////////////////////////////////////////////////////////

func NetworkManager() {
	// Get parameters from config
	broadCastPort = configuration.Config.BroadcastPort
	packetDuplicates = configuration.Config.NetworkPacketDuplicates
	maxUniqueSignatures = configuration.Config.MaxUniqueSignatures
	removePercent = configuration.Config.UniqueSignatureRemovalPercentage

	mode = datatypes.Network
	localID, _ = localip.LocalIP()
	localID += ":" + strconv.Itoa(os.Getpid())
	recentSignatures = make([]string, 0)

	//Start timer used for signatures
	start = time.Now()

	//Start connectionWatchdog to detect change in network (online/offline)
	go connectionWatchdog()

	//Send initialize messages on all relevant control-signal channels.
	channels.InitTransmitter <- struct{}{}
	channels.InitReceiver <- struct{}{}
	channels.InitDriverTX <- struct{}{}
	channels.InitDriverRX <- struct{}{}

	for {
		select {
		case <-channels.InitTransmitter:
			go transmitter(broadCastPort)
		case <-channels.InitReceiver:
			go receiver(broadCastPort)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// Private functions
////////////////////////////////////////////////////////////////////////////////

func connectionWatchdog() {
	for {
		time.Sleep(1000 * time.Millisecond)
		IPAddr, err := localip.LocalIP()
		if err != nil {
			//Not connected to internet, take action if has not taken action already
			if mode != datatypes.Localhost {
				localID = "LOCALHOST" + ":" + strconv.Itoa(os.Getpid())
				mode = datatypes.Localhost
				channels.KillTransmitter <- struct{}{}
				channels.KillReceiver <- struct{}{}
			}
		} else {
			//Connected to internet, take action if has not taken action already
			if mode != datatypes.Network {
				localID = IPAddr + ":" + strconv.Itoa(os.Getpid())
				mode = datatypes.Network
				channels.KillTransmitter <- struct{}{}
				channels.KillReceiver <- struct{}{}
			}
		}
	}
}

//Create a unique signature for the package.
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

// Function for applying packet redundancy before transmitting over network.
func transmitter(port int) {
	go bcast.Transmitter(port, mode, channels.SWOrderTX, channels.CostRequestTX,
		channels.CostAnswerTX, channels.OrderRecvAckTX, channels.OrderCompleteTX,
		channels.OrderRegisteredTX)
	for {
		select {
		case order := <-channels.SWOrderFomTnm:
			order.Signature = createUniqueSignature()
			order.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.SWOrderTX <- order
			}
		case costReq := <-channels.CostRequestFomTnm:
			costReq.Signature = createUniqueSignature()
			costReq.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.CostRequestTX <- costReq
			}
		case costAns := <-channels.CostAnswerFomTnm:
			costAns.Signature = createUniqueSignature()
			costAns.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.CostAnswerTX <- costAns
			}
		case orderRecvAck := <-channels.OrderRecvAckFomTnm:
			orderRecvAck.Signature = createUniqueSignature()
			orderRecvAck.SourceID = localID
			for i := 0; i < packetDuplicates; i++ {
				channels.OrderRecvAckTX <- orderRecvAck
			}
		case orderComplete := <-channels.OrderCompleteFomTnm:
			orderComplete.Signature = createUniqueSignature()
			for i := 0; i < packetDuplicates; i++ {
				channels.OrderCompleteTX <- orderComplete
			}
		case orderRegistered := <-channels.OrderRegisteredFomTnm:
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

// Function for removing redundant packages such that only unique packages
// recieved are communicated onwards in the system to the Order Manager
func receiver(port int) {
	go bcast.Receiver(port, channels.SWOrderRX, channels.CostRequestRX,
		channels.CostAnswerRX, channels.OrderRecvAckRX,
		channels.OrderCompleteRX, channels.OrderRegisteredRX)
	for {
		select {
		case order := <-channels.SWOrderRX:
			if addSignature(order.Signature) {
				if order.PrimaryID == localID && order.BackupID == localID {
					channels.SWOrderPrimaryFnmTom <- order
					channels.SWOrderBackupFnmTom <- order
				} else if order.PrimaryID == localID {
					channels.SWOrderPrimaryFnmTom <- order
				} else if order.BackupID == localID {
					channels.SWOrderBackupFnmTom <- order
				}
			}
		case costReq := <-channels.CostRequestRX:
			if addSignature(costReq.Signature) {
				costReq.DestinationID = localID
				channels.CostRequestFnmTom <- costReq
			}
		case costAns := <-channels.CostAnswerRX:
			if costAns.DestinationID != localID {
				continue
			}
			if addSignature(costAns.Signature) {
				channels.CostAnswerFnmTom <- costAns
			}
		case orderRecvAck := <-channels.OrderRecvAckRX:
			if orderRecvAck.DestinationID != localID {
				continue
			}
			if addSignature(orderRecvAck.Signature) {
				channels.OrderRecvAckFnmTom <- orderRecvAck
			}
		case orderComplete := <-channels.OrderCompleteRX:
			if addSignature(orderComplete.Signature) {
				channels.OrderCompleteFnmTom <- orderComplete
			}
		case orderRegistered := <-channels.OrderRegisteredRX:
			if addSignature(orderRegistered.Signature) {
				channels.OrderRegisteredFnmTom <- orderRegistered
			}
		case <-channels.KillReceiver:
			channels.KillDriverRX <- struct{}{}
			channels.InitReceiver <- struct{}{}
			return
		}
	}
}
