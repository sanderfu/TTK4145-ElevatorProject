package main

import (
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/watchdog"
)

func main() {
	go watchdog.WatchdogNode()

	for {
		time.Sleep(50 * time.Second)
	}
}
