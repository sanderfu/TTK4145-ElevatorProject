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

var lastFloor int
var newFloorFlag bool
var currentDir int

var currentOrder datatypes.QueueOrder
var currentState datatypes.State

var doorOpeningTime time.Time

const doorTimeout = 3

func FSM() {

	fsmInit()

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

func fsmInit() {

	// Wait for hardware manager to finish its setup
	hmInitStatus := <-channels.HMInitStatusTFSM

	if !hmInitStatus {
		fmt.Println("Hardware Manager failed to initialize")
		os.Exit(1)
	}

	// Go down until elevator arrives at known floor
	hwmanager.SetElevatorDirection(datatypes.MotorDown)
	lastFloor = <-channels.CurrentFloorTFSM
	hwmanager.SetElevatorDirection(datatypes.MotorStop)
	currentDir = datatypes.MotorStop

	go updateLastFloor()

	currentState = datatypes.IdleState
}

func idle() {

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

	// Check if we arrived at destination floor
	if currentOrder.Floor == lastFloor {
		hwmanager.SetElevatorDirection(datatypes.MotorStop)
		doorOpeningTime = time.Now()
		hwmanager.SetDoorOpenLamp(true)

		// Tell order manager that order was completed on given floor
		completedOrder := datatypes.OrderComplete{
			Floor:     currentOrder.Floor,
			OrderType: currentOrder.OrderType,
		}
		channels.OrderCompleteFFSM <- completedOrder

		currentState = datatypes.DoorOpenState

	} else if newFloorFlag == true {
		// Check if we arrived at a new floor and there is an order there
		if ordermanager.OrderToTakeAtFloor(lastFloor, motorDirToOrderDir(currentDir)) {
			hwmanager.SetElevatorDirection(datatypes.MotorStop)
			doorOpeningTime = time.Now()
			hwmanager.SetDoorOpenLamp(true)

			// Tell order manager that order was completed on given floor
			completedOrder := datatypes.OrderComplete{
				Floor: lastFloor,
			}
			channels.OrderCompleteFFSM <- completedOrder

			currentState = datatypes.DoorOpenState
		}
		newFloorFlag = false
	}
}

func motorDirToOrderDir(dir int) int {
	if dir == datatypes.MotorUp {
		return datatypes.UP
	} else if dir == datatypes.MotorDown {
		return datatypes.DOWN
	} else {
		return datatypes.INSIDE
	}
}

func doorOpen() {
	if time.Since(doorOpeningTime) > doorTimeout*time.Second {
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
