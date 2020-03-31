package ordermanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	primaryv1 = "/primaryv1.json"
	primaryv2 = "/primaryv2.json"

	backupv1 = "/backupv1.json"
	backupv2 = "/backupv2.json"

	assetDir = "./assets/"

	filePermissions = 0755
)

////////////////////////////////////////////////////////////////////////////////
// Private variables
////////////////////////////////////////////////////////////////////////////////

var primaryQueue []datatypes.QueueOrder
var backupQueue []datatypes.QueueOrder

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

// Only primary queue is accessable from outside the module
func OrderInQueue(order datatypes.QueueOrder) bool {
	return orderInQueue(&primaryQueue, order)
}

func GetFirstOrderInQueue() datatypes.QueueOrder {
	return primaryQueue[0]
}

func QueueEmpty() bool {
	return len(primaryQueue) == 0
}

////////////////////////////////////////////////////////////////////////////////
// Private functions
////////////////////////////////////////////////////////////////////////////////

// Listen for messages to add or remove elements from the queues, then save the
// updated queues
func queueModifier() {
	for {
		select {
		case primaryOrder := <-channels.PrimaryQueueAppend:
			primary := true
			addOrderToQueue(&primaryQueue, primaryOrder)
			saveQueue(primaryQueue, primary)
			sendOrderRecvAck(primaryOrder)
		case primaryOrder := <-channels.PrimaryQueueRemove:
			primary := true
			removeOrderFromQueue(&primaryQueue, primaryOrder)
			saveQueue(primaryQueue, primary)
		case backupOrder := <-channels.BackupQueueAppend:
			primary := false
			addOrderToQueue(&backupQueue, backupOrder)
			saveQueue(backupQueue, primary)
			sendOrderRecvAck(backupOrder)
		case backupOrder := <-channels.BackupQueueRemove:
			primary := false
			removeOrderFromQueue(&backupQueue, backupOrder)
			saveQueue(backupQueue, primary)
		}
	}
}

func orderInQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) bool {
	for _, orderFromQueue := range *queue {
		if ordersAreEqual(order, orderFromQueue) {
			return true
		}
	}
	return false
}

func ordersAreEqual(order1 datatypes.QueueOrder, order2 datatypes.QueueOrder) bool {
	return order1.Floor == order2.Floor && order1.OrderType == order2.OrderType
}

func addOrderToQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) {
	if !orderInQueue(queue, order) {
		*queue = append(*queue, order)
	}
}

func removeOrderFromQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) {
	j := 0
	for _, element := range *queue {
		if order.Floor != element.Floor {
			(*queue)[j] = element
			j++
		}
	}
	*queue = (*queue)[:j]
}

// Check if any backup orders have expired and if so move them into primary queue
func backupTimeoutListener() {
	for {
		for _, elem := range backupQueue {
			if time.Since(elem.RegistrationTime) > backupTakeoverTimeoutS*time.Second {
				channels.BackupQueueRemove <- elem
				channels.PrimaryQueueAppend <- elem
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func sendOrderRecvAck(queueOrder datatypes.QueueOrder) {
	var orderRecvAck = datatypes.OrderRecvAck{
		OrderType:     queueOrder.OrderType,
		Floor:         queueOrder.Floor,
		DestinationID: queueOrder.SourceID,
	}
	channels.OrderRecvAckFOM <- orderRecvAck
}

////////////////////////////////////////////////////////////////////////////////
// Functions for saving/loading queue
////////////////////////////////////////////////////////////////////////////////

func restoreQueues(lastPID string) {
	if lastPID != "NONE" {
		fmt.Println("Importing queue from crashed session")
		dir := "/" + lastPID
		loadQueue(&primaryQueue, true, dir)
		loadQueue(&backupQueue, false, dir)
		var orderReg datatypes.OrderRegistered
		saveQueue(primaryQueue, true)

		// Refresh order lights
		for i := 0; i < len(primaryQueue); i++ {
			orderReg.Floor = primaryQueue[i].Floor
			orderReg.OrderType = primaryQueue[i].OrderType
			channels.OrderRegisteredFOM <- orderReg
		}
		saveQueue(backupQueue, false)
		fmt.Println("Resume successful")
	}
}

func saveQueue(queue []datatypes.QueueOrder, primary bool) {
	pid := strconv.Itoa(os.Getpid())
	processAssetsDir := assetDir + pid

	result, err := json.MarshalIndent(queue, "", "")
	if err != nil {
		fmt.Println(err)
	}
	// If directory for storing queue does not exist, make it
	if _, err := os.Stat(processAssetsDir); os.IsNotExist(err) {
		err := os.MkdirAll(processAssetsDir, filePermissions)
		if err != nil {
			fmt.Println(err)
		}
	}
	// Store queue in new file, then delete the old one
	writefile, deletefile := selectFileNames(primary, pid)
	err = ioutil.WriteFile(processAssetsDir+writefile, result, filePermissions)
	if err != nil {
		fmt.Println(err)
	}
	os.Remove(processAssetsDir + deletefile)
}

func loadQueue(queue *[]datatypes.QueueOrder, primary bool, pid string) {
	_, readFile := selectFileNames(primary, pid)
	file, err := ioutil.ReadFile(assetDir + pid + "/" + readFile)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// Convert from JSON and load into queue
	_ = json.Unmarshal([]byte(file), queue)
}

func selectFileNames(primary bool, pid string) (string, string) {
	processAssetsDir := assetDir + pid

	if primary {
		if fileExists(processAssetsDir, primaryv1) {
			return primaryv2, primaryv1
		} else {
			return primaryv1, primaryv2
		}
	} else {
		if fileExists(processAssetsDir, backupv1) {
			return backupv2, backupv1
		} else {
			return backupv1, backupv2
		}
	}
}

func fileExists(dir string, filename string) bool {
	_, err := os.Stat(dir + filename)
	return err == nil
}
