package user

import (
	"fmt"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Sso_sub  string `json:"sso_sub"`
}

func FromRequest(r *http.Request) (User, error) {

	user := User{}

	user.Username = r.Header.Get("REMOTE_USER")
	if user.Username == "" {
		return user, fmt.Errorf("error creating user from Request: REMOTE_USER is empty")
	}

	user.Sso_sub = r.Header.Get("SSO_SUB")
	if user.Sso_sub == "" {
		return user, fmt.Errorf("error creating user from Request: SSO_SUB is empty")
	}

	return user, nil
}
