package main

import (
	"net/http"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/abort"
)

// ========================================================================== //
// CTL / CONTROL / DEV MUX
// ========================================================================== //

// types
type tokenMux struct {
	handler http.Handler
}

func (m *tokenMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	recToken := r.Header.Get(APP.Ctl.AuthHeader)
	if recToken == tokenMuxToken {
		m.handler.ServeHTTP(w, r)
		return
	}

	abort.Fatal(w, http.StatusUnauthorized, "invalid auth token",
		APP.Ctl.AuthHeader, recToken,
		"middleware", "tokenMux",
	)
}
