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
)

// Config ...
type Config struct {
	Radius struct {
		Listen string
		Secret string
	}
	Ldap struct {
		Addr        string
		UserDn      string
		BaseDn      string
		User        string
		Password    string
		GroupFilter string
	}
	Duo struct {
		Enabled bool
		APIHost string
		SKey    string
		IKey    string
		TimeOut int
	}
}

// DefaultTimeout - defines default time out for Duo client
const DefaultTimeout = 30

// Validate - makes sure that configuration is valid
// There should be better way
func (c *Config) Validate() (err error) {
	// Radius config
	log.Println("Validating RADIUS configuration")
	if c.Radius.Listen == "" {
		return fmt.Errorf("Failed to validate RADIUS configuration - Listen is required")
	}
	// Check that Listen is an address
	if c.Radius.Secret == "" {
		return fmt.Errorf("Failed to validate RADIUS configuration - Secret is required")
	}
	log.Println("RADIUS configuration is valid")
	log.Println("Validating LDAP configuration")
	if c.Ldap.Addr == "" {
		return fmt.Errorf("Failed to validate LDAP configuration - Addr is required")
	}
	// Check that Addr is an address
	if c.Ldap.UserDn == "" {
		return fmt.Errorf("Failed to validate RADIUS configuration - UserDn is required")
	}
	if c.Ldap.GroupFilter != "" {
		if c.Ldap.User == "" {
			return fmt.Errorf("Failed to validate LDAP configuration - GroupFilter specified but User is not. User is required when GroupFilter is not empty")
		}
		if c.Ldap.Password == "" {
			return fmt.Errorf("Failed to validate LDAP configuration - GroupFilter specified but Password is not. Password is required when GroupFilter is not empty")
		}
		if c.Ldap.BaseDn == "" {
			return fmt.Errorf("Failed to validate LDAP configuration - GroupFilter specified but BaseDn is not. BaseDn is required when GroupFilter is not empty")
		}
	}
	log.Println("LDAP configuration is valid")
	if c.Duo.Enabled {
		log.Println("DUO is enabled. Validating DUO configuration")
		// Check format
		if c.Duo.IKey == "" {
			return fmt.Errorf("Failed to validate DUO configuration - IKey is valid")
		}
		// Check format
		if c.Duo.SKey == "" {
			return fmt.Errorf("Failed to validate DUO configuration - SKey is required")
		}
		// Check that it is address
		if c.Duo.APIHost == "" {
			return fmt.Errorf("Failed to validate DUO configuration - APIHost is required")
		}
		if c.Duo.TimeOut == 0 {
			log.Printf("Duo client time out is not set in config. Fallback to default - %d sec", DefaultTimeout)
			c.Duo.TimeOut = DefaultTimeout
		}
	} else {
		log.Println("DUO is disabled. Skip DUO configuration validation")
	}
	return nil
}
