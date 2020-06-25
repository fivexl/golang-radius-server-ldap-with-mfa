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
}

// NewLdapConnection ...
func NewLdapConnection(config Config) (lc LdapConnection, err error) {
	lc = LdapConnection{addr: config.Ldap.Addr,
		userDn:      config.Ldap.UserDn,
		basedn:      config.Ldap.BaseDn,
		user:        config.Ldap.User,
		password:    config.Ldap.Password,
		groupFilter: config.Ldap.GroupFilter}
	return lc, nil
}

// Connect ...
func (l *LdapConnection) Connect() (conn *ldap.Conn, err error) {

	if conn, err = ldap.DialURL(l.addr); err != nil {
		return nil, fmt.Errorf("Failed to connect to LDAP server %s: %w", l.addr, err)
	}

	return conn, nil
}

// TestConnection - ...
func (l *LdapConnection) TestConnection() (err error) {

	log.Printf("Checking connection to LDAP server %s", l.addr)
	conn, err := l.Connect()
	if err != nil {
		return fmt.Errorf("Error while testing connection to %s: %w", l.addr, err)
	}
	defer conn.Close()

	// Only check bind if groupFilter is specified
	if l.groupFilter != "" {
		if err = conn.Bind(l.user, l.password); err != nil {
			return fmt.Errorf("Failed to bind to LDAP server %s using user %s: %w", l.addr, l.user, err)
		}
	}

	log.Printf("Successful connection to LDAP server")

	return nil
}

// CheckUser - binds to LDAP server using provided creds and returns auth status
func (l *LdapConnection) CheckUser(username string, password string, rLog *log.Logger) (result bool, err error) {

	conn, err := l.Connect()
	if err != nil {
		return false, fmt.Errorf("Error while checking user %s: %w", username, err)
	}
	defer conn.Close()

	userDn := strings.Replace(l.userDn, "{{username}}", username, -1)

	rLog.Printf("Checking LDAP auth for %s", userDn)
	if err = conn.Bind(userDn, password); err != nil {
		return false, err
	}

	return true, nil
}

// CheckGroupMembership - searches user
func (l *LdapConnection) CheckGroupMembership(username string, rLog *log.Logger) (result bool, err error) {

	conn, err := l.Connect()
	if err != nil {
		return false, fmt.Errorf("Error while checking user %s: %w", username, err)
	}
	defer conn.Close()

	// Bind using binder user so we can perform search
	if err = conn.Bind(l.user, l.password); err != nil {
		return false, err
	}

	searchFilter := strings.Replace(l.groupFilter, "{{username}}", username, -1)
	rLog.Printf("Checking user %s group membership using filter %s", username, searchFilter)

	searchRequest := ldap.NewSearchRequest(
		l.basedn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchFilter,
		[]string{},
		nil,
	)

	var searchResult *ldap.SearchResult
	if searchResult, err = conn.Search(searchRequest); err != nil {
		return false, fmt.Errorf("Error while searching for user %s using filter %s: %s", username, searchFilter, err)
	}

	if len(searchResult.Entries) > 1 {
		rLog.Printf("Found more than one record using filter %s", searchFilter)
		return false, nil
	}

	if len(searchResult.Entries) == 0 {
		rLog.Printf("Found no records using filter %s", searchFilter)
		return false, nil
	}

	rLog.Printf("Found one LDAP record that corresponds to provided filter. Consider check successful")

	return true, nil
}
