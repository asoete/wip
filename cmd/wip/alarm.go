package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db/sqlite"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/dispatch"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

func init() {
	http.HandleFunc("POST /alarm", alarmCreateHandler)
}

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---------------------------------------------------------------------------

func init() {
	http.HandleFunc("POST /alarm/{id}/cancel", alarmCancelHandler)
}

func alarmCancelHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle errors graceful
	// TODO: Input validation
	// TODO: AuthZ: may this user cancel the alarm??

	// user, err := user.FromRequest(r)
	// if err != nil {
	// 	log.Fatalf("unable to obtain user: %v", err)
	// }

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---------------------------------------------------------------------------

func init() {
	http.HandleFunc("POST /alarm/{id}/quickadd", alarmQuickaddHandler)
}

func alarmQuickaddHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle errors graceful
	// TODO: Input validation
	// TODO: AuthZ: may this user cancel the alarm??

	// user, err := user.FromRequest(r)
	// if err != nil {
	// 	log.Fatalf("unable to obtain user: %v", err)
	// }

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		log.Fatalf("/alarm/{id:%s}/quickadd : parsing alarm_id failed: %s", alarm_id_str, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		slog.Error("parsing form data failed", "request", r)
		return
	}

	duration_str := r.PostForm.Get("quickadd")

	if duration_str == "" {
		slog.Error("quickadd requires duration", "request", r)
		return
	}

	duration, err := time.ParseDuration(duration_str)
	if err != nil {
		slog.Error("parsing raw duration data failed", "duration", duration_str, "request", r)
		return
	}

	ctx := context.Background()
	alarm, err := dbRW.SelectAlarm(ctx, alarm_id)
	if err != nil {
		slog.Error("unable to select alarm", "id", alarm_id, "error", err)
	}

	old_deadline, err := alarm.DeadlineTime()
	if err != nil {
		slog.Error("unable to calculate new alarm deadline", "id", alarm_id, "error", err)
	}

	new_deadline := old_deadline.Add(duration)

	fmt.Printf("old deadline: %+v\n", old_deadline)
	fmt.Printf("new deadline: %+v\n", new_deadline)

	fmt.Printf(">>>> %d\n", time.Duration(time.Minute*30))

	alarm, err = dbRW.UpdateAlarmDeadline(ctx, db.UpdateAlarmDeadlineParams{
		Deadline: sqlite.Time{Time: new_deadline},
		AlarmID:  alarm_id,
	})

	if err != nil {
		slog.Error("unable to update alarm deadline", "id", alarm_id, "error", err)
	}

	dispatch.Register(alarm)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
