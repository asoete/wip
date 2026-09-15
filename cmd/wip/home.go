package main

import (
	"context"
	"html/template"
	"log"
	"log/slog"
	"net/http"

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
		"assets/templates/components/alarm/as-ul.html",
	))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	user, err := user.FromRequest(r)
	if err != nil {
		slog.Error("unable to retrieve user from request", "error", err)
	}

	ctx := context.Background()
	activeAlarms, err := dbRW.ListActiveUserAlarms(ctx, user.Username)
	if err != nil {
		slog.Error("unable to retrieve active user alarms", "user", user, "error", err)
	}

	log.Printf("user: %+v\n", user)
	log.Printf("activeAlarms: %+v\n", activeAlarms)

	cancelledAlarms, err := dbRW.ListCancelledUserAlarms(ctx, user.Username)
	if err != nil {
		slog.Error("unable to retrieve cancelled user alarms", "user", user, "error", err)
	}

	homePage.Execute(w, map[string]any{
		"user": user,
		"alarms": map[string]any{
			"active":    activeAlarms,
			"cancelled": cancelledAlarms,
		},
	})
}
