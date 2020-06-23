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
	addr        string `validate:"nonzero"`
	userDn      string `validate:"nonzero"`
	basedn      string
	user        string
	password    string
	groupFilter string
	conn        *ldap.Conn
}

// NewLdapConnection ...
func NewLdapConnection(config Config) (lc LdapConnection, err error) {
	lc = LdapConnection{addr: config.Ldap.Addr,
		userDn:      config.Ldap.UserDn,
		basedn:      config.Ldap.BaseDn,
		user:        config.Ldap.User,
		password:    config.Ldap.Password,
		groupFilter: config.Ldap.GroupFilter,
		conn:        nil}
	if err = validator.Validate(lc); err != nil {
		return LdapConnection{}, err
	}
	if lc.groupFilter != "" {
		if lc.user == "" {
			return LdapConnection{}, fmt.Errorf("Failed to validate LDAP configuration - GroupFilter specified but User is not. User is required when GroupFilter is not empty")
		}
		if lc.password == "" {
			return LdapConnection{}, fmt.Errorf("Failed to validate LDAP configuration - GroupFilter specified but Password is not. Password is required when GroupFilter is not empty")
		}
		if lc.basedn == "" {
			return LdapConnection{}, fmt.Errorf("Failed to validate LDAP configuration - GroupFilter specified but BaseDn is not. BaseDn is required when GroupFilter is not empty")
		}
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

	// Only check bind if groupFilter is specified
	if l.groupFilter != "" {
		var loginOk bool
		if loginOk, err = l.CheckUser(l.user, l.password); err != nil {
			return fmt.Errorf("Failed to bind to LDAP server %s using user %s: %w", l.addr, l.user, err)
		}

		if !loginOk {
			return fmt.Errorf("Unsuccessful bind using user %s. Please check username and pasword in configuration file", l.user)
		}
	}

	log.Printf("Successful connection to LDAP server\n")

	return err
}

// CheckUser - binds to LDAP server using provided creds and returns auth status
func (l *LdapConnection) CheckUser(username string, password string) (result bool, err error) {

	if l.conn == nil {
		log.Printf("LdapConnection.CheckUser invoked but LdapConnection.conn is not yet intialized. Connecting to LDAP server")
		if err = l.Connect(); err != nil {
			return false, err
		}
	}

	userDn := strings.Replace(l.userDn, "{{username}}", username, -1)

	log.Printf("Checking LDAP auth for %s", userDn)
	if err = l.conn.Bind(userDn, password); err != nil {
		return false, err
	}

	return true, nil
}

// CheckGroupMembership - searches user
func (l *LdapConnection) CheckGroupMembership(username string) (result bool, err error) {

	searchFilter := strings.Replace(l.groupFilter, "{{username}}", username, -1)
	log.Printf("Checking user %s group membership using filter %s", username, searchFilter)

	searchRequest := ldap.NewSearchRequest(
		l.basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{},
		nil,
	)

	var searchResult *ldap.SearchResult
	if searchResult, err = l.conn.Search(searchRequest); err != nil {
		return false, fmt.Errorf("Error while searching for user %s using filter %s: %s", username, searchFilter, err)
	}

	if len(searchResult.Entries) > 1 {
		log.Printf("Found more than one record using filter %s", searchFilter)
		return false, nil
	}

	if len(searchResult.Entries) == 0 {
		log.Printf("Found no records using filter %s", searchFilter)
		return false, nil
	}

	log.Printf("Found one LDAP record that corresponds to provided filter. Consider check successful")

	return true, nil
}
