package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/go-validator/validator.v2"
	"gopkg.in/ldap.v3"
)

// LdapConnection ...
type LdapConnection struct {
	addr   string `validate:"nonzero"`
	userDn string `validate:"nonzero"`
	conn   *ldap.Conn
}

// NewLdapConnection ...
func NewLdapConnection(addr string, userDn string) (lc LdapConnection, err error) {
	lc = LdapConnection{addr: addr, userDn: userDn, conn: nil}
	if err = validator.Validate(lc); err != nil {
		return LdapConnection{}, err
	}
	return lc, nil
}

// Connect ...
func (l *LdapConnection) Connect() (err error) {

	if l.conn != nil {
		log.Printf("LdapConnection.Connection invocked but LdapConnection.conn is not nil, i.e. already connected. Do nothing")
		return nil
	}

	log.Printf("Attempting to connect to LDAP server %s\n", l.addr)
	if l.conn, err = ldap.DialURL(l.addr); err != nil {
		return fmt.Errorf("Failed to connect to LDAP server %s: %w", l.addr, err)
	}

	log.Printf("Successful connection to LDAP server\n")
	return nil
}

// CheckUser ...
func (l *LdapConnection) CheckUser(username string, password string) (result bool, err error) {

	if l.conn == nil {
		log.Printf("LdapConnection.CheckUser invoked but LdapConnection.conn is not yet intialized. Connecting to LDAP server")
		if err = l.Connect(); err != nil {
			return false, err
		}
	}

	userDn := strings.Replace(l.userDn, "{{username}}", username, -1)

	log.Printf(userDn)
	if err = l.conn.Bind(userDn, password); err != nil {
		return false, err
	}

	return true, nil
}
