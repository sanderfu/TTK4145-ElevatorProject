package ordermanager

import (
	"math"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

const (
	c1 = 1
	c2 = 1
)

func genCostAnswer(costReq datatypes.CostRequest) datatypes.CostAnswer {
	var costAns datatypes.CostAnswer
	costAns.DestinationID = costReq.SourceID
	var newDirection int

	channels.FloorAndDirectionRequestFOM <- struct{}{}
	var lastFloor int = <-channels.FloorFFSM
	var currentDirection int = <-channels.DirectionFFSM

	if lastFloor > costReq.Floor {
		newDirection = 0
	} else if lastFloor < costReq.Floor {
		newDirection = 1
	} else {
		newDirection = 2
	}

	var directionMatch int
	if currentDirection != newDirection {
		directionMatch = 1
	} else {
		directionMatch = 0
	}

	if costReq.OrderType == datatypes.OrderInside && costReq.SourceID != costReq.DestinationID {
		costAns.CostValue = maxCostValue + 1
	} else {

		costAns.CostValue = c1*int(math.Abs(float64(costReq.Floor-lastFloor))) + c2*(directionMatch)
	}
	return costAns
}

func costRequestListener() {
	var costReq datatypes.CostRequest
	for {
		costReq = <-channels.CostRequestFNM
		channels.CostAnswerFOM <- genCostAnswer(costReq)
	}
}
