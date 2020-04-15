package channels

import "github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"

// Format for channels:
// var <channelName>F<sourceModule>T<destinationModule>
// F = from, T = to
// Modules:
// 		nm  - Network Manager
//		om  - Order Manager
// 		fsm - FSM
// 		hm  - Hardware Manager

////////////////////////////////////////////////////////////////////////////////
// Network Manager Channels
////////////////////////////////////////////////////////////////////////////////

// Internal signalling channels in Network Manager
var KillTransmitter = make(chan struct{}, 1)
var KillReceiver = make(chan struct{}, 1)
var InitTransmitter = make(chan struct{}, 1)
var InitReceiver = make(chan struct{}, 1)

var SWOrderPrimaryFnmTom = make(chan datatypes.Order, 10)
var SWOrderBackupFnmTom = make(chan datatypes.Order, 10)
var CostRequestFnmTom = make(chan datatypes.CostRequest, 10)
var CostAnswerFnmTom = make(chan datatypes.CostAnswer, 10)
var OrderRecvAckFnmTom = make(chan datatypes.OrderRecvAck, 10)
var OrderCompleteFnmTom = make(chan datatypes.OrderComplete, 10)
var OrderRegisteredFnmTom = make(chan datatypes.OrderRegistered, 10)

////////////////////////////////////////////////////////////////////////////////
// Order manager channels
////////////////////////////////////////////////////////////////////////////////

var SWOrderFomTnm = make(chan datatypes.Order, 10)
var CostRequestFomTnm = make(chan datatypes.CostRequest, 10)
var CostAnswerFomTnm = make(chan datatypes.CostAnswer, 10)
var OrderRecvAckFomTnm = make(chan datatypes.OrderRecvAck, 10)
var OrderCompleteFomTnm = make(chan datatypes.OrderComplete, 10)
var OrderRegisteredFomTnm = make(chan datatypes.OrderRegistered, 10)

// Signals for Network driver from Order Manager
var KillDriverTX = make(chan struct{}, 1)
var KillDriverRX = make(chan struct{}, 1)
var InitDriverTX = make(chan struct{}, 1)
var InitDriverRX = make(chan struct{}, 1)

// Internal channels for queue modifications
var PrimaryQueueAppend = make(chan datatypes.QueueOrder, 10)
var PrimaryQueueRemove = make(chan datatypes.QueueOrder, 10)
var BackupQueueAppend = make(chan datatypes.QueueOrder, 10)
var BackupQueueRemove = make(chan datatypes.QueueOrder, 10)

var FloorAndDirectionRequestFomTfsm = make(chan struct{}, 1)

////////////////////////////////////////////////////////////////////////////////
// Hardware manager channels
////////////////////////////////////////////////////////////////////////////////

var OrderFhmTom = make(chan datatypes.Order, 1)
var CurrentFloorFhmTfsm = make(chan int, 1)
var HMInitStatusFhmTfsm = make(chan bool, 1)
var ClearLightsFomThm = make(chan datatypes.OrderComplete, 1)
var SetLightsFomThm = make(chan datatypes.OrderRegistered, 1)

////////////////////////////////////////////////////////////////////////////////
// FSM channels
////////////////////////////////////////////////////////////////////////////////

var OrderCompleteFfsmTom = make(chan datatypes.OrderComplete, 1)
var FloorFfsmTom = make(chan int, 1)
var DirectionFfsmTom = make(chan int, 1)

////////////////////////////////////////////////////////////////////////////////
// Network driver channels
////////////////////////////////////////////////////////////////////////////////

// All between Network Manager and Network driver
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
