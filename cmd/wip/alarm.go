package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/dispatch"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

func alarmPostHandler(w http.ResponseWriter, r *http.Request) {

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

func init() {
	http.HandleFunc("POST /alarm", alarmPostHandler)
}
