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

	configuration.ParseFlags()
	configuration.ReadConfig("./config.json")

	// start managers
	go watchdog.SenderNode(configuration.Flags.WatchdogPort)

	go networkmanager.NetworkManager()

	go ordermanager.OrderManager(configuration.Flags.LastPID)

	go hwmanager.HardwareManager(configuration.Flags.ElevatorPort)

	go fsm.FSM()

	//Go to sleep
	select {}

}
