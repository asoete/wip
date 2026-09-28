package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "modernc.org/sqlite"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/abort"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/config"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/dispatch"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/event"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

var ctlMux = http.NewServeMux()
var tokenMuxToken string

// database handles
var dbRO *db.Queries
var dbRW *db.Queries

var APP config.Application

// types
type tokenMux struct {
	handler http.Handler
}

func (m *tokenMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	recToken := r.Header.Get(APP.Ctl.AuthHeader)
	if recToken == tokenMuxToken {
		m.handler.ServeHTTP(w, r)
		return
	}

	abort.Fatal(w, http.StatusUnauthorized, "invalid auth token",
		APP.Ctl.AuthHeader, recToken,
		"middleware", "tokenMux",
	)
}

func init() {
	APP.Init()
}

func main() {

	APP.Boot()

	currentLogLevel := slog.SetLogLoggerLevel(slog.LevelDebug)
	defer slog.SetLogLoggerLevel(currentLogLevel) // revert changes after the example

	// fetch config from ENV
	if os.Getenv("WIP_CTL_TOKEN") != "" {
		tokenMuxToken = os.Getenv("WIP_CTL_TOKEN")
	}

	if APP.DumpConfig || slog.Default().Handler().Enabled(context.Background(), slog.LevelDebug) {

		APP.Dump()

		if APP.DumpConfig {
			os.Exit(0)
		}
	}

	if APP.Pidfile != "" {

		// TODO: add plumbing to handle os.Exit, interupts and sigterms from other workers
		defer os.Remove(APP.Pidfile)

		err := os.MkdirAll(filepath.Dir(APP.Pidfile), os.ModePerm)
		if err != nil {
			log.Fatal("create pidfile (parent) dir failed:", err)
		}

		err = os.WriteFile(APP.Pidfile, []byte(strconv.Itoa(os.Getpid())), 0644)
		if err != nil {
			log.Fatal("create pidfile failed:", err)
		}
	}

	// init database handles
	dbhandle, err := sql.Open("sqlite", APP.Db.Dsn)
	if err != nil {
		log.Fatalf("unable to parse database(%s): %w", APP.Db.Dsn, err)
	}

	dbRW = db.New(dbhandle)

	dispatch.Boot(dbRW)

	event.SetDefault(event.New(dbRW))
	go event.Log(event.AppStart, user.System, nil)

	// Start (separate) server to listen for control commands
	go func() {
		slog.Info("[CTL] start control service", "--ctl.address", APP.Ctl.Address)
		slog.Debug("[CTL]", "X-WiP-Ctl-Token", tokenMuxToken)
		slog.Error("[CTL] server stopped", "error", http.ListenAndServe(APP.Ctl.Address, &tokenMux{ctlMux}))
	}()

	// Start main webserver
	slog.Info("[WEB] start web service", "--web.address", APP.Web.Address)
	slog.Error("[WEB] server stopped", "error", http.ListenAndServe(APP.Web.Address, nil))

	go event.Log(event.AppStop, user.System, nil)
}
