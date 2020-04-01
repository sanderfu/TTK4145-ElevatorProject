# Finite state machine (FSM)

## Overview

The FSM keeps track of and controls the state of the elevator. To do this, it 
uses the queue from the order manager to get the next objective to go to, while
using the hardware manager to control the hardware of the elevator.

## Exported functions

* `FSM`
    * Starts the FSM and all needed to run it. Keeps track of the state of the
    elevator and chooses which state function to run
    * **Input argument(s):** None
    * **Return value(s):** None

## Implementation

The FSM is implemented by each state having its own function associated with it.
A global variable with the current state is the parameter used to determine what
state function to call next and this solution was chosen due to its simplicity
after finding that a complicated solution was not needed in this case.

The three states included in the FSM are

* **idle:** The state when the elevator has no orders in its primary queue and 
thus no orders to attend. Here it is simply checking if a new order has arrived
in the queue and if so it will take the elevator to the moving state.
* **moving:** The moving state represents the state where the elevator is 
serving an order, but has not yet arrived at the destination floor. If the 
elevator reaches the destination floor or passes by a floor with and order it 
can take on the way, the moving state will take the elevator to the door open 
state.
* **door open:** The door open state is when the elevator has stoppet at a floor
with an order and is servicing that order. This state will simply wait for a 
fixed amount of time (given in the config file) and then close the door and 
inform the order manager that the order was completed. It will then return the 
elevator to the idle state.

The FSM module also keeps track of the last floor the elevator was at and the 
current direction it is heading. This information is sent to the order manager
when requested as parameters in the cost function.