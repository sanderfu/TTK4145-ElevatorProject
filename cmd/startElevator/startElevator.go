package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
)

const (
	connHost = ":"
	connType = "tcp"
)

func findOpenPort() (int, net.Listener) {
	connPort := 16698
	addr := connHost
	addr += strconv.Itoa(connPort)
	fmt.Println(addr)
	l, err := net.Listen("tcp", addr)

	for err != nil {
		fmt.Printf("Port %v already in use, increments..\n", connPort)
		connPort++
		addr = connHost + strconv.Itoa(connPort)
		fmt.Println(addr)
		l, err = net.Listen("tcp", addr)
	}
	return connPort, l
}

func choosePorts() (string, string) {
	wport, lw := findOpenPort()
	eport, le := findOpenPort()
	defer lw.Close()
	defer le.Close()
	watchdogport := strconv.Itoa(wport)
	elevport := strconv.Itoa(eport)
	return watchdogport, elevport
}

func main() {

	watchdogport, elevport := choosePorts()
	fmt.Printf("Watchdogport: %v\n elevport: %v\n", watchdogport, elevport)

	cmdWatchdog := exec.Command("gnome-terminal", "-e", "build/watchdog -watchdogport "+watchdogport+" -elevport "+elevport)
	cmdWatchdog.Run()

	cmdElevatorHardware := exec.Command("gnome-terminal", "-e", "./SimElevatorServer --port "+elevport)
	cmdElevatorHardware.Run()

	cmdElevatorSoftware := exec.Command("gnome-terminal", "-e", "build/elevator -elevport "+elevport+" -watchdogport "+watchdogport)
	cmdElevatorSoftware.Run()

}
