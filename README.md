## web-knock: An alternative to port-knocking

### ABOUT  

This is a web-knock server.
Instead of some special port sequence it uses another method -- web-knocking.

Lets' imagine that we want to secure SSH port 22 on our server, so that no one can
connect to SSH server w/o knocking.

In order to open the access you send a special web request to our host.
When the program discovers the valid web-request in the access.log of the web server
it opens the port by adding origin's IP to the whitelist ipset of firewall.



### Install

    `make web-knock`
    `make setup-fw-rules`

## Run the daemon

    `./web-knock -key=secretcode`

## Knocking your server

    https://yourserver/secretcode
