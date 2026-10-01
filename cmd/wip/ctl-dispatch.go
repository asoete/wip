//go:build debug
// +build debug

package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/watchdog"
)

// ----------------------------------------------------------------------------

func init() {
	ctlMux.HandleFunc("GET /ctl/watchdog/list", ctlAlarmWatchdogListAlarms)
}

func ctlAlarmWatchdogListAlarms(w http.ResponseWriter, r *http.Request) {

	printable := make(map[int64]any)
	for i, w := range watchdog.GetWatchers() {
		printable[i] = map[string]any{
			"Alarm":   w.Alarm,
			"Channel": w.Channel,
			"Countdown": map[string]any{
				"Timer":    "<time.AfterFunc>",
				"Deadline": w.Timer.Deadline,
			},
		}
	}
	json, err := json.MarshalIndent(printable, "", "  ")

	if err != nil {
		slog.Error("unable to Marshal watchdog.alarms", "error", err)
		os.Exit(1)
	}

	fmt.Fprint(w, string(json))
}
