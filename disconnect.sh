# Drop broadcast port
sudo iptables -A INPUT -p udp --dport 16569 -m statistic --mode random --probability 1 -j DROP

# Drop tcp port for finding local IP
sudo iptables -A INPUT -p tcp --dport 53 -m statistic --mode random --probability 1 -j DROP
sudo iptables -A INPUT -p tcp --sport 53 -m statistic --mode random --probability 1 -j DROP
