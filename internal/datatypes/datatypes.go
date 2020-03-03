package datatypes

// Datatypes goes here
//TODO: Fix everything to be CamelCase

//Basic types
type Floor int
type Direction int

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

//Struct types

type Cost_request struct {
	Source_id string
	Floor     Floor
	Direction Direction
}

type Cost_answer struct {
	Source_id  string
	Cost_value float64
}

type SW_Order struct {
	Primary_id string
	Backup_id  string
	Floor      Floor
	Dir        Direction
}

type Order_recv_ack struct {
	Source_id string
	Floor     Floor
	Dir       Direction
}

// This type can be the same for Order_complete and HW_Order since they will be
// the same (Martin thinks). So TODO: merge these into one
type Order_complete struct {
	Floor Floor
	Dir   Direction
}

type HW_Order struct {
	Floor Floor
	Dir   Direction
}
