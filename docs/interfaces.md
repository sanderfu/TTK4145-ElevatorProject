# Modules interfaces

## Description

Contains the interfaces between the modules, including number of go channels and
data types sent in them.

## Network Manager - Order Manager

All types go both ways between the modules and every type has it's own go 
channel both ways.

### Elements
- Cost requests
- Cost answers
- SW orders
- Order received acknowledge
- Primary order complete broadcast

---

### Cost Request

A struct containing neccesary information to calculate the cost. 

Struct `cost_request`:
- source_id
- floor
- direction

#### Comments
`source_id` is the IP address of the elevator requesting the cost.

---

### Cost Answer

A struct containing the answer to a calculated cost.

Struct `cost_answer`:
- source_id
- cost_value

#### Comments
`source_id` is the IP address of the elevator requesting the cost.

---

### SW Order

A struct containing a software order

Struct `sw_order`:
- primary_id
- backup_id
- floor
- direction

---

### Order Received Acknowledge

A struct containing an acknowledge that and order is received.

Struct `order_recv_ack`:
- source_id
- floor
- direction

#### Comments
`source_id` is the IP address of the elevator who received the order.

---

### Order Complete

A struct containing the message that an order has been completed.

Struct `order_complete`:
- floor
- direction

---
---

## Order Manager - Queue Manager

### Elements
- SW Orders
- Order Complete
- Cost Value

### From Order Manager to Queue Manager

#### SW Order

A struct containing a software order

Struct `sw_order`:
- primary_id
- backup_id
- floor
- direction

---

### From Queue Manager to Order Manager

#### Cost Value

A struct containing the cost value for a given order. May be overkill with a 
struct for only one data type, but this is done to generalize.

Struct `cost_value`:
- cost_value

---

### Between both modules

### Order Complete

A struct containing the message that an order has been completed.

Struct `order_complete`:
- floor
- direction

---
---

## Queue Manager - FSM

### Elements
- Current Objective
- Current State

### From Queue Manager to FSM

#### Current Objective

A struct containing the next objective (floor) for the FSM

FINISH THIS LATER

