# Channels

## Overview
This package is a collection of all the channels that are used between packages,
and some inside packages. They are all on the format:
```go
<channelName>F<sourceModule>T<destinationModule>
```
where F means "from" and T means "to". This is done to easier see which modules
the channel connects.
