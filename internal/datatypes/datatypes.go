package datatypes

import "time"

// Datatypes goes here
//TODO: Fix everything to be CamelCase

//Basic types
type Floor int
type Direction int
type StructType int

const (
	FIRST  Floor = 0
	SECOND Floor = 1
	THIRD  Floor = 2
	FOURTH Floor = 3
)

// Change these to match with values from elevator_io (just for simplicity)
// Martin thinks this works. TODO: Fix this comment
const (
	UP     Direction = 0
	DOWN   Direction = 1
	INSIDE Direction = 2
)

const (
	MotorUp   Direction = 1
	MotorDown Direction = -1
	MotorStop Direction = 0
)

const ()

//Struct types
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values over the network.

type CostRequest struct {
	Signature string //Used by networkmanager to remove duplicates
	SourceID  string //ID of sender, to direct answer back.
	Floor     Floor
	Direction Direction
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
	Floor     Floor
	Dir       Direction
}

type OrderRecvAck struct {
	Signature     string
	SourceID      string
	DestinationID string
	Floor         Floor
	Dir           Direction
}

type OrderComplete struct {
	Signature string
	Floor     Floor
	Dir       Direction
}

type LightCommand struct {
	Signature string
	Floor     Floor
	Dir       Direction
}

type NWMMode int

const (
	Network   NWMMode = 0
	Localhost NWMMode = 1
)

type QueueOrder struct {
	SourceID         string
	Floor            Floor
	Dir              Direction
	RegistrationTime time.Time
}
