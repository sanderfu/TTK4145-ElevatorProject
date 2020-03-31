package main

import (
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/fsm"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/hwmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/ordermanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/watchdog"
)

func main() {

	// initialize system parameters
	configuration.ParseFlags()
	configuration.ReadConfig("./config.json")

	// start managers
	go watchdog.ElevatorNode(configuration.Flags.WatchdogPort)

	go networkmanager.NetworkManager()

	go ordermanager.OrderManager()

	go hwmanager.HardwareManager()

	go fsm.FSM()

	// block program for exiting
	select {}

}
