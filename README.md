# TTK4145-ElevatorProject

## Notes
* Might need to know which elevators is online. Not sure yet. If this becomes 
only a convenience, we will drop it
* Orders are divided into HW and SW orders. HW orders are orders from the 
physical elevator buttons. SW orders are commands to execute HW orders sent 
between elevators

## Modules

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
who should execute these orders. From this, SW Orders are generated with master and backup details are broadcasted on the network.

##### Notes
* Lights should only be lit when the elevator knows it can execute the order. 
This is indicated by the HW order OK status

#### SW Order Registration
This module listens for SW orders sent on the network and gives message to Queue Manager to registrate them in to the correct local queue on the elevator machine.

##### Notes
* 

#### Cost Calculation
This module listens for cost requests on the network. If it receives a request 
it calculates its cost and sends it back to the requestor. 

##### Notes
* How the cost function calculates this and what dependencies it needs must be
specified.
* It must request its own elevator's order queue and state
* Inside orders must have infinite cost if elevator id does not match elevator id in order.


### Queue manager

#### Notes

It needs to be able to
* Add orders coming from order manager to queue
* Delete completed orders from hardware
* Check for expired backup orders and move these to primary queue if expired
* Save queue locally
* Load queue from local file
* Send message to order manager when primary order complete
* Receive message from order manager that backup order was completed on another
elevator
* Be able to calculate cost of an order