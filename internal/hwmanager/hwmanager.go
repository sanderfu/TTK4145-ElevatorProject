package hwmanager

import (
	"fmt"
	"time"

	"github.com/TTK4145/driver-go/elevio"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

var totalFloors int

func Init(numFloors int) {
	// TODO: Find out if this function should take addr and numFloors as args
	addr := "192.168.0.163:15657"
	totalFloors = numFloors

	elevio.Init(addr, numFloors)
	SetAllLights(false)

	go fsmMock()
	go omMock()

}

func PollCurrentFloor(curFloorChan chan<- datatypes.Floor) {

	floorSensorChan := make(chan int)
	go elevio.PollFloorSensor(floorSensorChan)

	for {
		floor := <-floorSensorChan

		elevio.SetFloorIndicator(floor)

		curFloorChan <- datatypes.Floor(floor)
	}

}

func PollHWORder(hwOrderChan chan<- datatypes.HW_Order) {

	btnChan := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(btnChan)

	for {

		btnValue := <-btnChan

		hwOrder := datatypes.HW_Order{
			Floor: datatypes.Floor(btnValue.Floor),
			Dir:   datatypes.Direction(btnValue.Button),
		}

		hwOrderChan <- hwOrder
	}
}

func SetLight(element datatypes.HW_Order, value bool) {
	elevio.SetButtonLamp(elevio.ButtonType(element.Dir), int(element.Floor),
		value)
}

func SetAllLights(value bool) {
	for floor := 0; floor < totalFloors; floor++ {
		for btn := elevio.BT_HallUp; btn <= elevio.BT_Cab; btn++ {
			if !(floor == 0 && btn == elevio.BT_HallDown) &&
				!(floor == totalFloors-1 && btn == elevio.BT_HallUp) {
				elevio.SetButtonLamp(btn, floor, value)
			}
		}
	}
}

func SetElevatorDirection(dir datatypes.Direction) {
	elevio.SetMotorDirection(elevio.MotorDirection(dir))
}

// Mocks below

func fsmMock() {
	go fsmPollFloorMock()
	go fsmSetElevatorDirectionMock()
}

func fsmPollFloorMock() {

	// Poll current floor
	floorChan := make(chan datatypes.Floor)

	go PollCurrentFloor(floorChan)

	for {
		floor := <-floorChan
		fmt.Println("Reached floor", floor)
	}
}

func fsmSetElevatorDirectionMock() {

	// Simulate an arbitrary sequence to see that directions are set correctly
	SetElevatorDirection(datatypes.MotorUp)
	time.Sleep(time.Second * 3)
	SetElevatorDirection(datatypes.MotorStop)
	time.Sleep(time.Second * 3)
	SetElevatorDirection(datatypes.MotorDown)
	time.Sleep(time.Second * 3)
	SetElevatorDirection(datatypes.MotorStop)
}

func omMock() {
	go omMockGetHWOrders()
}

func omMockGetHWOrders() {

	// Poll HW orders
	hwOrderChan := make(chan datatypes.HW_Order)

	go PollHWORder(hwOrderChan)

	for {
		hwOrder := <-hwOrderChan

		fmt.Println("HW Order: Floor", hwOrder.Floor, "Direction:", hwOrder.Dir)

		// Turn off that order again
		go omMockLightControl(hwOrder)
	}
}

func omMockLightControl(order datatypes.HW_Order) {

	// Set that light on
	SetLight(order, true)

	time.Sleep(time.Second * 3)

	SetLight(order, false)

}
