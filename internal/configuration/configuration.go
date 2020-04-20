package configuration

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Configuration struct {
	NumberOfFloors   int
	DoorOpenDuration int

	BroadcastPort                    int
	NetworkPacketDuplicates          int
	MaxUniqueSignatures              int
	UniqueSignatureRemovalPercentage int

	CostRequestTimeoutMS     int
	OrderReceiveAckTimeoutMS int
	MaxCostValue             int
	BackupTakeoverTimeoutS   int
}

type CommandLineFlags struct {
	ElevatorPort   string
	WatchdogPort   string
	LastPID        string
	BcastLocalPort string
}

////////////////////////////////////////////////////////////////////////////////
// Global variables storing data from config file and command line flags
////////////////////////////////////////////////////////////////////////////////

var Config Configuration
var Flags CommandLineFlags

////////////////////////////////////////////////////////////////////////////////
// Public functions
////////////////////////////////////////////////////////////////////////////////

func ReadConfig(filename string) {
	// Open config file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read config file into global variable
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ParseFlags() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Argument(s) missing. See -h")
		os.Exit(1)
	}

	// Define flags
	elevPortPtr := flag.String("elevport", "", "elevator server port (mandatory)")
	watchdogPortPtr := flag.String("watchdogport", "", "watchdog port (mandatory)")
	lastPIDPtr := flag.String("lastpid", "NONE", "process ID of last running program")
	bcastlocalPortPtr := flag.String("bcastlocalport", "NONE", "Port for local host broadcast")
	flag.Parse()

	// Load flags into global variable
	Flags.ElevatorPort = *elevPortPtr
	Flags.WatchdogPort = *watchdogPortPtr
	Flags.LastPID = *lastPIDPtr
	Flags.BcastLocalPort = *bcastlocalPortPtr
}
