package main

import (
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/ordermanager"
	//"github.com/TTK4145/Network-go/network/peers"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
)

func main() {

	go networkmanager.NetworkManager()
	go ordermanager.OrderManager()

	go ordermanager.ConfigureAndRunTest()
	for {

	}
}
