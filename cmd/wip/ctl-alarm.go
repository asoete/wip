//go:build debug
// +build debug

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/watchdog"
)

// ----------------------------------------------------------------------------

func init() {
	ctlMux.HandleFunc("GET /ctl/watchdog/{id}/exists", ctlAlarmWatchdogHasAlarm)
}

func ctlAlarmWatchdogHasAlarm(w http.ResponseWriter, r *http.Request) {

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		log.Fatalf("ctlAlarmWatchdogHasAlarm(%s) failed: %s", alarm_id_str, err)
	}

	_, ok := watchdog.GetWatchdog(alarm_id)

	fmt.Fprint(w, strconv.FormatBool(ok))
}
