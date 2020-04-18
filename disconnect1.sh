#Simulator 1
sudo iptables -A INPUT -p tcp --dport 15657 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 15657 -j ACCEPT
#Watchdog 1
sudo iptables -A INPUT -p tcp --dport 15658 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 15658 -j ACCEPT

#Simulator 2
sudo iptables -A INPUT -p tcp --dport 15659 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 15659 -j ACCEPT
#Watchdog 2
sudo iptables -A INPUT -p tcp --dport 15660 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 15660 -j ACCEPT

#Localhost port for elevator 1 to talk to itself
sudo iptables -A INPUT -p udp --dport 15661 -j ACCEPT
sudo iptables -A INPUT -p udp --sport 15661 -j ACCEPT

#Localhost port for elevator 2 to talk to itself
sudo iptables -A INPUT -p udp --dport 15662 -j ACCEPT
sudo iptables -A INPUT -p udp --sport 15662 -j ACCEPT

#Drop all other network connections
sudo iptables -A INPUT -j DROP
