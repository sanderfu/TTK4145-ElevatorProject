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

	permissionRW = 0755
)

var primaryQueue []datatypes.QueueOrder
var backupQueue []datatypes.QueueOrder

func removeFromQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) {
	j := 0
	for _, element := range *queue {
		if order.Floor != element.Floor {
			(*queue)[j] = element
			j++
		}
	}

	*queue = (*queue)[:j]
}

func orderInQueue(queue *[]datatypes.QueueOrder, order datatypes.QueueOrder) bool {
	for _, elem := range *queue {
		if elem.Floor == order.Floor && elem.OrderType == order.OrderType {
			return true
		}
	}
	return false
}

func queueModifier() {
	for {
		select {
		case primaryOrder := <-channels.PrimaryQueueAppend:
			primary := true
			if !orderInQueue(&primaryQueue, primaryOrder) {
				primaryQueue = append(primaryQueue, primaryOrder)
				saveQueue(primaryQueue, primary)
				generateOrderRecvAck(primaryOrder)
			}
		case primaryOrder := <-channels.PrimaryQueueRemove:
			primary := true
			removeFromQueue(&primaryQueue, primaryOrder)
			saveQueue(primaryQueue, primary)
		case backupOrder := <-channels.BackupQueueAppend:
			primary := false
			backupQueue = append(backupQueue, backupOrder)
			saveQueue(backupQueue, primary)
			generateOrderRecvAck(backupOrder)
		case backupOrder := <-channels.BackupQueueRemove:
			primary := false
			removeFromQueue(&backupQueue, backupOrder)
			saveQueue(backupQueue, primary)
		}
	}
}

func GetFirstOrderInQueue() datatypes.QueueOrder {
	return primaryQueue[0]
}

func QueueEmpty() bool {
	return len(primaryQueue) == 0
}

func OrderToTakeAtFloor(floor int, ordertype int) bool {

	for _, order := range primaryQueue {
		if order.Floor == floor && (order.OrderType == ordertype || order.OrderType == datatypes.OrderInside) {
			return true
		}
	}
	return false
}

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

////////////////////////////////////////////////////////////////////////////////
// Functions for queue backup
////////////////////////////////////////////////////////////////////////////////

func restoreQueues(lastPID string) {
	if lastPID != "NONE" {
		fmt.Println("Importing queue from crashed session")
		dir := "/" + lastPID
		loadQueue(&primaryQueue, true, dir)
		loadQueue(&backupQueue, false, dir)
		var orderReg datatypes.OrderRegistered
		saveQueue(primaryQueue, true)

		// Set lights
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
	if _, err := os.Stat(processAssetsDir); os.IsNotExist(err) {
		err := os.MkdirAll(processAssetsDir, permissionRW)
		if err != nil {
			fmt.Println(err)
		}
	}
	writefile, deletefile := selectFileNames(primary, pid)
	err = ioutil.WriteFile(processAssetsDir+writefile, result, permissionRW)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove(processAssetsDir + deletefile)
	if err != nil {
		//fmt.Println(err)
	}
}

func loadQueue(queue *[]datatypes.QueueOrder, primary bool, pid string) {
	_, readFile := selectFileNames(primary, pid)
	file, err := ioutil.ReadFile(assetDir + pid + "/" + readFile)
	if err != nil {
		fmt.Println("Error: ", err)
	}
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
