package datatypes

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
	Signature string
	SourceID  string
	Floor     Floor
	Direction Direction
}

type CostAnswer struct {
	Signature string
	SourceID  string
	CostValue float64
}

type SWOrder struct {
	Signature string
	PrimaryID string
	BackupID  string
	Floor     Floor
	Dir       Direction
}

type OrderRecvAck struct {
	Signature string
	SourceID  string
	Floor     Floor
	Dir       Direction
}

type OrderComplete struct {
	Signature string
	Floor     Floor
	Dir       Direction
}

type NWMMode int

const (
	Network   NWMMode = 0
	Localhost NWMMode = 1
)
