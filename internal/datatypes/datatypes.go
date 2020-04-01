package datatypes

import "time"

////////////////////////////////////////////////////////////////////////////////
// Definitions
////////////////////////////////////////////////////////////////////////////////

const (
	OrderUp     int = 0
	OrderDown   int = 1
	OrderInside int = 2
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

const (
	Network   NWMMode = 0
	Localhost NWMMode = 1
)

////////////////////////////////////////////////////////////////////////////////
// Data types
////////////////////////////////////////////////////////////////////////////////

type State int
type NWMMode int // Network module mode

// Cost structures

type CostRequest struct {
	Signature     string //Used by networkmanager to remove duplicates
	SourceID      string //ID of sender, to direct answer back.
	DestinationID string //ID of answer receiver
	Floor         int
	OrderType     int
}

type CostAnswer struct {
	Signature string //Used by networkmanager to remove duplicates
	SourceID  string //ID of answer sender.
	ArrivalID string //ID of where the answer arrived
	CostValue int
}

// Order structures

type Order struct {
	Signature string
	SourceID  string
	PrimaryID string
	BackupID  string
	Floor     int
	OrderType int
}

type OrderRecvAck struct { // order received acknowledgment
	Signature     string
	SourceID      string
	DestinationID string
	Floor         int
	OrderType     int
}

type OrderComplete struct {
	Signature string
	Floor     int
	OrderType int
}

type OrderRegistered struct {
	Signature string
	SourceID  string
	ArrivalID string
	Floor     int
	OrderType int
}

type QueueOrder struct {
	SourceID         string
	Floor            int
	OrderType        int
	RegistrationTime time.Time
}
