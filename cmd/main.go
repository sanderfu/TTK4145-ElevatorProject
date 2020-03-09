package main

import (
	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/network/conn"
	//"github.com/TTK4145/Network-go/network/bcast"
	//"github.com/TTK4145/Network-go/network/localip"
	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/driver-go/elevio"
	//"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"

	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/hwmanager"
)

//"github.com/TTK4145/Network-go/network/peers"

func main() {

	// go networkmanager.NetworkManager()
	// go ordermanager.OrderManager()
	// go ordermanager.ConfigureAndRunTest()

	go hwmanager.HardwareManager()

	for {
		time.Sleep(10 * time.Second)
	}

}
