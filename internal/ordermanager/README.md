# Order Manager

## Overview
The order manager is parted in three files, whivch handles orders, queues and calculates cost for taking an order. It keeps track of both the primary and the backup queue, and is the only package with accsess to the queues. This means that it is interfacing between the `FSM`, the `network manager` and getting order request from the `hardware manager`.

## Exported functions
* `OrderManager`
    * Main function of order manager, which starts all subrutines associated with this package
    * Input argument(s): None
    * Return argument(s): None

* `OrderInQueue`
    * Checks if the order is in the primary queue
    * Input argument(s): `order`
    * Return argument(s): `bool`, True if the order is in the queue

* `GetFirstOrderInQueue`
    * Returns the first order in the primary queue
    * Input argument(s): None
    * Return argument(s): `order`

* `QueueEmpty`
    * Checks if the primary queue is empty
    * Input argument(s): None
    * Return argument(s): `bool`, True if the queue is empty

## Implementation
The order manager is handling all logic for the local elevator to function. It preps order request from the hardware manager for the network manager, so network manager only is interfacing with order manager. The package is sliced into three files, `order.go`, `queue.go` and `cost.go`, so not one file is doing several different computations.

### order.go
This is the main file for the package, which handles all new orders (and order requests) coming either from the network manager or the hardware manager. It is always listning for new orders and requestfrom the network manager and handles them respectivly.

### queue.go
This file is only working with the primary- and backup queue, adding/deleting orders from the queues. In addition to also saving the queue in a json file and loading the file if the system would crash and need to reboot.

### cost.go
When called upon in `order.go`, it generates a cost for taken a given order at the place and state the elevator is in at the moment. The cost is then sent back to `order.go`.