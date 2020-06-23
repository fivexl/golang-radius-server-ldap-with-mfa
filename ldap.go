/*

Copyright 2020 Andrey Devyatkin.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package main

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/ldap.v3"
)

// LdapConnection ...
type LdapConnection struct {
	addr        string
	userDn      string
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
