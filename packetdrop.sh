#!/bin/bash
echo Packet drop rate: $1
sudo iptables -A INPUT -p udp --dport 16569 -m statistic --mode random --probability $1 -j DROP
