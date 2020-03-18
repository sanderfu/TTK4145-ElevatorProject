package main

import (
	"encoding/json"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/ordermanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/watchdog"

	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/network/conn"
	//"github.com/TTK4145/Network-go/network/bcast"
	//"github.com/TTK4145/Network-go/network/localip"
	//"github.com/TTK4145/Network-go/network/peers"
	//"github.com/TTK4145/Network-go/driver-go/elevio"

	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/fsm"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/hwmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
)

func main() {

	
	// flag parsing
	elevPortPtr := flag.String("elevport", "", "elevator server port (mandatory)")
	watchdogPortPtr := flag.String("watchdogport", "", "watchdog port (mandatory)")
	lastPIDPtr := flag.String("lastpid", "NONE", "process ID of last running program")
	flag.Parse()
	
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Argument(s) missing. See -h")
		os.Exit(1)
	}
	

	fmt.Println("The Process ID is: ", os.Getpid())

	var resuming bool

	if *lastPIDPtr != "NONE" {
		resuming = true
	} else {
		resuming = false
	}
	fmt.Println("PID to resume from: ", *lastPIDPtr)

	// config parsing
	readConfig("./config.json")
	
	
	// start managers
	go watchdog.SenderNode(*watchdogPortPtr)
	
	go networkmanager.NetworkManager()

	go ordermanager.OrderManager(resuming, *lastPIDPtr)
	//go ordermanager.ConfigureAndRunTest()

	go hwmanager.HardwareManager(*elevPortPtr)
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
