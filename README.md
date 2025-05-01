## web-knock: An alternative to port-knocking

### ABOUT  

This is a web-knock server.
Instead of some special port sequence it uses another method -- web-knocking.

Lets' imagine that we want to secure SSH port 22 on our server, so that no one can
connect to SSH server w/o knocking.

In order to open the access you send a special web request to our host.
When the program discovers the valid web-request in the access.log of the web server
it opens the port by adding origin's IP to the whitelist ipset of firewall.



### Build binary

    `make web-knock`

## Setup firewall rules (ufw)

At the end of /etc/ufw/before.rules find the last COMMIT line and paste the following before it:

    -A ufw-before-input -p tcp -m tcp --dport 8443 -m set --match-set web_blacklist src -j REJECT
    -A ufw-before-input -p tcp -m tcp --dport 22 -m set --match-set ssh_whitelist src -j ACCEPT
    -A ufw-before-input -p tcp --dport 22 -j REJECT


## Run the daemon

    `./web-knock -key=secretcode`

## Knocking your server

    https://yourserver/secretcode
