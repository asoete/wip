package main

import (
	"context"
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

var APP config.Application

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

	err := db.BootRO(APP.Db.Dsn)
	if err != nil {
		slog.Error("unable to bootstrap RO database", "dsn", APP.Db.Dsn, "error", err)
	}

	err = db.BootRW(APP.Db.Dsn)
	if err != nil {
		slog.Error("unable to bootstrap RW database", "dsn", APP.Db.Dsn, "error", err)
	}

	dispatch.Boot()

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
