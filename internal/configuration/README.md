# Configuration

## Overview
This package is configurating the constants needed for the system to run and define and parse flags to be used in the system.

## Exported functions
* `ReadConfig`
    * Reads the config file and decodes its constants into global variables for the whole system to use.
    * **Input argument(s):** `string` name of the file 
    * **Retrun value(s):** None

* `ParseFlags`
    * Initialize flags which is then parsed and loaded into global variables.
    * **Input argument(s):** None
    * **Retrun value(s):** None

## Implementation
All the constants for the system is stored in the file `config.json`, this package is reading that file and updating the struct `Configuration`, with the correct constants. The flags is generated in the function `ParseFlags` which is then loaded into the struct `CommandLineFlags`. These two structs is accessible for all the packages in the system.