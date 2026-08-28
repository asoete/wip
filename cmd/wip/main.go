package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// flags
var listenAddress string
var listenPort int

var dbDSN string

var ctlAddress string
var ctlMux = http.NewServeMux()

var tokenMuxTokenHeader = "X-WiP-Ctl-Token"
var tokenMuxToken = "wip-ctl-token-please-change!"

type tokenMux struct {
	handler http.Handler
}

func (m *tokenMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	recToken := r.Header.Get(tokenMuxTokenHeader)
	if recToken == tokenMuxToken {
		m.handler.ServeHTTP(w, r)
		return
	}

	abort(w, http.StatusUnauthorized, "invalid token")
}

func init() {
	// Web service
	flag.StringVar(&listenAddress, "web.address", "127.0.0.1:8080", "bind to this address")

	// Database
	flag.StringVar(&dbDSN, "db.dsn", "sqlite::memory:", "use this database")

	// Control service
	flag.StringVar(&ctlAddress, "ctl.address", "127.0.0.1:8100", " service control (ctl) endpoints")

	// Control flags
	flag.StringVar(&ctlAddress, "ctl-address", "127.0.0.1:8100", " Service control (ctl) endpoints")
}

func main() {

	currentLogLevel := slog.SetLogLoggerLevel(slog.LevelDebug)
	defer slog.SetLogLoggerLevel(currentLogLevel) // revert changes after the example

	// fetch config from ENV
	if os.Getenv("WIP_CTL_TOKEN") != "" {
		tokenMuxToken = os.Getenv("WIP_CTL_TOKEN")
	}

	// Handle command line parameters
	flag.Parse()

	// Start (separate) server to listen for control commands
	go func() {
		slog.Info("[CTL] start control service", "--ctl.address", ctlAddress)
		slog.Debug("[CTL]", "X-WiP-Ctl-Token", tokenMuxToken)
		log.Fatal(http.ListenAndServe(ctlAddress, &tokenMux{ctlMux}))
	}()

	// Start main webserver
	slog.Info("[WEB] start web service", "--web.address", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
