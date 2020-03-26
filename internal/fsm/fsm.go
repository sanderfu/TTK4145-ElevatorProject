package fsm

import (
	"fmt"
	"os"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
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

var doorTimeout time.Duration

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

	doorTimeout = time.Duration(configuration.Config.DoorOpenDuration)
	// Wait for hardware manager to finish its setup
	hmInitStatus := <-channels.HMInitStatusFHM

	if !hmInitStatus {
		fmt.Println("Hardware Manager failed to initialize")
		os.Exit(1)
	}

	// Go down until elevator arrives at known floor
	hwmanager.SetElevatorDirection(datatypes.MotorDown)
	lastFloor = <-channels.CurrentFloorFHM
	hwmanager.SetElevatorDirection(datatypes.MotorStop)
	currentDir = datatypes.MotorStop

	go updateLastFloor()
	go costValueListener()

	currentState = datatypes.IdleState
}

////////////////////////////////////////////////////////////////////////////////
// State functions
////////////////////////////////////////////////////////////////////////////////

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

	// Check if elevator arrived at destination floor
	if currentOrder.Floor == lastFloor {
		hwmanager.SetElevatorDirection(datatypes.MotorStop)
		currentDir = datatypes.MotorStop
		doorOpeningTime = time.Now()
		hwmanager.SetDoorOpenLamp(true)

		// Inform order manager that order was completed on given floor
		completedOrder := datatypes.OrderComplete{
			Floor:     currentOrder.Floor,
			OrderType: currentOrder.OrderType,
		}
		channels.OrderCompleteFFSM <- completedOrder

		currentState = datatypes.DoorOpenState

		// Check if elevator arrived at a new floor and there is an order there
	} else if newFloorFlag == true {
		if ordermanager.OrderToTakeAtFloor(lastFloor, motorDirToOrderType(currentDir)) {
			hwmanager.SetElevatorDirection(datatypes.MotorStop)
			currentDir = datatypes.MotorStop
			doorOpeningTime = time.Now()
			hwmanager.SetDoorOpenLamp(true)

			// Inform order manager that order was completed on given floor
			completedOrder := datatypes.OrderComplete{
				Floor: lastFloor,
			}
			channels.OrderCompleteFFSM <- completedOrder

			currentState = datatypes.DoorOpenState
		}
		newFloorFlag = false
	}
}

func doorOpen() {
	if time.Since(doorOpeningTime) > doorTimeout*time.Second {
		hwmanager.SetDoorOpenLamp(false)
		currentState = datatypes.IdleState
	}
}

////////////////////////////////////////////////////////////////////////////////
// Other functions
////////////////////////////////////////////////////////////////////////////////

func motorDirToOrderType(dir int) int {
	if dir == datatypes.MotorUp {
		return datatypes.OrderUp
	} else if dir == datatypes.MotorDown {
		return datatypes.OrderDown
	} else {
		return datatypes.OrderInside
	}
}

func updateLastFloor() {
	for {
		floor := <-channels.CurrentFloorFHM
		if floor != lastFloor {
			lastFloor = floor
			newFloorFlag = true
		}
	}
}

//TODO: New function name
func costValueListener() {
	for {
		<-channels.FloorAndDirectionRequestFOM
		channels.FloorFFSM <- lastFloor
		channels.DirectionFFSM <- motorDirToOrderType(currentDir)
	}
}
