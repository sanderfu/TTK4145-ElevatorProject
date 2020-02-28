package main

import (
	//"github.com/TTK4145/Network-go/network/conn"
	//"github.com/TTK4145/Network-go/network/bcast"
	//"github.com/TTK4145/Network-go/network/localip"
	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/driver-go/elevio"
	//"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"fmt"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/hwmanager"
)

func main() {
	fmt.Println("Starting HW Manager")

	hwmanager.Init()

	for {

	}

	fmt.Println("Done main")

}
