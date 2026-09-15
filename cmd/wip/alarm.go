package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/dispatch"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

func init() {
	http.HandleFunc("POST /alarm", alarmCreateHandler)
	http.HandleFunc("POST /alarm/{id}/cancel", alarmCancelHandler)
}

// ---------------------------------------------------------------------------

func alarmCreateHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle errors graceful
	// TODO: Input validation

	user, err := user.FromRequest(r)
	if err != nil {
		log.Fatalf("unable to obtain user: %v", err)
	}

	err = r.ParseForm()
	if err != nil {
		slog.Error("parsing form data failed", "request", r)
	}

	deadline := r.PostForm.Get("deadline")
	description := r.PostForm.Get("description")

	ctx := context.Background()

	dbAlarm, err := dbRW.InsertAlarm(ctx, db.InsertAlarmParams{
		User:        user.Username,
		Datetime:    deadline,
		Description: sql.NullString{String: description, Valid: description != ""},
	})

	if err != nil {
		slog.Error("inserting new alarm failed", "error", err)
	}

	//slog.Info("insert new timer", "result", res, "request", r)
	slog.Info("DB: new alarm inserted", "db.Alarm", dbAlarm)

	dispatch.Register(dbAlarm)
}

// ---------------------------------------------------------------------------

func alarmCancelHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle errors graceful
	// TODO: Input validation
	// TODO: AuthZ: may this user cancel the alarm??

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		log.Fatalf("/alarm/{id:%s}/cancel : parsing alarm_id failed: %s", alarm_id_str, err)
	}

	ctx := context.Background()
	dbAlarm, err := dbRW.CancelAlarm(ctx, alarm_id)

	if err != nil {
		log.Fatalf("/alarm/{id:%d}/cancel : database: update alarm failed: %s", alarm_id, err)
	}

	dispatch.Cancel(dbAlarm)
}
