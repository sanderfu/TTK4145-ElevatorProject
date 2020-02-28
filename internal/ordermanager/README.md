# Order manager

### TODO - HW Order Registration
 * Set up datatypes(?) & channels for comm. with HW Manager and Queue Manager
 * Make function to recieve HW orders, broadcast cost request, recieve cost answers and decide primary and backup
    * Assumptions:
        * Networkmanager must write ID to structs requiring this - IMPLEMENTED
    * Test reciving costAns and signing Primary and Backup ID - PASSED
    * Test reciving costAns alone on network - PASSED
    * Test working cost signature - PASSED AND REMOVED
    * Test interaction between two PCs where one recieves HW order - PASSED
 * Modify prev. function to wait for confirmation from primary and backup that order is recieved.
    * Assumptions:
        * Networkmanager must only let SWOrders where we are primary or backup through to Order Manager - IMPLEMENTED
    * Test interaction between three PCs where one recieves HW order and all have different cost - FAIL
        * Wrong Primary/Secondary combination is choosen - FIXED AND PASSED
    * Major bug was present with TOM channels filling up with no emptying. Fixed by introducing DestionationID.
 * Test that most extensive test from network manager still produces correct result
 * Implement "broadcast lightcommand"
    * Make this feature unnecessary by using the recieved SW order instead. Idea: Send a duplicate of the SW order on a dedicated channel from networkmanager to ordermanagerfor this.
 * Reduce the waittimes to break after 250ms instead of after 1s and test functionality on 3 Terminals
 * Repeat the same test on 3 pcs on the sanntidssal

### TODO - SW Order Registration
 * Set up routine to listen for incoming SW orders
 * Redirect these WS orders to primary or backup channel for registration in queuemanager
 * Set up waitloop