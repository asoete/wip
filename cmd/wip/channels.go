package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/abort"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

// ---------------------------------------------------------------------------

var channelsPage *template.Template

func init() {

	// register handler
	http.HandleFunc("GET /channels", channelsHandler)

	// register templates
	channelsPage = template.Must(template.ParseFiles(
		"assets/templates/layouts/default.html",
		"assets/templates/pages/channels.html",
		"assets/templates/components/page/header.html",
		"assets/templates/icons.html",
	))
}

func channelsHandler(w http.ResponseWriter, r *http.Request) {

	abrt := abort.New("handler", "channelHandler")

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.AuthError(w, err)
		return
	}

	abrt.Append("user", user)

	ctx := context.Background()
	channels, err := dbRW.ListChannels(ctx)
	if err != nil {
		abrt.DbError(w, err, "unable to list channels")
	}

	err = channelsPage.Execute(w, map[string]any{
		"user":     user,
		"channels": channels,
	})

	if err != nil {
		slog.Error("rendering channelsPage failed", "error", err)
	}
}

// ---------------------------------------------------------------------------

func init() {

	// register handler
	http.HandleFunc("POST /channels", newChannelHandler)
}

func newChannelHandler(w http.ResponseWriter, r *http.Request) {

	abrt := abort.New("handler", "newChannelHandler")

	user, err := user.FromRequest(r)
	if err != nil {
		abrt.AuthError(w, err)
		return
	}

	abrt.Append("user", user)

	err = r.ParseForm()
	if err != nil {
		abrt.InputError(w, err, "invalid input")
		return
	}

	channel_name := r.PostForm.Get("channel_name")

	if channel_name == "" {
		abrt.InputError(w, fmt.Errorf("channel is empty"), "invalid channel name",
			"channel_name", channel_name,
		)
		return
	}

	url, err := APP.Ntfy.SecretUrlFrom(channel_name)
	if err != nil {
		abrt.Error(w, err, "creating channel url failed", "channel_name", channel_name)
		return
	}

	ctx := context.Background()
	_, err = dbRW.InsertChannel(ctx, db.InsertChannelParams{
		Name: channel_name,
		Url:  url.String(),
	})

	if err != nil {
		abrt.DbError(w, err, "unable to create channel")
	}

	http.Redirect(w, r, "/channels", http.StatusSeeOther)
}
