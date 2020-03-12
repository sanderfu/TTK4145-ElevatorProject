package hwmanager

import (
	"fmt"
	"strconv"

	"github.com/TTK4145/Network-go/network/localip"
	"github.com/TTK4145/driver-go/elevio"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

var numberOfFloors int

func HardwareManager() {
	setup()

	go pollCurrentFloor()
	go pollHWORder()
	go lightWatch()
}

func setup() {
	addr, err := localip.LocalIP()

	if err != nil {
		fmt.Println("Error: hwmanager (setup):", err)
	}

	addr += ":" + strconv.Itoa(datatypes.Config.ElevatorPort)
	numberOfFloors = datatypes.Config.NumberOfFloors
	elevio.Init(addr, numberOfFloors)

	for floor := 0; floor < numberOfFloors; floor++ {
		setAllLightsAtFloor(floor, false)
	}
	SetDoorOpenLamp(false)

	channels.HMInitStatusTFSM <- true
}

func pollCurrentFloor() {

	floorSensorChan := make(chan int)
	go elevio.PollFloorSensor(floorSensorChan)

	for {
		floor := <-floorSensorChan
		elevio.SetFloorIndicator(floor)
		channels.CurrentFloorTFSM <- floor
	}
}

func pollHWORder() {

	btnChan := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(btnChan)

	for {
		btnValue := <-btnChan
		hwOrder := datatypes.Order{
			Floor: btnValue.Floor,
			Dir:   int(btnValue.Button),
		}
		channels.OrderFHM <- hwOrder
	}
}

func setLight(element datatypes.Order, value bool) {
	elevio.SetButtonLamp(elevio.ButtonType(element.Dir), int(element.Floor),
		value)
}

func setAllLightsAtFloor(floor int, value bool) {
	for btn := datatypes.UP; btn <= datatypes.INSIDE; btn++ {
		if !(int(floor) == 0 && btn == datatypes.DOWN) &&
			!(int(floor) == numberOfFloors-1 && btn == datatypes.UP) {
			elevio.SetButtonLamp(elevio.ButtonType(btn), int(floor), value)
		}
	}

}

func SetElevatorDirection(dir int) {
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
