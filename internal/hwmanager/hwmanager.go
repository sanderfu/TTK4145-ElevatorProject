package hwmanager

import (
	"github.com/TTK4145/driver-go/elevio"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/configuration"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

var numberOfFloors int

func HardwareManager(port string) {
	hwInit(port)

	go pollCurrentFloor()
	go pollHWORder()
	go updateOrderLights()
}

func SetElevatorDirection(dir int) {
	elevio.SetMotorDirection(elevio.MotorDirection(dir))
}

func SetDoorOpenLamp(value bool) {
	elevio.SetDoorOpenLamp(value)
}

func hwInit(port string) {
	addr := ":" + port
	numberOfFloors = configuration.Config.NumberOfFloors
	elevio.Init(addr, numberOfFloors)

	for floor := 0; floor < numberOfFloors; floor++ {
		setAllLightsAtFloor(floor, false)
	}
	SetDoorOpenLamp(false)

	channels.HMInitStatusFHM <- true
}

func pollCurrentFloor() {

	floorSensorChan := make(chan int)
	go elevio.PollFloorSensor(floorSensorChan)

	for {
		floor := <-floorSensorChan
		elevio.SetFloorIndicator(floor)
		channels.CurrentFloorFHM <- floor
	}
}

func pollHWORder() {

	btnChan := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(btnChan)

	for {
		btnValue := <-btnChan
		hwOrder := datatypes.Order{
			Floor:     btnValue.Floor,
			OrderType: int(btnValue.Button),
		}
		channels.OrderFHM <- hwOrder
	}
}

func updateOrderLights() {
	for {
		select {
		case orderComplete := <-channels.ClearLightsFOM:
			setAllLightsAtFloor(orderComplete.Floor, false)
		case orderRegistered := <-channels.SetLightsFOM:
			elevio.SetButtonLamp(elevio.ButtonType(orderRegistered.OrderType),
				orderRegistered.Floor, true)
		}
	}
}

func setAllLightsAtFloor(floor int, value bool) {
	for btn := datatypes.OrderUp; btn <= datatypes.OrderInside; btn++ {
		if !(floor == 0 && btn == datatypes.OrderDown) &&
			!(floor == numberOfFloors-1 && btn == datatypes.OrderUp) {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, value)
		}
	}
}
