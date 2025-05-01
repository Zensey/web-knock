BINARY = web-knock

$(BINARY):
	go build -ldflags="-s -w" -o $(BINARY)
	sudo setcap 'cap_net_raw=+ep cap_net_admin=+ep' $(BINARY)

$(BINARY)-:
	chmod +x $(BINARY)
	sudo setcap 'cap_net_raw=+ep cap_net_admin=+ep' $(BINARY)

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)
