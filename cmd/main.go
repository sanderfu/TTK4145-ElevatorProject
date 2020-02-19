package main

import (
	//"github.com/TTK4145/Network-go/network/conn"
	//"github.com/TTK4145/Network-go/network/bcast"
	//"github.com/TTK4145/Network-go/network/localip"
	//"github.com/TTK4145/Network-go/network/peers"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
)

func main() {
	go networkmanager.NetworkManager()
	go networkmanager.TestSending()
	go networkmanager.TestRecieving()
	for {
	}
}
