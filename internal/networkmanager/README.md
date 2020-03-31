# networkmanager

## Overview
The networkmanager package handles all communication with the network driver `Network-go`
and it is thus the link between the network of elevators and our custom elevator 
software. The purpose of the network manager is to guarantee that information 
sent through it arrives at all elevators connected to the network in a 
reliable fashion. It also ensures that the measures to satisfy reliability does 
not put unnecessary strain on the network or other modules.

## Exported functions

### `NetworkManager`
* NetworkManager is the main routine-function for network communication.
Initializes all package-variables and starts all subroutines needed
for normal operation. Starts the network driver. Stays alive waiting to start 
transmitter and reciever again if the network status changes. Please note that 
as all managers in our software, the network manager communicates only via 
channels accessed through the global channels package implemented in this 
project software.
* Input argument(s): No input arguments
* Return value(s) : No return values

## Implementation
The networkmanager is the middle point between our custom local elevator softtware and our network of online elevators. It therefore implements transceiving infrastructure for all relevant datatypes, error handling and packet redundancy to ensure reliable communication between the elevators.

By only exposing the `NetworkManager` function, it allows the `ordermanager` package to transmit and recieve from ordinary golang channels defined in the `channels` package as if all elevators were connected directly through these. It thus effectively makes the complexity of transmitting with udp on a network invisible to the rest of the elevator software.

The error handling that the package accounts for is a loss/regain of network connection. If the network connection is lost, this is automatically detected by a internal listener routine, which instructs the `NetworkManager` to restart the driver in `localhost` mode. Similarly, if network connection is regained, we restart drivers again in the online mode. To account for the restaring time, all channels to and from `ordermanager` to the `networkmanager` are buffered to ensure normal operation during this restart period. It should be noted that the driver `Network-go` has been modified in a minor way to allow for the mode instruction and restarting scheme. 

With regards to packet redundancy, the package implements a uniform system for implementing this redundancy. Every packet that is to be sent on the network is sent 10 times over. To avoid the `ordermanager` package reciving 9 duplicates of every package, a signature scheme has been created which ensures a unique signature for every unique package. When a package arrives, the signature is checked against an overview of recently recievd signatures, and if it is deemed a duplicate it is discarded. 