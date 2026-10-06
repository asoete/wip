package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	embedder "vsc.irc.ugent.be/itsupport/work-in-peace/assets"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/abort"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

// ---------------------------------------------------------------------------

var homePage *template.Template

func init() {

	// register handler
	http.HandleFunc("GET /{$}", homeHandler)

	// register templates
	homePage = template.Must(template.ParseFS(embedder.Templates,
		"templates/layouts/default.html",
		"templates/pages/home.html",
		"templates/components/alarm/card.html",
		"templates/components/page/header.html",
		"templates/icons.html",
	))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	abrt := abort.New("handler", "homeHandler")

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.AuthError(w, err)
		return
	}

	abrt.Append("user", user)

	ctx := context.Background()
	activeAlarms, err := db.ROQ.ListActiveUserAlarms(ctx, user.Username)
	if err != nil {
		abrt.DbError(w, err, "error retrieving active alarms")
		return
	}

	cancelledAlarms, err := db.ROQ.ListCancelledUserAlarms(ctx, user.Username)
	if err != nil {
		abrt.DbError(w, err, "error retrieving history")
		return
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
