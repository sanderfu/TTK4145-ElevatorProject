package channels

import "github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"

////////////////////////////////////////////////////////////////////////////////
// Network driver channels
////////////////////////////////////////////////////////////////////////////////

var SWOrderTX = make(chan datatypes.Order, 1)
var SWOrderRX = make(chan datatypes.Order, 1)
var CostRequestTX = make(chan datatypes.CostRequest, 1)
var CostRequestRX = make(chan datatypes.CostRequest, 1)
var CostAnswerTX = make(chan datatypes.CostAnswer, 1)
var CostAnswerRX = make(chan datatypes.CostAnswer, 1)
var OrderRecvAckTX = make(chan datatypes.OrderRecvAck, 1)
var OrderRecvAckRX = make(chan datatypes.OrderRecvAck, 1)
var OrderCompleteTX = make(chan datatypes.OrderComplete, 1)
var OrderCompleteRX = make(chan datatypes.OrderComplete, 1)
var OrderRegisteredTX = make(chan datatypes.OrderRegistered, 1)
var OrderRegisteredRX = make(chan datatypes.OrderRegistered, 1)

////////////////////////////////////////////////////////////////////////////////
// Order manager channels
////////////////////////////////////////////////////////////////////////////////

var SWOrderFNMPrimary = make(chan datatypes.Order, 1)
var SWOrderFNMBackup = make(chan datatypes.Order, 1)
var SWOrderFOM = make(chan datatypes.Order, 1)
var CostRequestFNM = make(chan datatypes.CostRequest, 10)
var CostRequestFOM = make(chan datatypes.CostRequest, 1)
var CostAnswerFNM = make(chan datatypes.CostAnswer, 10)
var CostAnswerFOM = make(chan datatypes.CostAnswer, 1)
var OrderRecvAckFNM = make(chan datatypes.OrderRecvAck, 10)
var OrderRecvAckFOM = make(chan datatypes.OrderRecvAck, 1)
var OrderCompleteFNM = make(chan datatypes.OrderComplete, 10)
var OrderCompleteFOM = make(chan datatypes.OrderComplete, 1)
var OrderRegisteredFOM = make(chan datatypes.OrderRegistered, 10)
var OrderRegisteredFNM = make(chan datatypes.OrderRegistered, 10)

var KillDriverTX = make(chan struct{}, 1)
var KillDriverRX = make(chan struct{}, 1)
var InitDriverTX = make(chan struct{}, 1)
var InitDriverRX = make(chan struct{}, 1)

var PrimaryQueueAppend = make(chan datatypes.QueueOrder, 1)
var PrimaryQueueRemove = make(chan datatypes.QueueOrder, 1)
var BackupQueueAppend = make(chan datatypes.QueueOrder, 1)
var BackupQueueRemove = make(chan datatypes.QueueOrder, 1)

////////////////////////////////////////////////////////////////////////////////
// Hardware manager channels
////////////////////////////////////////////////////////////////////////////////

var OrderFHM = make(chan datatypes.Order, 10)
var CurrentFloorFHM = make(chan int, 1)
var HMInitStatusFHM = make(chan bool, 1)
var ClearLightsFOM = make(chan datatypes.OrderComplete, 10)
var SetLightsFOM = make(chan datatypes.OrderRegistered, 10)

////////////////////////////////////////////////////////////////////////////////
// FSM channels
////////////////////////////////////////////////////////////////////////////////

var OrderCompleteFFSM = make(chan datatypes.OrderComplete, 10)
