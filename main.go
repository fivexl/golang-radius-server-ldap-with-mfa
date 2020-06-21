package main

import (
	"fmt"
	"log"

	duoauthapi "github.com/duosecurity/duo_api_golang/authapi"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"

	"gopkg.in/gcfg.v1"
)

// AuthRequest ...
func AuthRequest(username string, password string, lc LdapConnection, dc *duoauthapi.AuthApi) (result bool, err error) {
	isLdapCheckOk := false
	isDuoCheckOk := true
	if isLdapCheckOk, err = lc.CheckUser(username, password); err != nil {
		return false, fmt.Errorf("LDAP auth failed for user %s: %w", username, err)
	}
	log.Printf("LDAP auth for user %s: %t", username, isLdapCheckOk)
	if dc != nil {
		if isDuoCheckOk, err = DuoPushAuth(dc, username); err != nil {
			return false, fmt.Errorf("DUO auth failed for user %s: %w", username, err)
		}
		log.Printf("DUO auth for user %s: %t", username, isDuoCheckOk)
	}
	result = isLdapCheckOk && isDuoCheckOk
	log.Printf("Overal status for use %s: %t", username, result)
	return result, nil
}

func main() {
	var configPath string = "./config.gcfg"
	var config Config
	var err error
	if err = gcfg.ReadFileInto(&config, configPath); err != nil {
		log.Fatalf("Failed to read config file %s: %s", configPath, err)
	}

	var lc LdapConnection
	if lc, err = NewLdapConnection(config.Ldap.Addr, config.Ldap.UserDn); err != nil {
		log.Fatalf("Failed to initiate connection to LDAP server: %s", err)
	}

	var dc *duoauthapi.AuthApi = nil
	if config.Duo.Enabled {
		if dc, err = GetDuoAuthClient(config.Duo.IKey, config.Duo.SKey, config.Duo.APIHost); err != nil {
			log.Fatalf("Failed to initiate Duo client: %s", err)
		}
	}

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		username := rfc2865.UserName_GetString(r.Packet)
		password := rfc2865.UserPassword_GetString(r.Packet)
		authResult := false
		code := radius.CodeAccessReject

		if authResult, err = AuthRequest(username, password, lc, dc); err != nil {
			log.Printf("Error while performing auth for user %s: %s", username, err)
		} else if authResult {
			code = radius.CodeAccessAccept
		}
		log.Printf("Writing %v to %v", code, r.RemoteAddr)
		w.Write(r.Response(code))
	}

	server := radius.PacketServer{
		Addr:         config.Radius.Listen,
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(config.Radius.Secret)),
	}

	log.Printf("Starting server on %s", config.Radius.Listen)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
