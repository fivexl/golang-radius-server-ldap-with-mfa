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
