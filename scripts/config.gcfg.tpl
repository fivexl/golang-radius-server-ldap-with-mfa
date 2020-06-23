[Radius]
Listen=127.0.0.1:1812
Secret={{ FivexL/github.com/golang-radius-server-ldap/radius-secret }}
[LDAP]
Addr={{ FivexL/github.com/golang-radius-server-ldap/ldap-addr }}
UserDN={{ FivexL/github.com/golang-radius-server-ldap/ldap-userdn }}
BaseDN={{ FivexL/github.com/golang-radius-server-ldap/ldap-basedn }}
User={{ FivexL/github.com/golang-radius-server-ldap/ldap-user }}
Password={{ FivexL/github.com/golang-radius-server-ldap/ldap-password }}
GroupFilter={{ FivexL/github.com/golang-radius-server-ldap/ldap-group-filter }}
[DUO]
Enabled=true
IKey={{ FivexL/github.com/golang-radius-server-ldap/duo-ikey }}
SKey={{ FivexL/github.com/golang-radius-server-ldap/duo-skey }}
APIHost={{ FivexL/github.com/golang-radius-server-ldap/duo-apihost }}