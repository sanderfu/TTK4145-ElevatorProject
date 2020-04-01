# TTK4145-ElevatorProject

## Usage

This is a *Makefile* project which exists of three different programs, combining
to form the complete system. The usage of the Makefile is the following is 
the following

* `make` builds all programs needed to run the elevator
* `make elevator` builds just the elevator program
* `make watchdog` builds just the watchdog program
* `make startelevator` builds just the program starting the entire program
* `make run` starts the `startElevator` program, thus starting the entire 
system
* `make clean` deletes all executables in `build` folder and the entire `assets` 
folder.

The project includes three executables in `build`, namely `elevator`, `watchdog`
and `startElevator` where the functionality of each program is

* `elevator` is the program running the elevator
* `watchdog` a program for making sure that the `elevator` program is restarted
if it crashes
* `startElevator` a start up program which makes sure that both the `elevator`
and `watchdog` programs are started correctly

## Overview

This elevator project implements a modified mesh network where each node is 
the exact same. This gives the following advantages

* The code base will be identical for all nodes in the system
* The system is easily scalable for N number of elevators

The system does not operate with a shared world view, meaning everyone does not
know everything about the others. Rather the system is *auction based*, meaning
that once an order is received it is broadcasted to the other elevators and they
in turn respond with their cost value associated with that order. The two 
elevators (for redundancy) with the lowest cost are chosen to handle that order.
The one with the lowest cost will execute the order and the one with the second
lowest cost will be a backup elevator that takes the order in case the first 
elevator fails to serve it. This means that all elevators only know their own 
order queues, where each elevator has the following order queues

* **primary queue:** the queue with the orders that the elevator is executing
* **backup queue:** the queue with the orders that the elevator is ready to take
over if something should fail and they are not executed

## Repository structure

The structure of the repository is inspired by the Golang standard project 
layout, found [here](https://github.com/golang-standards/project-layout). The 
various folders in the project are used for the following

* `cmd` includes all main-programs for the various applications. In this case it
has the three programs `elevator`, `watchdog` and `startelevator`.
* `internal` here the custom modules needed for all the applications in `cmd` 
are stored.
* `vendor` the same as internal with the difference being that these are not 
self written, and handouts for the project. This includes the files for the 
elevator driver and network driver.
* `build` where all executables are stored
* `assets` created in run time and contains the order queue logs
* `docs` documentation

## Modules

This project consists of the following modules in the `internal` folder, further
documented in their own READMEs

* **Network manager:** handles broadcasting packages over UDP including packet 
loss. ([README](./internal/networkmanager/README.MD))
* **Order manager:** handles all order related actions including order auctions
and order queues. ([README](./internal/ordermanager/README.md))
* **Finite state machine (FSM):** Handles the state of the elevator.
([README](./internal/fsm/README.md))
* **Hardware manager:** handles all hardware related actions of the elevator.
([README](./internal/hwmanager/README.md))
* **Watchdog:** Makes sure the main elevator program does not stop running.
([README](./internal/watchdog/README.md))
* **Configuration:** Handles the initialization of system parameters upon 
startup. ([README](./internal/configuration/README.md))

The handed out drivers are the following and documentation can be found on
Github

* **Network-go:** Network driver, found on Github 
[here](https://github.com/TTK4145/Network-go)
* **driver-go:** Elevator driver, found on Github 
[here](https://github.com/TTK4145/driver-go) 

