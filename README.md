# TTK4145-ElevatorProject

## Notes
* Might need to know which elevators is online. Not sure yet. If this becomes 
only a convenience, we will drop it
* Orders are divided into HW and SW orders. HW orders are orders from the 
physical elevator buttons. SW orders are commands to execute HW orders sent 
between elevators

## Modules

All modules are the exact same on every elevator.

### Order Interpreter

* Orders are divided into HW and SW orders. HW orders are orders from the 
physical elevator buttons. SW orders are commands to execute HW orders sent 
between elevators
* All orders will be sent to both a primary and backup
* If only one elevator on network, this elevator becomes primary and backup
* The backup will ensure that the order is handled, either by the primary or if
the primary fails it handles the order itself. 
* This can lead to two elevators handling the same order, but this is okay, as
we focus on no order losses rather than performance. 

#### HW Order Registration 

##### Description
This module listens for hardware orders from the physical layer and decides 
who should execute these orders. From this, SW Orders are generated with master 
and backup details are broadcasted on the network.

##### Notes
* Lights should only be lit when the elevator knows it can execute the order. 
This is indicated by the HW order OK status

#### SW Order Registration
This module listens for SW orders sent on the network and gives message to Queue
 Manager to registrate them in to the correct local queue on the elevator 
 machine.

##### Notes
* 

#### Cost Calculation
This module listens for cost requests on the network. If it receives a request 
it calculates its cost and sends it back to the requestor. 

##### Notes
* How the cost function calculates this and what dependencies it needs must be
specified.
* It must request its own elevator's order queue and state
* Inside orders must have infinite cost if elevator id does not match elevator 
id in order.


### Queue manager

#### Notes

Abilities
* Add orders coming from order manager to queue
* Delete completed orders from hardware
* Check for expired backup orders and move these to primary queue if expired
* Save queue locally
* Load queue from local file
* Send message to order manager when primary order complete
* Receive message from order manager that backup order was completed on another
elevator
* Be able to calculate cost of an order

### Hardware manager

#### Notes
Abilities:
* Initialize hardware
* Listen for hardware orders 
* Receive control commands from FSM (go up, go down, stay)
* Send HW Orders to Order Manager
* Turn on or off lamps when receiving messages from order manager
* Tell FSM the state of the elevator (going up, going down, reached floor x)

### FSM/Control

#### Notes
Abilities:
* Intialize the elevator to a known state (e.g drive down to know floor etc.)
* Initialize the other managers and receive status from these
* Keep track of the elevator's state (going up, going down, etc.)
* Ask queue manager for first next objective
* Ask queue manager if we can stop for other orders
* Provide queue manager with the current state of the elevator

##### Startup note

FSM -> Hardware init -> elev init
Hardware send status ok to FSM
FSM -> queue init
Queue -> order init
Order -> Network module init
Network init send status ok to order
Order init send status ok to queue 
Queue init send status ok to FSM 

If status not ok, reset.

### Network module

#### Notes

Technical Implementation:
  * Topology: Mesh network
  * All messages are broadcasted via UDP
  * Must choose blocking or select
  * Core reliability built into module
  * Messages are structs and are packed/unpacked in the network module.
  * Detection and handling of lost packages are dealt with in the network module
  * Lost nodes are not handeled as messages are broadcasted

Guarantees about the elevator:
  * Other nodes do not care, the faulty node will try to reconnect and stay on 
localhost in meantime (will be invisible for order manager)


### Elevator driver
Given in https://github.com/TTK4145/driver-go. 


