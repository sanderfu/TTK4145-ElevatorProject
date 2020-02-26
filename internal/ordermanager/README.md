# Order manager

### TODO
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
    * Test interaction between three PCs where one recieves HW order and all have different cost
 * Test that most extensive test from network manager still produces correct result