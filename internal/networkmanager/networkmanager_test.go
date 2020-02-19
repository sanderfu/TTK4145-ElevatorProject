package networkmanager_test

import (
	"fmt"
	"time"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/networkmanager"
)

func NetworkManagerTestSending() {
	for {
		var testOrdre datatypes.SW_Order
		testOrdre.Primary_id = "12345"
		testOrdre.Backup_id = "67890"
		testOrdre.Dir = datatypes.INSIDE
		testOrdre.Floor = datatypes.FIRST
		networkmanager.SWOrderTX <- testOrdre
		time.Sleep(1 * time.Second)
	}
}

func NetworkManagerTestRecieving() {
	for {
		select {
		case order := <-networkmanager.SWOrderRX:
			fmt.Printf("Received: %#v\n", order)

		}
	}
}
