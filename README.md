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
