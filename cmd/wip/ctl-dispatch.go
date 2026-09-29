//go:build debug
// +build debug

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/dispatch"
)

// ----------------------------------------------------------------------------

func init() {
	ctlMux.HandleFunc("GET /ctl/dispatch/list-alarms", ctlAlarmDispatchListAlarms)
}

func ctlAlarmDispatchListAlarms(w http.ResponseWriter, r *http.Request) {

	json, err := json.MarshalIndent(dispatch.GetAlarms(), "", "  ")

	if err != nil {
		slog.Error("unable to Marshal dispatch.alarms", "error", err)
		os.Exit(1)
	}

	fmt.Fprint(w, string(json))
}

// ----------------------------------------------------------------------------

func init() {
	ctlMux.HandleFunc("GET /ctl/dispatch/list-timers", ctlTimerDispatchListTimers)
}

func ctlTimerDispatchListTimers(w http.ResponseWriter, r *http.Request) {

	var list []int64

	for i, _ := range dispatch.GetTimers() {

		list = append(list, i)
	}

	json, err := json.MarshalIndent(list, "", "  ")

	if err != nil {
		slog.Error("unable to Marshal dispatch.alarms", "error", err)
		os.Exit(1)
	}

	fmt.Fprint(w, string(json))
}
