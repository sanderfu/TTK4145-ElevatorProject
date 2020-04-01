# Hardware Manager

## Overview
The hardware manager handles all interacting with the hardware of the elevator.
This includes interfacing with the `elevio` driver and using this to expose the
necessary parameters and functionality of the elevator to the rest of the 
system in a meaningful way.

## Exported functions

* `HardwareManager`
    * Main function for starting the hardware manager. Starts all subroutines
    associated with the hardware manager
    * **Input argument(s):** None
    * **Retrun value(s):** None
* `SetElevatorDirection`
    * Sets the motor direction of the elevator
    * **Input arguments(s):** `dir int` Direction (-1 = down, 0 = stop, 1 = up)
    * **Return value(s):** None
* `SetDoorOpenLamp`
    * Sets the door open lamp
    * **Input argument(s):** `value bool` Set lamp on or off
    * **Return argument(s):** None

## Implementation

The hardware manager is a middle point between the hardware and the rest of the 
system. This means that the manager takes the raw hardware data from the 
`elevio` driver and converts it to the different data types used by the rest of 
the system. 

It exposes the `SetElevatorDirection` and `SetDoorOpenLamp` functions in 
`elevio` with wrapper functions. This is to totally hide the `elevio` driver to 
the rest of the system. 

The hardware manager also polls the elevator buttons and sensors in order to 
watch for new order and relay these to the order manager and to provide the rest
of the system with information on the current floor of the elevator. 

In the opposite way the hardware manager listens for registered orders or 
completed order from the order manager and manages the order lights accordingly.
