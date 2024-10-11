package main

import (
	"log"
	"os"
)

const (
	dbPath uint8 = iota
	listenAddr
	allowedOrigins
	hCaptchaSecret
	hCaptchaSiteKey
	verboseLogging
	confLen
)

// env variable, default pairing
var confEnv = [confLen][2]string{
	{"DB", "db"},
	{"LISTEN_ADDR", "127.0.0.1:3000"},
	{"ALLOWED_ORIGINS", "https://hizla.io"},
	{"HCAPTCHA_SECRET", "unset"},
	{"HCAPTCHA_SITE_KEY", "unset"},
	{"VERBOSE", "1"},
}

// resolved config values
var conf [confLen]string

var verbose bool

func init() {
	for i := 0; i < int(confLen); i++ {
		if v, ok := os.LookupEnv(confEnv[i][0]); !ok {
			conf[i] = confEnv[i][1]
		} else {
			conf[i] = v
		}
	}

	switch conf[verboseLogging] {
	case "0":
		verbose = false
	case "1":
		verbose = true
	default:
		log.Printf("invalid verbose value %q", conf[verboseLogging])
	}
}
