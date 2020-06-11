package main

// Config ...
type Config struct {
	Radius struct {
		Listen string
		Secret string
	}
	Ldap struct {
		Addr   string
		UserDn string
	}
	Duo struct {
		Enabled bool
		APIHost string
		SKey    string
		IKey    string
	}
}
