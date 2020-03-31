# External Watchdog

## Overview
This package has as its only responsibility to make sure that the elevator is always running. If the elvator software for any reason has crashed, the Watchdog detects this and restarts the elevator without crashing itself. 

## Exported functions

### `WatchdogNode`
* Configures the watchdog based on the spesification in the `configuration` package and starts timeouthandler subroutine for normal operation. This function runs as a standalone application. Has the responsibility of updating the timestamp to the value in the most recent watchdog-message from the `ElevatorNode` and restarting the elevator application if the duration since the last timestamp exceeds a configurable treshold.
* **Input argument(s):** 
* * `watchdogport string` The tcp port for communication between watchdog and elevator 
* * `elevport string` The tcp port for communication between elevator and elevator driver.
* **Output argument(s):**
* * No output arguments

### `ElevatorNode`
* Runs as a subroutine in the elevator applicaion. Sends timestamped messages to `WatchdogNode` to inform that the elevator application has not crashed.
* **Input argument(s):** 
* * `port string` The tcp port for communication between watchdog and elevator 
* **Output argument(s):**
* * No output arguments

## Implementation
The watchdog runs as its own standalone application, communicating with the `ElevatorNode` subroutine in the elevator application through a tcp connection on the localhost of the machine running the applications, where the subroutine is a client and the `WatchdogNode` is a server. Periodically, `ElevatorNode` sends timestamped messages and `WatchdogNode` updates its record of the latest timestamp. More frequently than messages are sent, `WatchdogNode` checks how much time has passed since the last recorded timestamp, and if more time has passed than the aforementioned treshold, then it will assume that the elevator application has crashed and start a new instance of the elevator application to ensure normal operation. When it does this, it informs the application through the elevator application flags which ports are used for communication with elevator driver and watchdog, and also informs the application of the process ID (PID) of the instance that crashed.
