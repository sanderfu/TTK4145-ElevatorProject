# Feedback regarding design from studass

## 05.02.2020
### SW Order "Interpeter":
* Have a routine for master and one for backup instead of one for every order. Have a master queue and backup queue in these routines.
* Status: Implemented in design

### HW Order Interpeter:
* Does not need special case for inside orders, can just broadcast to all elevators (the others will give cost inf because it is an inside order naturally)
* Status: Implemented in design

### Overall design
* Take care to not make a huge control/FSM module but divide code as close to equally as possible between modules.
* Status: Noted