package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/abort"
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

	abrt := abort.New("handler", "alarmCreateHandler")

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.Fatal(w, http.StatusForbidden,
			"error", err,
			"user", user,
		)
		return
	}
	abrt.Append("user", user)

	err = r.ParseForm()
	if err != nil {
		abrt.Fatal(w, http.StatusBadRequest,
			"info", "parsing form-data failed",
			"error", err,
		)
		return
	}

	deadline := r.PostForm.Get("deadline")
	description := r.PostForm.Get("description")
	channel_name := r.PostForm.Get("channel")

	if channel_name == "" {
		abrt.InputError(w, fmt.Errorf("channel is empty"), "invalid channel provided",
			"channel_name", channel_name,
		)
		return
	}

	fmt.Printf("channel_name: %v\n", channel_name)

	ctx := context.Background()

	var channel db.NtfyChannel
	channel, err = dbRW.SelectChannel(ctx, channel_name)

	if err != nil && err != sql.ErrNoRows {
		abrt.DbError(w, err, "invalid channel")
		return
	}

	if err == sql.ErrNoRows {
		url, err := APP.Ntfy.SecretUrlFrom(channel_name)
		if err != nil {
			abrt.Error(w, err, "creating channel failed", "channel", channel)
		}
		channel, err = dbRW.InsertChannel(ctx, db.InsertChannelParams{
			Name: channel_name,
			Url:  url.String(),
		})
	}

	dbAlarm, err := dbRW.InsertAlarm(ctx, db.InsertAlarmParams{
		User:        user.Username,
		Datetime:    deadline,
		Description: sql.NullString{String: description, Valid: description != ""},
		Channel:     channel.Name,
	})

	if err != nil {
		abrt.DbError(w, err, "storing new alarm failed")
		return
	}

	//slog.Info("insert new timer", "result", res, "request", r)
	slog.Info("DB: new alarm inserted", "db.Alarm", dbAlarm)

	err = dispatch.Register(dbAlarm)
	if err != nil {
		abrt.Error(w, err, "deadline update failed", "info", "dispatch.Register() failed")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---------------------------------------------------------------------------

func init() {
	http.HandleFunc("POST /alarm/{id}/cancel", alarmCancelHandler)
}

func alarmCancelHandler(w http.ResponseWriter, r *http.Request) {

	abrt := abort.New("handler", "alarmCancelHandler")

	// TODO: handle errors graceful
	// TODO: Input validation
	// TODO: AuthZ: may this user cancel the alarm??

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.AuthError(w, err)
		return
	}

	abrt.Append("user", user)

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		abrt.Fatal(w, http.StatusBadRequest, "invalid alarm_id",
			"{id}", alarm_id_str,
			"error", err,
		)
		return
	}

	abrt.Append("alarm_id", alarm_id)

	ctx := context.Background()
	dbAlarm, err := dbRW.CancelAlarm(ctx, alarm_id)

	if err != nil {
		abrt.DbError(w, err, "alarm cancel failed")
		return
	}

	err = dispatch.Cancel(dbAlarm)
	if err != nil {
		abrt.Error(w, err, "alarm cancel failed")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ---------------------------------------------------------------------------

func init() {
	http.HandleFunc("POST /alarm/{id}/quickadd", alarmQuickaddHandler)
}

func alarmQuickaddHandler(w http.ResponseWriter, r *http.Request) {

	abrt := abort.New("handler", "alarmQuickaddHandler")

	// TODO: handle errors graceful
	// TODO: Input validation
	// TODO: AuthZ: may this user cancel the alarm??

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.Fatal(w, http.StatusForbidden,
			"error", err,
			"user", user,
		)
		return
	}

	abrt.Append("user", user)

	alarm_id_str := r.PathValue("id")
	alarm_id, err := strconv.ParseInt(alarm_id_str, 10, 64)
	if err != nil {
		abrt.Fatal(w, http.StatusBadRequest, "invalid alarm_id",
			"{id}", alarm_id_str,
			"error", err,
		)
		return
	}

	abrt.Append("alarm_id", alarm_id)

	err = r.ParseForm()
	if err != nil {
		abrt.Fatal(w, http.StatusBadRequest, "invalid request data",
			"request", r,
			"error", err,
		)
		return
	}

	duration_str := r.PostForm.Get("quickadd")

	if duration_str == "" {
		abrt.Fatal(w, http.StatusBadRequest, "invalid duration",
			"quickadd/duration", duration_str,
			"error", "duration is empty",
		)
		return
	}

	duration, err := time.ParseDuration(duration_str)
	if err != nil {
		abrt.Fatal(w, http.StatusBadRequest, "invalid duration",
			"input", duration_str,
			"duration", duration,
			"error", err,
		)
		return
	}

	abrt.Append("quickadd", duration)

	ctx := context.Background()
	alarm, err := dbRW.SelectAlarm(ctx, alarm_id)
	if err != nil {
		abrt.DbError(w, err, "no such alarm", "info", "fetch alarm from database failed")
		return
	}

	old_deadline, err := alarm.DeadlineTime()
	if err != nil {
		abrt.Error(w, err, "invalid deadline")
		return
	}

	new_deadline := old_deadline.Add(duration)

	abrt.Append("old_deadline", old_deadline)
	abrt.Append("new_deadline", new_deadline)

	fmt.Printf(">>>> %d\n", time.Duration(time.Minute*30))

	alarm, err = dbRW.UpdateAlarmDeadline(ctx, db.UpdateAlarmDeadlineParams{
		Deadline: sqlite.Time{Time: new_deadline},
		AlarmID:  alarm_id,
	})

	if err != nil {
		abrt.DbError(w, err, "deadline update failed")
		return
	}

	err = dispatch.Register(alarm)
	if err != nil {
		abrt.Error(w, err, "deadline update failed", "info", "dispatch.Register() failed")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
