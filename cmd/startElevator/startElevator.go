package main

import "os/exec"

func main() {
	cmdWatchdog := exec.Command("gnome-terminal", "-e", "./watchdog")
	cmdElevatorHardware := exec.Command("gnome-terminal", "-e", "./SimElevatorServer")
	cmdElevatorSoftware := exec.Command("gnome-terminal", "-e", "./main")

	cmdWatchdog.Run()
	cmdElevatorHardware.Run()
	cmdElevatorSoftware.Run()
}
