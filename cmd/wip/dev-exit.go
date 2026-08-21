//go:build debug
// +build debug

package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func exitHandler(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")
	if code == "" {
		code = "128"
	}

	codei, err := strconv.Atoi(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 Internal Server Error")
		log.Fatalf("exitHandler(): invalid exit code specified: %v. strconv.Atoi: %v\n", code, err)
	}

	delay := 100 * time.Millisecond

	// log message to console
	slog.Warn("[CTL] service got shutdown command", "exitcode", codei, "delay", delay)

	// return message to HTTP request
	fmt.Fprint(w, "Shutting down in 100ms\n")

	// Allow request to return: async exit
	go func() {
		time.Sleep(delay)
		slog.Warn("[CTL] shutting down NOW")
		os.Exit(codei)
	}()
}

func init() {
	ctlMux.HandleFunc("POST /ctl/exit/{code}", exitHandler)
}
