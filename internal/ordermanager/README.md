# Order manager

### TODO
 * Set up datatypes(?) & channels for comm. with HW Manager and Queue Manager
 * Make function to recieve HW orders and send SW Order
    * Assumptions:
        * Networkmanager must write ID to structs requiring this
    * Test reciving costAns and signing Primary and Backup ID - PASSED
    * Test reciving costAns alone on network - PASSED
    * Test working cost signature - PASSED