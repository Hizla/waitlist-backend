package main

import (
	"log"
	"os"
)

const (
	dbPath uint8 = iota
	listenAddr
	allowedURL
	hCaptchaSiteKey
	hCaptchaSecret
	verboseLogging
)

// env variable, default pairing
var confEnv = [...][2]string{
	dbPath:          {"DB", "db"},
	listenAddr:      {"LISTEN_ADDR", "127.0.0.1:3000"},
	allowedURL:      {"ALLOWED_URL", "https://hizla.io"},
	hCaptchaSiteKey: {"HCAPTCHA_SITE_KEY", "unset"},
	hCaptchaSecret:  {"HCAPTCHA_SECRET", "unset"},
	verboseLogging:  {"VERBOSE", "1"},
}

// resolved config values
var conf [len(confEnv)]string

var verbose bool

func init() {
	for i := 0; i < len(confEnv); i++ {
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
