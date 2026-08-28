package main

import (
	"context"
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

var dumpConfig = false

// types
type tokenMux struct {
	handler http.Handler
}

func (m *tokenMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	recToken := r.Header.Get(tokenMuxTokenHeader)
	if recToken == tokenMuxToken {
		m.handler.ServeHTTP(w, r)
		return
	}

	abort(w, http.StatusUnauthorized, "invalid token (host=%s ; path=%s ; user=%s)", r.URL.Host, r.URL.Path, r.URL.User)
}

func init() {
	// Web service
	flag.StringVar(&listenAddress, "web.address", "127.0.0.1:8080", "bind to this address")

	// Database
	flag.StringVar(&dbDSN, "db.dsn", "sqlite::memory:", "use this database")

	// Control service
	flag.StringVar(&ctlAddress, "ctl.address", "127.0.0.1:8100", " service control (ctl) endpoints")

	// Misc options
	flag.BoolVar(&dumpConfig, "dump-config", false, "dump the active config and exit")
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

	if dumpConfig || slog.Default().Handler().Enabled(context.Background(), slog.LevelDebug) {

		if !dumpConfig {
			fmt.Printf("== DEBUG: dump config =====================================\n")
		}

		fmt.Printf("    --web.address = %s\n", listenAddress)

		fmt.Printf("         --db.dsn = %s\n", dbDSN)

		fmt.Printf("    --ctl.address = %s\n", ctlAddress)
		fmt.Printf("  X-WiP-Ctl-Token = %s\n", tokenMuxToken)

		if !dumpConfig {
			fmt.Printf("===========================================================\n")
		}

		if dumpConfig {
			os.Exit(0)
		}
	}

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
