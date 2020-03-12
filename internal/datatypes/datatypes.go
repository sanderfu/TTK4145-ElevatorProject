package datatypes

import "time"

// Datatypes goes here
//TODO: Fix everything to be CamelCase

//Basic types
type StructType int
type State int

// Change these to match with values from elevator_io (just for simplicity)
// Martin thinks this works. TODO: Fix this comment
const (
	UP     int = 0
	DOWN   int = 1
	INSIDE int = 2
)

const (
	MotorUp   int = 1
	MotorDown int = -1
	MotorStop int = 0
)

const (
	IdleState     State = 0
	MovingState   State = 1
	DoorOpenState State = 2
)

//Struct types
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values over the network.

type CostRequest struct {
	Signature string //Used by networkmanager to remove duplicates
	SourceID  string //ID of sender, to direct answer back.
	Floor     int
	Direction int
}

type CostAnswer struct {
	Signature     string //Used by networkmanager to remove duplicates
	SourceID      string //ID of answer sender.
	DestinationID string //ID of answer receiver
	CostValue     int
}

type Order struct {
	Signature string
	SourceID  string
	PrimaryID string
	BackupID  string
	Floor     int
	Dir       int
}

type OrderRecvAck struct {
	Signature     string
	SourceID      string
	DestinationID string
	Floor         int
	Dir           int
}

type OrderComplete struct {
	Signature string
	Floor     int
	Dir       int
}

type LightCommand struct {
	Signature string
	Floor     int
	Dir       int
}

type NWMMode int

const (
	Network   NWMMode = 0
	Localhost NWMMode = 1
)

type QueueOrder struct {
	SourceID         string
	Floor            int
	Dir              int
	RegistrationTime time.Time
}

// Configuration struct
type Configuration struct {
	NumberOfFloors int
	ElevatorPort   int

	NetworkPacketDuplicates          int
	MaxUniqueSignatures              int
	UniqueSignatureRemovalPercentage int

	CostRequestTimeoutMS     int
	OrderReceiveAckTimeoutMS int
	MaxCostValue             int
	BackupTakeoverTimeoutS   int
}

var Config Configuration
