#!/bin/bash
if [$# -lt "1"]
then
echo hi
    build/startElevator
else
    build/startElevator -bcastlocalport $1
fi