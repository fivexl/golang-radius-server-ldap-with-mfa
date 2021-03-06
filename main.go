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
	"os"
	"time"

	duoauthapi "github.com/duosecurity/duo_api_golang/authapi"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"

	"gopkg.in/gcfg.v1"

	"github.com/urfave/cli/v2"
)

// VERSION will be passed via compliation flag
var VERSION string

// AuthRequest - encapsulates approval logic
func AuthRequest(username string, password string, rLog *log.Logger, lc LdapConnection, dc *duoauthapi.AuthApi) (soFarSoGood bool, err error) {
	soFarSoGood = false
	if lc.groupFilter != "" {
		if soFarSoGood, err = lc.CheckGroupMembership(username, rLog); err != nil {
			return false, fmt.Errorf("LDAP group membership check for user %s failed with error: %w", username, err)
		}
		rLog.Printf("LDAP group membership check status for user %s: %t", username, soFarSoGood)
		if !soFarSoGood {
			rLog.Printf("Reject auth request from user %s", username)
			return false, nil
		}
	}
	if soFarSoGood, err = lc.CheckUser(username, password, rLog); err != nil {
		return false, fmt.Errorf("LDAP auth for user %s failed with error: %w", username, err)
	}
	rLog.Printf("LDAP auth status for user %s: %t", username, soFarSoGood)
	if !soFarSoGood {
		rLog.Printf("Reject auth request from user %s", username)
		return false, nil
	}
	if dc != nil {
		if soFarSoGood, err = DuoPushAuth(dc, username); err != nil {
			return false, fmt.Errorf("DUO auth for user %s failed with error: %w", username, err)
		}
		rLog.Printf("DUO auth status for user %s: %t", username, soFarSoGood)
		if !soFarSoGood {
			rLog.Printf("Reject auth request from user %s", username)
			return false, nil
		}
	}
	rLog.Printf("Final auth status for user %s: %t", username, soFarSoGood)
	return soFarSoGood, nil
}

func main() {
	app := &cli.App{
		Name:      "LDAP based RADIUS server with MFA support",
		Usage:     "Provide a valid config file and server will take care of the rest",
		UsageText: "https://github.com/FivexL/golang-radius-server-ldap-with-mfa/README.md",
		Version:   VERSION,
		Action:    Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "./config.gcfg",
				Usage:   "Path to the config file",
				EnvVars: []string{"RADIUS_SERVER_CONFIG_PATH"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Run - function to be called by cli.App wrapper
func Run(c *cli.Context) (err error) {

	log.Printf("LDAP based RADIUS server with MFA support version %s", VERSION)
	log.Printf("Documentation and installation instructions: https://github.com/FivexL/golang-radius-server-ldap-with-mfa/blob/master/README.md")

	var configPath string = c.String("config")
	var config Config
	log.Printf("Parsing config file %s", configPath)
	if err = gcfg.ReadFileInto(&config, configPath); err != nil {
		return fmt.Errorf("Failed to read config file %s: %w", configPath, err)
	}
	if err = config.Validate(); err != nil {
		return fmt.Errorf("Config validation failed: %w", err)
	}

	var lc LdapConnection
	if lc, err = NewLdapConnection(config); err != nil {
		return fmt.Errorf("Failed to configure LDAP server connection: %w", err)
	}

	if err := lc.TestConnection(); err != nil {
		return fmt.Errorf("Failed to connect to LDAP server due to: %w", err)
	}

	var dc *duoauthapi.AuthApi = nil
	if config.Duo.Enabled {
		log.Printf("DUO MFA is enabled. Initiating DUO client for API endpoint %s", config.Duo.APIHost)
		if dc, err = GetDuoAuthClient(config.Duo.IKey, config.Duo.SKey, config.Duo.APIHost, config.Duo.TimeOut); err != nil {
			return fmt.Errorf("Failed to initiate Duo client: %w", err)
		}
	}

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)
		requestID := time.Now().UnixNano() / int64(time.Millisecond)
		requestLogger := log.New(os.Stdout, fmt.Sprintf("%s [Request-%d]", log.Prefix(), requestID), log.LstdFlags|log.Lmsgprefix)
		authResult := false
		code := radius.CodeAccessReject

		requestLogger.Println("Handling new request")

		if authResult, err = AuthRequest(username, password, requestLogger, lc, dc); err != nil {
			requestLogger.Printf("Error while performing auth for user %s: %s", username, err)
		} else if authResult {
			code = radius.CodeAccessAccept
		}
		requestLogger.Printf("Writing %v to %v", code, r.RemoteAddr)
		err = w.Write(r.Response(code))
		if err != nil {
			requestLogger.Printf("Encountered error when responding to request: %s", err)
		}
	}

	server := radius.PacketServer{
		Addr:         config.Radius.Listen,
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(config.Radius.Secret)),
	}

	log.Printf("Starting server on %s", config.Radius.Listen)
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
