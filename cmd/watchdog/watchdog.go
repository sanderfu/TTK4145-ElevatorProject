package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/watchdog"
)

func main() {

	// flag parsing
	elevPortPtr := flag.String("elevport", "", "elevator server port (mandatory)")
	watchdogPortPtr := flag.String("watchdogport", "", "watchdog port (mandatory)")
	bcastlocalPortPtr := flag.String("bcastlocalport", "NONE", "Port for local host broadcast")
	flag.Parse()
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Argument(s) missing. See -h")
		os.Exit(1)
	}

	go watchdog.WatchdogNode(*watchdogPortPtr, *elevPortPtr, *bcastlocalPortPtr)

	//Block thread from exiting
	select {}
}
