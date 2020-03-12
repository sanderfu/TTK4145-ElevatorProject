package main

import (
	"encoding/json"
	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/network/conn"
	//"github.com/TTK4145/Network-go/network/bcast"
	//"github.com/TTK4145/Network-go/network/localip"
	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/driver-go/elevio"

	"fmt"
	"os"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/fsm"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/hwmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/ordermanager"
)

func main() {

	readConfig("./config.json")

	go networkmanager.NetworkManager()

	go ordermanager.OrderManager()
	//go ordermanager.ConfigureAndRunTest()

	go hwmanager.HardwareManager()
	go fsm.FSM()

	for {
		time.Sleep(10 * time.Second)
	}

}

func readConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&datatypes.Config)
	if err != nil {
		fmt.Println(err)
	}
}
