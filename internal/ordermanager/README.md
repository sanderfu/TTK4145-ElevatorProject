# Order manager

### Note:
The decision has been made to merge ordermanager and queuemanager under the ordermanager name. Docs will be updated after this module is completed.

### TODO - HW Order Registration
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
 * Reduce the waittimes to break after 250ms instead of after 1s and test functionality on 3 Terminals - PASSED
 * Repeat the same test on 3 pcs on the sanntidssal - PASSED
 * Repeat same test on 3 PCs ont he sanntidssal w. packet loss 20% - PASSED

### TODO - SW Order Registration
 * Set up routine to listen for incoming SW orders
 * Redirect these WS orders to primary or backup channel for registration in queuemanager - Deligated to Network Manager
 * Create queuestructs and arrays. Keep only necessary info in these structs - DONE
 * Test logging queues - DONE 
 
 ### TODO - Order complete watchers
 * Utilize channels and help functions to avoid race conditions when modifying queues
 * Change generateSignature in networkmanager to delay by 1 ms to guarantee unique signatures

 ### TODO - Major midway testing
 * Test sending orders and ordercompletes on one terminal on one computer - PASSED
 * Test sending orders and ordercompletes on 3 terminals on one computer - PASSED
 * Test sending order and ordercompletes on 3 PCs on the sanntidlab with 20% packet loss - PASSED

 