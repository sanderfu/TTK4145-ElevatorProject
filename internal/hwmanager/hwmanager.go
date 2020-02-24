package hwmanager

import (
	"github.com/TTK4145/driver-go/elevio"
)

func Init() {
	addr := "192.168.0.163:15657"
	numFloors := 4

	elevio.Init(addr, numFloors)

	testPolling()
}

func testPolling() {

	print("Polling stop button... ")

	stopBtnChan := make(chan bool)
	go elevio.PollStopButton(stopBtnChan)

	pressed := <-stopBtnChan

	println("Pressed = ", pressed)

}
