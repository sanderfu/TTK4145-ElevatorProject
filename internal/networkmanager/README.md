# Network Manager

## Introduction
The network manager package handels all communication with the network driver
and it is thus the link between the network and our custom elevator software.
The purpose of the network manager is to make packet loss and other errors regarding the network such as connection drop/regaining connection invisible to the rest of the system by silently restarting the driver in an offline/online state whenever such a change in network connection happens. To guarantee normal operation during the restart time, the channels the manager use to communicate with the order manager are buffered to allow the order-manager to keep on working during the restarting procedure of the network manager.

## Functions

### NetworkManager(**"channels"**)
NetworkManager is the main routine-function for network communication.
Initializes all package-variables and starts all subroutines needed
for normal operation. Stays alive waiting to start transmitter and reciever again if the necessary 