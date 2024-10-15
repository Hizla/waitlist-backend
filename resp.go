package main

import "encoding/json"

// pre-generate constant responses
func mustConstResp(a any) string {
	if p, err := json.Marshal(a); err != nil {
		panic(err.Error())
	} else {
		return string(p)
	}
}

type respM struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func newMessage(success bool, message string) *respM {
	return &respM{success, message}
}
