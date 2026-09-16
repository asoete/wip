//go:build debug
// +build debug

package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/dispatch"
)

// ----------------------------------------------------------------------------

func init() {
	ctlMux.HandleFunc("GET /ctl/alarm/{id}/dispatch-has-alarm", ctlAlarmDispatchHasAlarm)
}

func ctlAlarmDispatchHasAlarm(w http.ResponseWriter, r *http.Request) {

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		log.Fatalf("/ctl/alarm/{id:%s}/dispatch-has-alarm : parsing alarm_id failed: %s", alarm_id_str, err)
	}

	_, ok := dispatch.GetAlarm(alarm_id)

	fmt.Fprint(w, strconv.FormatBool(ok))
}

// ----------------------------------------------------------------------------

func init() {
	ctlMux.HandleFunc("GET /ctl/alarm/{id}/dispatch-has-timer", ctlAlarmDispatchHasTimer)
}

func ctlAlarmDispatchHasTimer(w http.ResponseWriter, r *http.Request) {

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		log.Fatalf("/ctl/alarm/{id:%s}/dispatch-has-timer : parsing alarm_id failed: %s", alarm_id_str, err)
	}

	_, ok := dispatch.GetTimer(alarm_id)

	fmt.Fprint(w, strconv.FormatBool(ok))
}
