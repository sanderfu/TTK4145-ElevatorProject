package fsm

import (
	"fmt"
	"os"
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

func FSM() {

	// Add init here
	fsmInit(4)

	for {
		switch currentState {
		case datatypes.IdleState:
			idle()
		case datatypes.MovingState:
			moving()
		case datatypes.DoorOpenState:
			doorOpen()
		default:
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func fsmInit(numFloors int) {
	fmt.Println("Starting FSM")

	totalFloors = numFloors

	// Wait for hardware manager to finish its setup
	hmInitStatus := <-channels.HMInitStatusTFSM

	if !hmInitStatus {
		println("Hardware Manager failed to initialize")
		os.Exit(1)
	}

	// Go down until elevator arrives at known floor
	hwmanager.SetElevatorDirection(datatypes.MotorDown)
	lastFloor = <-channels.CurrentFloorTFSM
	hwmanager.SetElevatorDirection(datatypes.MotorStop)
	currentDir = datatypes.MotorStop

	fmt.Println("Came to rest at floor", lastFloor, "with last dir", currentDir)

	go updateLastFloor()

	currentState = datatypes.IdleState
}

func idle() {

	fmt.Println("State idle")

	// Check for new orders
	if ordermanager.QueueEmpty() {
		return
	}

	currentOrder = ordermanager.GetFirstOrderInQueue()

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

func moving() {

	fmt.Println("State moving")

	// Check if we arrived at destination floor
	if currentOrder.Floor == lastFloor {
		fmt.Println("Arrived at destination floor", lastFloor)
		hwmanager.SetElevatorDirection(datatypes.MotorStop)
		doorOpeningTime = time.Now()
		hwmanager.SetDoorOpenLamp(true)

		// Tell order manager that order was completed on given floor
		completedOrder := datatypes.OrderComplete{
			Floor: lastFloor,
		}
		channels.OrderCompleteTOM <- completedOrder

		currentState = datatypes.DoorOpenState
	} else if newFloorFlag == true {
		// Check if we arrived at a new floor and there is an order there
		if ordermanager.OrderToTakeAtFloor(lastFloor, currentDir) {
			fmt.Println("Stopping at floor even though its not destination")
			hwmanager.SetElevatorDirection(datatypes.MotorStop)
			doorOpeningTime = time.Now()
			hwmanager.SetDoorOpenLamp(true)

			// Tell order manager that order was completed on given floor
			completedOrder := datatypes.OrderComplete{
				Floor: lastFloor,
			}
			channels.OrderCompleteTOM <- completedOrder

			currentState = datatypes.DoorOpenState
		}
		newFloorFlag = false
	}
}

func doorOpen() {
	fmt.Println("State door open")

	if time.Since(doorOpeningTime) > doorTimeout*time.Second {
		fmt.Println("Door closing")
		hwmanager.SetDoorOpenLamp(false)
		currentState = datatypes.IdleState
	}

}

func updateLastFloor() {
	for {
		floor := <-channels.CurrentFloorTFSM
		if floor != lastFloor {
			lastFloor = floor
			newFloorFlag = true
		}
	}
}
