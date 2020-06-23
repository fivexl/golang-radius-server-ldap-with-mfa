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
	}
}
