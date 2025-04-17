BINARY = web-knock

$(BINARY):
	go build -ldflags="-s -w" -o $(BINARY)
	sudo setcap 'cap_net_raw=+ep cap_net_admin=+ep' $(BINARY)

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

setup-fw-rules:
	sudo iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT
	sudo iptables -A INPUT -p tcp -m tcp --dport 22 -m set --match-set ssh_whitelist src -j ACCEPT
	sudo iptables -A INPUT -p tcp --dport 22 -j REJECT
