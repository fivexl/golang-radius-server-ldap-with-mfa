package main

import (
	"fmt"
	"net/url"

	duoapi "github.com/duosecurity/duo_api_golang"
	duoauthapi "github.com/duosecurity/duo_api_golang/authapi"
)

// GetDuoAuthClient ...
func GetDuoAuthClient(iKey string, sKey string, apiHost string) (client *duoauthapi.AuthApi, err error) {
	duoClient := duoapi.NewDuoApi(
		iKey,
		sKey,
		apiHost,
		"",
	)
	client = duoauthapi.NewAuthApi(*duoClient)

	var check *duoauthapi.CheckResult
	if check, err = client.Check(); err != nil {
		return nil, err
	}
	if check == nil {
		return nil, fmt.Errorf("could not connect to Duo; got nil result back from API check call")
	}
	var msg, detail string
	if check.StatResult.Message != nil {
		msg = *check.StatResult.Message
	}
	if check.StatResult.Message_Detail != nil {
		detail = *check.StatResult.Message_Detail
	}
	if check.StatResult.Stat != "OK" {
		return nil, fmt.Errorf("Could not connect to Duo: %s (%s)", msg, detail)
	}
	return client, nil
}

// DuoPushAuth ...
func DuoPushAuth(client *duoauthapi.AuthApi, username string) (result bool, err error) {

	options := []func(*url.Values){duoauthapi.AuthUsername(username), duoauthapi.AuthDevice("auto")}

	var authResult *duoauthapi.AuthResult
	if authResult, err = client.Auth("push", options...); err != nil {
		return false, err
	}

	if authResult == nil {
		return false, fmt.Errorf("Duo Auth request returned no error but response is empty")
	}

	if authResult.StatResult.Stat != "OK" {
		return false, fmt.Errorf("Could not authenticate Duo user: %s (%s)", *authResult.StatResult.Message, *authResult.StatResult.Message_Detail)
	}

	if authResult.Response.Result != "allow" {
		return false, nil
	}

	return true, nil
}
