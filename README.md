# LDAP based RADIUS server with MFA Support
* Written in Golang
* LDAP based
* DUO support (Push notification)
* Okta support is coming
* Reading configuration from HashiCorp Vault is coming

# Installation
Copy binary to server and place config file next to it. That is it

# Configuration format and example

## Radius section

Contains parametes that allows configuration of RADIUS server connection
* `Listen` - address that sever should listen for requests. Recommended value - `127.0.0.1:1812`. Mandatory.
* `Secret` - password to connect to RADIUS server. Mandatory.

## LDAP section
* `Addr`        - address of the LDAP server. Mandatory.
* `UserDN`      - user DN template that will be used to bind to LDAP server. Use `{{username}}` as a placeholder where to fill in user name. Server will use this to validate user LDAP auth by rendering user name instead of `{{username}}` and then performing bind to the LDAP server. Mandatory.
* `User`        - binder user DN that will be used to run LDAP searches when `GroupFilter` is specified. Not used otherwise. Optional.
* `Password`    - bind user password that will be used to run LDAP searches when `GroupFilter` is specified. Not used otherwise. Optional.
* `BaseDN`      - used to scope down group membership checkes when `GroupFilter` is specified. Not used otherwise. Optional.
* `GroupFilter` - LDAP search query that could be used to check user group memebership. Or whatever you want to check. Use `{{username}}` as a placeholder where to fill in user name. Search considered successful if only one entry returned. Will not run any searches if not specified. Optional. Requires - `User`, `Password`, `BaseDN`

## DUO section
Optional section that enables MFA with DUO. Read [here](https://duo.com/docs/protecting-applications) how to get credentials
If you do not use DUO just omit this section.
Note! DUO MFA only works for users with mobile push-capable device configured, i.e. server will attempt to send push-based auth request via DUO
and if user will have no device configured then auth request will be automatically rejected.
* `Enabled` - should MFA check be performed
* `IKey`    - DUO integration key
* `SKey`    - DUO security key
* `APIHost` - DUO API host
* `TimeOut` - DUO client time out. Defines how long duo client will wait for user to react on push notification. Default to 30 sec

## Minmal config file example
Will only use LDAP auth
```
[Radius]
Listen=127.0.0.1:1812
Secret=secret
[LDAP]
Addr="ldaps://ldap.server.com:636"
UserDN="uid={{username}},ou=Users,o=xxxxxxxxxx,dc=xxxxx,dc=com"
```

## LDAP section config example for Jumpcloud LDAP with group membership check
```
[LDAP]
Addr=ldaps://ldap.jumpcloud.com:636
UserDN=uid={{username}},ou=Users,o=0aaa0a000aa00a0baaaa00aa,dc=jumpcloud,dc=com
BaseDN=ou=Users,o=0aaa0a000aa00a0baaaa00aa,dc=jumpcloud,dc=com
User=uid=binder,ou=Users,o=0aaa0a000aa00a0baaaa00aa,dc=jumpcloud,dc=com
Password=EK0ToHQ4NWJ9nFHLEK0ToHQ4NWJ9nFH!
GroupFilter=(&(objectClass=inetOrgPerson)(uid={{username}})(memberOf=cn=group-name,ou=Users,o=0aaa0a000aa00a0baaaa00aa,dc=jumpcloud,dc=com))
```

## DUO section config example
```
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

Checking packets sent to server
```
sudo tshark -f "udp port 1812" -i any -V
```
