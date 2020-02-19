package datatypes
// Datatypes goes here

//Basic types
type Floor int
type Direction int

const (
	FIRST Floor = 0
	SECOND Floor = 1
	THIRD Floor = 2
	FOURTH Floor = 3
)

const (
	UP Direction = 1
	DOWN Direction = -1
	INSIDE Direction = 0
)

//Struct types

type Cost_request struct {
	Source_id string
	Floor Floor
	Direction Direction
}

type Cost_answer struct {
	Source_id string
	Cost_value float64
}

type SW_Order struct {
	Primary_id string
	Backup_id string
	Floor Floor
	Dir Direction
}

type Order_recv_ack struct {
	Source_id string
	Floor Floor
	Dir Direction
}

type Order_complete struct {
	Floor Floor
	Dir Direction
}