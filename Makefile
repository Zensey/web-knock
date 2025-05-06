BINARY = web-knock
BINARY2 = enrichgeo

$(BINARY):
	go build -ldflags="-s -w" -o $(BINARY) ./cmd/$(BINARY)
	sudo setcap 'cap_net_raw=+ep cap_net_admin=+ep' $(BINARY)

$(BINARY2):
	go build -ldflags="-s -w" -o $(BINARY2) ./cmd/$(BINARY2)

clean:
	rm $(BINARY)

strip-$(BINARY): $(BINARY)
	strip -s $(BINARY)

getgeodb:
	wget https://git.io/GeoLite2-ASN.mmdb

all: getgeodb $(BINARY) $(BINARY2)
