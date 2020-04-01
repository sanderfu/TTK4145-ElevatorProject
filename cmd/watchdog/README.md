# watchdog

## Overview
The application code of the `watchdog` package. This application starts a watchdog node with configuartion defined from incoming flags.

## Implementation
The application firstly parses flags that were set when the application was called and translate these flags into the relevant input arguments for the `WatchdogNode`. The application then starts the `WatchdogNode` with these parameters and the main thread then has nothing more to do. For a detailed description of the `WatchdogNode`, see the relevant documentation in the internals directory. Lastly the main thread is suspended indefinitely to avoid terminating the `WatchdogNode` routine.
