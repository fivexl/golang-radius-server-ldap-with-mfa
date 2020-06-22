# RADIUS server based on LDAP with out-of-the-box support for MFA
* Written in Golang
* LDAP based
* DUO support (Push notification)
* Okta support is coming
* Reading configuration from HashiCorp Vault is coming

# Installation
Copy binary to server and place config file next to it. That is it

# Configuration example
```
[Radius]
Listen=127.0.0.1:1812
Secret=secret
[LDAP]
Addr="ldaps://ldap.jumpcloud.com:636"
UserDN="uid={{username}},ou=Users,o=xxxxxxxxxx,dc=jumpcloud,dc=com"
[DUO]
Enabled=true
IKey=XXXXXXXXXXXX
SKey=XXXXXXXXXXXXXXX
APIHost=api-xxxxx.duosecurity.com
```
# Development

# Testing

First get server up and running
```
# Build the server
bash scripts/build.sh

# Render configuration (requires secretshub-cli)
sudo apt-get install -y curl gnupg2
curl -fsSL https://apt.secrethub.io/pub | sudo apt-key add -
echo "deb https://apt.secrethub.io stable main" | sudo tee -a /etc/apt/sources.list.d/secrethub.sources.list
sudo apt-get update
sudo apt-get install -y secrethub-cli
bash scripts/render-config.sh

# Run it
./build/rserver-linux-amd64 -c ./build/config.gcfg
 ```

For the test we need a client
```
 sudo apt-get install freeradius-utils
 radtest <username> <password> localhost 1812 <radius secret>
 ```