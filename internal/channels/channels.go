package channels

import "github.com/sanderfu/TTK4145-ElevatorProject/internal/datatypes"

//SWOrderTX for transmitting to network via driver
var SWOrderTX chan datatypes.SWOrder = make(chan datatypes.SWOrder)

//SWOrderRX for recieveing from network via driver
var SWOrderRX chan datatypes.SWOrder = make(chan datatypes.SWOrder)

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

//SWOrderTOM channel from Network Manager to Order Manager
var SWOrderTOM chan datatypes.SWOrder = make(chan datatypes.SWOrder, 10)

//SWOrderFOM channel from Order Manager to Network Manager
var SWOrderFOM chan datatypes.SWOrder = make(chan datatypes.SWOrder)

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

//Order manager

//OrderFHM delivers orders from Hardware Manager to Order Manager
var OrderFHM chan datatypes.SWOrder = make(chan datatypes.SWOrder)

//Unfinished
var LigthCommand chan datatypes.LightCommand = make(chan datatypes.LightCommand)
