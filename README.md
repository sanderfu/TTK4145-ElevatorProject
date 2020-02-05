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
* All orders will be sent to both a primary and backup (unless there is only one
elevator on the network)
* The backup will ensure that the order is handled, either by the primary or if
the primary fails it handles the order itself. 
* This can lead to two elevators handling the same order, but this is okay, as
we focus on no order losses rather than performance. 

#### HW Order Interpreter

##### Description
This module listens for hardware orders from the physical buttons and decides 
who should execute these orders.

##### Notes
* Lights should only be lit when the elevator knows it can execute the order. 
This is indicated by the HW order OK status
* Cab calls must always be saved locally

#### SW Order Interpreter
This module listens for SW orders from other elevators. And handles the 
difference between primary and backup routines.

##### Notes
* **It is important that these routines are not blocking / that a new routine is
created for each order**

#### Cost Calculation
This module listens for cost requests on the network. If it receives a request 
it calculates its cost and sends it back to the requestor. 

##### Notes
* How the cost function calculates this and what dependencies it needs must be
specified.
* It must somehow request its own elevator's order queue and state