package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Configuration struct
type Configuration struct {
	NumberOfFloors int

	NetworkPacketDuplicates          int
	MaxUniqueSignatures              int
	UniqueSignatureRemovalPercentage int

	CostRequestTimeoutMS     int
	OrderReceiveAckTimeoutMS int
	MaxCostValue             int
	BackupTakeoverTimeoutS   int
}

type CommandLineFlags struct {
	ElevatorPort string
	WatchdogPort string
	LastPID      string
}

// Global variable storing all data from config file
var Config Configuration
var Flags CommandLineFlags

func ReadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println(err)
	}
}

func ParseFlags() {
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
	Flags.ElevatorPort = *elevPortPtr
	Flags.WatchdogPort = *watchdogPortPtr
	Flags.LastPID = *lastPIDPtr
}
