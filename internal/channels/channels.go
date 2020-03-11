package channels

import "github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"

//SWOrderTX for transmitting to network via driver
var SWOrderTX chan datatypes.Order = make(chan datatypes.Order)

//SWOrderRX for recieveing from network via driver
var SWOrderRX chan datatypes.Order = make(chan datatypes.Order)

//CostRequestTX for transmitting to network via driver
var CostRequestTX chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostRequestRX for recieveing from network via driver
var CostRequestRX chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostAnswerTX for transmitting to network via driver
var CostAnswerTX chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//CostAnswerRX for recieveing from network via driver
var CostAnswerRX chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//OrderRecvAckTX ...
var OrderRecvAckTX chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderRecvAckRX ...
var OrderRecvAckRX chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderCompleteTX ...
var OrderCompleteTX chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

//OrderCompleteRX ...
var OrderCompleteRX chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

//SWOrderTOMPrimary channel for delivering primary orders to Order Manager from Network Manager
var SWOrderTOMPrimary chan datatypes.Order = make(chan datatypes.Order)

//SWOrderTOMBackup is channel for delivering backup orders to Order Manager from Network manager
var SWOrderTOMBackup chan datatypes.Order = make(chan datatypes.Order)

//SWOrderFOM channel from Order Manager to Network Manager
var SWOrderFOM chan datatypes.Order = make(chan datatypes.Order)

//CostRequestTOM ...
var CostRequestTOM chan datatypes.CostRequest = make(chan datatypes.CostRequest, 10)

//CostRequestFOM ...
var CostRequestFOM chan datatypes.CostRequest = make(chan datatypes.CostRequest)

//CostAnswerTOM ...
var CostAnswerTOM chan datatypes.CostAnswer = make(chan datatypes.CostAnswer, 10)

//CostAnswerFOM ...
var CostAnswerFOM chan datatypes.CostAnswer = make(chan datatypes.CostAnswer)

//OrderRecvAckTOM ...
var OrderRecvAckTOM chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck, 10)

//OrderRecvAckFOM ...
var OrderRecvAckFOM chan datatypes.OrderRecvAck = make(chan datatypes.OrderRecvAck)

//OrderCompleteTOM ...
var OrderCompleteTOM chan datatypes.OrderComplete = make(chan datatypes.OrderComplete, 10)

//OrderCompleteFOM ...
var OrderCompleteFOM chan datatypes.OrderComplete = make(chan datatypes.OrderComplete)

var KillDriverTX = make(chan struct{}, 1)
var KillDriverRX = make(chan struct{}, 1)
var InitDriverTX = make(chan struct{}, 1)
var InitDriverRX = make(chan struct{}, 1)

//Hardware manager

//OrderFHM delivers orders from Hardware Manager to Order Manager
var OrderFHM chan datatypes.Order = make(chan datatypes.Order, 10)

var CurrentFloorTFSM chan int = make(chan int)

// Channel to send init status to FSM from HM
var HMInitStatusTFSM chan bool = make(chan bool)

//OrderCompleteTHM delivers order complete messages from Order Manager to Hardware manager to clear lights
var OrderCompleteTHM chan datatypes.OrderComplete = make(chan datatypes.OrderComplete, 10)

//OrderRegisteredTHM delivers orders from Order Manager to Hardware Manager that has been registered by Primary & Backup in non-volatile memory
var OrderRegisteredTHM chan datatypes.Order = make(chan datatypes.Order, 10)
