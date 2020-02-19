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
 * Implement redundancy in packets sent to combat packet loss with unique packet signature and multiple packet duplicates
 * Create test for packet loss transmission
 * Implement checking such that only unique packets are sent on channels to Order Manager
 * Create test for packet uniqueness
