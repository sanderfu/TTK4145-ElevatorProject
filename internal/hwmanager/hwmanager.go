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
}

func GetCurrentState(curFloorChan chan<- datatypes.Floor) {

	floorSensorChan := make(chan int)
	go elevio.PollFloorSensor(floorSensorChan)

	for {
		floor := <-floorSensorChan

		fmt.Println("Arrived at floor", floor)

		elevio.SetFloorIndicator(floor)

		curFloorChan <- datatypes.Floor(floor)
	}

}

func GetHWORder(hwOrderChan chan<- datatypes.HW_Order) {

	// Poll buttons
	btnChan := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(btnChan)

	for {

		btnValue := <-btnChan

		fmt.Println("Floor:", btnValue.Floor, "Button:", btnValue.Button)

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
