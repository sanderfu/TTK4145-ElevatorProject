package main

import (
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/ordermanager"
)

//"github.com/TTK4145/Network-go/network/peers"

func main() {
	go networkmanager.NetworkManager()
	go ordermanager.OrderManager()
	go ordermanager.ConfigureAndRunTest()
	for {
		time.Sleep(10 * time.Second)
	}
}
