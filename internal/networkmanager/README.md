# Network Manager

### TODO
 * Create channels for communicating with network driver, set up basic transmission - DONE
 * Create test for basic transmission - DONE
    * PASSED
 * Implement buffered channels - DONE AND REMOVED
 * Create test for buffered channels - DONE
    * Test showed that network driver has buffer, therefore no need for buffered channels. Reverting to non-buffered to save memory
 * Create channels for communication with Order Manager - DONE
 * Implement functionality for generating unique packet signature and maintaining list of recent signatures - DONE
 * Create test for unique signature functionality - DONE
    * PASSED
 * Implement redundancy in packets sent to combat packet loss with unique packet signature and multiple packet duplicates - DONE
 * Create test for packet loss transmission (comparing with nonredundant packets) - DONE
    * PASSED
 * Implement checking such that only unique packets are sent on channels to Order Manager - DONE
 * Create test for packet uniqueness - DONE
    * PASSED
    * Comment: Same test as above
 * Modify signature to be able to send more than one order every second - DONE
 * Verify that everything still working
    * PASSED
 * Test how the manager reacts to network loss - DONE
    * It is not detected, with base functionality of the driver everythign just stops working until network is back
 * Modify bcast package to allow for localhost broadcasting - DONE 
 * Implement checking for network connection loss and handle (restart Network Manager silently in localhost mode)
    * WORK IN PROGRESS, current state NOT WORKING
 * Test network connection loss handling