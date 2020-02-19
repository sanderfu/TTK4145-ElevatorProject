package datatypes

// Datatypes goes here

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

const (
	UP     Direction = 1
	DOWN   Direction = -1
	INSIDE Direction = 0
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
