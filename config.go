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
