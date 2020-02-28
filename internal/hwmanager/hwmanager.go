package hwmanager

import (
	"fmt"

	"github.com/TTK4145/driver-go/elevio"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

func Init() {
	// TODO: Find out if this function should take addr and numFloors as args
	addr := "192.168.0.163:15657"
	numFloors := 4

	elevio.Init(addr, numFloors)

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

func SetOrderComplete(completedOrder datatypes.Order_complete) {

	fmt.Println("Got Order Complete message from Order Manager")
	elevio.SetButtonLamp(elevio.ButtonType(completedOrder.Dir),
		int(completedOrder.Floor), false)

}

func convertOrderDirToMotorDir(dir datatypes.Direction) datatypes.Direction {
	switch dir {
	case datatypes.UP:
		dir = 1
	case datatypes.DOWN:
		dir = -1
	case datatypes.INSIDE:
		dir = 0
	}
	return dir
}

// Mocks below

func fsmMock() {
	go fsmPollFloorMock()
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

//func fsmSetDirectionMock() {
//
//	for {
//
//	}
//}

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
	}
}
