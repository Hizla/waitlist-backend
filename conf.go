package main

import "os"

const (
	dbPath uint8 = iota
	listenAddr
	allowedOrigins

	confLen
)

// env variable, default pairing
var confEnv = [confLen][2]string{
	{"DB", "db"},
	{"LISTEN_ADDR", "127.0.0.1:3000"},
	{"ALLOWED_ORIGINS", "https://hizla.io"},
}

// resolved config values
var conf [confLen]string

func init() {
	for i := 0; i < int(confLen); i++ {
		if v, ok := os.LookupEnv(confEnv[i][0]); !ok {
			conf[i] = confEnv[i][1]
		} else {
			conf[i] = v
		}
	}
}
