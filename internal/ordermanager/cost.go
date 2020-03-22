package ordermanager

import (
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/channels"
	"github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"
)

func genCostAns(costReq datatypes.CostRequest) datatypes.CostAnswer {
	var costAns datatypes.CostAnswer
	costAns.DestinationID = costReq.SourceID
	if costReq.OrderType == datatypes.INSIDE && costReq.SourceID != costReq.DestinationID {
		costAns.CostValue = maxCost + 1
	} else {
		costAns.CostValue = 2*len(primaryQueue) + 1*len(backupQueue)
	}

	return costAns
}

func costReqWatch() {
	var costReq datatypes.CostRequest
	var costAns datatypes.CostAnswer
	for {
		costReq = <-channels.CostRequestTOM
		costAns.CostValue = 2*len(primaryQueue) + 1*len(backupQueue)
		costAns.DestinationID = costReq.SourceID
		channels.CostAnswerFOM <- costAns
	}
}
