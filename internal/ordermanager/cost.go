package ordermanager

import (
	"math"

	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

// Constants to tune the equation for cost generation
const (
	weightFloor    = 1
	weightDirMatch = 1
)

func genCostAnswer(costReq datatypes.CostRequest) datatypes.CostAnswer {
	var costAns datatypes.CostAnswer
	var newDirection int
	var directionMatch int

	costAns.DestinationID = costReq.SourceID

	// requesting last floor and direction from FSM
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

	if currentDirection != newDirection {
		directionMatch = 1
	} else {
		directionMatch = 0
	}

	if costReq.OrderType == datatypes.OrderInside && costReq.SourceID != costReq.DestinationID {
		costAns.CostValue = maxCostValue + 1
	} else {
		costAns.CostValue = weightFloor*int(math.Abs(float64(costReq.Floor-lastFloor))) + weightDirMatch*directionMatch
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
