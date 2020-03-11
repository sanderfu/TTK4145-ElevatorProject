package hwmanager

import (
	"fmt"
	"time"

	"github.com/TTK4145/Network-go/network/localip"
	"github.com/TTK4145/driver-go/elevio"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	numFloors = 4
)

var totalFloors int

func HardwareManager() {

	setup(4)

	go pollCurrentFloor()
	go pollHWORder()
	go lightWatch()

}

func setup(numFloors int) {
	// TODO: Find out if this function should take addr and numFloors as args
	addr, err := localip.LocalIP()

	if err != nil {
		fmt.Println("Error: hwmanager (setup):", err)
	}

	addr += ":15657"
	totalFloors = numFloors

	elevio.Init(addr, numFloors)
	for floor := datatypes.FIRST; floor < datatypes.Floor(totalFloors); floor++ {
		setAllLightsAtFloor(floor, false)
	}
	SetDoorOpenLamp(false)

	channels.HMInitStatusTFSM <- true
	//go fsmMock()
	//go omMock()

}

func pollCurrentFloor() {

	floorSensorChan := make(chan int)
	go elevio.PollFloorSensor(floorSensorChan)

	for {
		floor := <-floorSensorChan

		elevio.SetFloorIndicator(floor)

		channels.CurrentFloorTFSM <- datatypes.Floor(floor)
	}

}

func pollHWORder() {

	btnChan := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(btnChan)

	for {

		btnValue := <-btnChan
		fmt.Println("Button pressed")
		hwOrder := datatypes.Order{
			Floor: datatypes.Floor(btnValue.Floor),
			Dir:   datatypes.Direction(btnValue.Button),
		}

		channels.OrderFHM <- hwOrder
	}
}

func setLight(element datatypes.Order, value bool) {
	elevio.SetButtonLamp(elevio.ButtonType(element.Dir), int(element.Floor),
		value)
}

func setAllLightsAtFloor(floor datatypes.Floor, value bool) {
	for btn := datatypes.UP; btn <= datatypes.INSIDE; btn++ {
		if !(int(floor) == 0 && btn == datatypes.DOWN) &&
			!(int(floor) == totalFloors-1 && btn == datatypes.UP) {
			elevio.SetButtonLamp(elevio.ButtonType(btn), int(floor), value)
		}
	}

}

func SetElevatorDirection(dir datatypes.Direction) {
	elevio.SetMotorDirection(elevio.MotorDirection(dir))
}

func SetDoorOpenLamp(value bool) {
	elevio.SetDoorOpenLamp(value)
}

func lightWatch() {
	for {
		select {
		case orderComplete := <-channels.OrderCompleteTHM:
			setAllLightsAtFloor(orderComplete.Floor, false)
		case orderRegistered := <-channels.OrderRegisteredTHM:
			setLight(orderRegistered, true)
		}
	}
}

// Mocks below

func fsmMock() {
	go fsmPollFloorMock()
	go fsmsetElevatorDirectionMock()
}

func fsmPollFloorMock() {

	for {
		floor := <-channels.CurrentFloorTFSM
		fmt.Println("Reached floor", floor)
	}
}

func fsmsetElevatorDirectionMock() {

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
	for {
		hwOrder := <-channels.OrderFHM

		fmt.Println("HW Order: Floor", hwOrder.Floor, "Direction:", hwOrder.Dir)

		// Turn off that order again
		go omMockLightControl(hwOrder)
	}
}

func omMockLightControl(order datatypes.Order) {

	// Set that light on
	setLight(order, true)

	time.Sleep(time.Second * 3)

	setLight(order, false)

}
