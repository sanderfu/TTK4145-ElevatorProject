package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/logger"
)

//"github.com/TTK4145/Network-go/network/peers"

func main() {
	var swOrderArray []datatypes.SWOrder
	var swOrderArray2 []datatypes.SWOrder
	for i := 0; i < 10; i++ {
		var testOrdre datatypes.SWOrder
		testOrdre.Signature = "TestSignature" + strconv.Itoa(i)
		testOrdre.SourceID = "TestSourceID"
		testOrdre.PrimaryID = "12345"
		testOrdre.BackupID = strconv.Itoa(i)
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.SECOND
		swOrderArray = append(swOrderArray, testOrdre)
	}
	fmt.Printf("%#v\n", swOrderArray)
	//go networkmanager.NetworkManager()
	//go ordermanager.OrderManager()
	//go ordermanager.ConfigureAndRunTest()

	logger.WriteLog(swOrderArray, "/logs/primaryqueue/")
	logger.ReadLogQueue(&swOrderArray2, "/logs/primaryqueue/")
	fmt.Printf("%#v\n", swOrderArray)
	for {
		time.Sleep(10 * time.Second)
	}
}
