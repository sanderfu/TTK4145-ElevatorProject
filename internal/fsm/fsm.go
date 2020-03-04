package fsm

import (
	"fmt"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/hwmanager"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/ordermanager"
)

var totalFloors int

var lastFloor datatypes.Floor
var newFloorFlag bool
var currentDir datatypes.Direction

var currentOrder datatypes.QueueOrder
var currentState datatypes.State

var doorOpeningTime time.Time

const doorTimeout = 3

func Run() {
	for {
		switch currentState {
		case datatypes.IdleState:
			Idle()
		case datatypes.MovingState:
			Moving()
		case datatypes.DoorOpenState:
			DoorOpen()
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func Init(numFloors int) {

	fmt.Println("Starting FSM")

	totalFloors = numFloors

	hwmanager.Init(numFloors)

	// Go down until elevator arrives at known floor
	hwmanager.SetElevatorDirection(datatypes.MotorDown)
	go hwmanager.PollCurrentFloor(channels.CurrentFloorTFSM)
	lastFloor = <-channels.CurrentFloorTFSM
	hwmanager.SetElevatorDirection(datatypes.MotorStop)
	currentDir = datatypes.MotorStop

	fmt.Println("Came to rest at floor", lastFloor, "with last dir", currentDir)

	go updateLastFloor()

	Run()
}

func Idle() {

	//fmt.Println("State idle")

	// Check for new orders
	var firstOrder datatypes.QueueOrder
	ordermanager.PeekFirstOrderInQueue(&firstOrder)

	if &firstOrder == nil {
		return
	}

	currentOrder = firstOrder

	// Calculate direction to move in
	if currentOrder.Floor > lastFloor {
		currentDir = datatypes.MotorUp
	} else if currentOrder.Floor < lastFloor {
		currentDir = datatypes.MotorDown
	} else {
		currentDir = datatypes.MotorStop
	}

	// Start moving
	hwmanager.SetElevatorDirection(currentDir)

	currentState = datatypes.MovingState
}

func Moving() {

	//fmt.Println("State moving")

	// Check if we arrived at destination floor
	if currentOrder.Floor == lastFloor {
		fmt.Println("Arrived at floor", lastFloor)
		hwmanager.SetElevatorDirection(datatypes.MotorStop)
		doorOpeningTime = time.Now()
		hwmanager.SetDoorOpenLamp(true)
		currentState = datatypes.DoorOpenState

		// TODO: Tell order manager that the order is complete
	}

	// TODO: Check if elevator should stop at floor if newFloorFlag == true
}

func DoorOpen() {
	fmt.Println("The door should now open")

	if time.Since(doorOpeningTime) > doorTimeout*time.Second {
		fmt.Println("Door closing")
		hwmanager.SetDoorOpenLamp(false)
		currentState = datatypes.IdleState
	}

}

func updateLastFloor() {
	for {
		lastFloor = <-channels.CurrentFloorTFSM
		newFloorFlag = true
	}
}

// States
// Idle
// Moving
// Door open
//
//
