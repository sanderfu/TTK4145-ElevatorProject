sudo iptables -A INPUT -p tcp --dport 15657 -j ACCEPT
sudo iptables -A INPUT -p tcp --sport 15657 -j ACCEPT

sudo iptables -A INPUT -j DROP
