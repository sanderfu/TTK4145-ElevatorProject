# elevator

## Overview
The elevator application code. Reads the configuration from flags and loads 
additional configuration into memory from configuration JSON. Starts all 
applications needed for its operation in goroutines and then the main thread
suspends indefinitely as it has no more to do.

## Implementation
Parses flags and reads JSON to get the correct configuration and starts the 
`ElevatorNode` to communicate with the `watchdog`.

Then, the manager functions of the different packages in internals directory 
are invoked. It is worth noting that these managers handle the functionality of 
the respective package themselves and does not make any work go through the 
main application thread. The managers communicate thorugh the global channels 
infrastructure defined in the `channels` package designed for this project. The 
`channels` package holds the details of how this communication infrastructure 
is set up, and it will therefore not be mentioned further here.

Lastly, the Finite State Machine (FSM) of the elevator is invoked in its own 
gorutine and the main thread is thereafter suspended indefinitely as it after 
this point has no more work to do.
