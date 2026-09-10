//go:build debug
// +build debug

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

func ctlUserinfoHandler(w http.ResponseWriter, r *http.Request) {

	// user := user.User{
	// 	Username: "arnes",
	// 	Sso_sub:  "ac9020c8-80dc-0121-0cb8-796a89e1f05a",
	// }

	user, err := user.FromRequest(r)
	if err != nil {
		panic(err)
	}

	json, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(json))
}

func init() {
	ctlMux.HandleFunc("GET /ctl/userinfo", ctlUserinfoHandler)
}
