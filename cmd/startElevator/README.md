# startElevator

## Overview
The startElevator is a small standalone go-executable script that finds available tcp localhost ports for communication between the elevator application and the watcchdog and for communication between between the elevator application and the elevator hardware interface application(currently configured to strat the driver simulator). Moreover, this application starts the watchdog application, the elevator hardware interface application (simulator currently) and elevator application to run the elevator service on the current node.

## Implementation
The ports for localhost tcp communication are found by trying to set up a tcp listener at a pre-configured base port and incrementing the port number every time that fails due to the port being busy. After two available ports are found, these are stored and the listeners that have been set up when searching for them are closed so that the ports will in fact be free when the Watchdog application and the elevator hardware interface applications want to use them. 

Information is passed from the `startElevator` application to the other applications by the use of flags.

When the three applications has been opened, the `startElevator` application terminates as it has no more relevant tasks to do.