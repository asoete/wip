package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/abort"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

// ---------------------------------------------------------------------------

var homePage *template.Template

func init() {

	// register handler
	http.HandleFunc("GET /{$}", homeHandler)

	// register templates
	homePage = template.Must(template.ParseFiles(
		"assets/templates/layouts/default.html",
		"assets/templates/pages/home.html",
		"assets/templates/components/alarm/card.html",
		"assets/templates/components/page/header.html",
		"assets/templates/icons.html",
	))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	abrt := abort.New("handler", "homeHandler")

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.AuthError(w, err)
		return
	}

	ctx := context.Background()
	activeAlarms, err := dbRW.ListActiveUserAlarms(ctx, user.Username)
	if err != nil {
		slog.Error("unable to retrieve active user alarms", "user", user, "error", err)
	}

	cancelledAlarms, err := dbRW.ListCancelledUserAlarms(ctx, user.Username)
	if err != nil {
		slog.Error("unable to retrieve cancelled user alarms", "user", user, "error", err)
	}

	// log.Printf("user: %+v\n", user)
	// log.Printf("activeAlarms: %+v\n", activeAlarms)
	// log.Printf("cancelledAlarms: %+v\n", cancelledAlarms)

	proposed_datetime := time.Now().Add(time.Hour * 2).Truncate(time.Hour)

	err = homePage.Execute(w, map[string]any{
		"user": user,
		"new_alarm": map[string]any{
			"proposed_datetime": proposed_datetime.Format("2006-01-02 15:04"),
			"min_datetime":      time.Now().Truncate(time.Minute * 15).Format("2006-01-02 15:04"),
		},
		"alarms": map[string]any{
			"active":    activeAlarms,
			"cancelled": cancelledAlarms,
		},
	})

	if err != nil {
		slog.Error("rendering homePage failed", "error", err)
	}
}
