//go:build debug
// +build debug

package main

import (
	"fmt"
	"net/http"
)

func ctlPongHandler(w http.ResponseWriter, r *http.Request) {

	// return message to HTTP request
	fmt.Fprint(w, "pong\n")
}

func init() {
	ctlMux.HandleFunc("GET /ctl/ping", ctlPongHandler)
}
