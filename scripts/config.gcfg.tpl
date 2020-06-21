[Radius]
Listen=127.0.0.1:1812
Secret={{ FivexL/github.com/golang-radius-server-ldap/radius-secret }}
[LDAP]
Addr={{ FivexL/github.com/golang-radius-server-ldap/ldap-addr }}
UserDN={{ FivexL/github.com/golang-radius-server-ldap/ldap-userdn }}
[DUO]
Enabled=true
IKey={{ FivexL/github.com/golang-radius-server-ldap/duo-ikey }}
SKey={{ FivexL/github.com/golang-radius-server-ldap/duo-skey }}
APIHost={{ FivexL/github.com/golang-radius-server-ldap/duo-apihost }}