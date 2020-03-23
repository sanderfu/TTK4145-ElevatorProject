package ordermanager

import (
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

func genCostAnswer(costReq datatypes.CostRequest) datatypes.CostAnswer {
	var costAns datatypes.CostAnswer
	costAns.DestinationID = costReq.SourceID
	if costReq.OrderType == datatypes.OrderInside && costReq.SourceID != costReq.DestinationID {
		costAns.CostValue = maxCostValue + 1
	} else {
		costAns.CostValue = 2*len(primaryQueue) + 1*len(backupQueue)
	}

	return costAns
}

func costRequestListener() {
	var costReq datatypes.CostRequest
	var costAns datatypes.CostAnswer
	for {
		costReq = <-channels.CostRequestFNM
		costAns.CostValue = 2*len(primaryQueue) + 1*len(backupQueue)
		costAns.DestinationID = costReq.SourceID
		channels.CostAnswerFOM <- costAns
	}
}
