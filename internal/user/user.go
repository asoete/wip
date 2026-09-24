package user

import (
	"fmt"
	"net/http"
	"strings"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/config"
)

type User struct {
	Username   string
	RemoteUser string
	Groups     []string
}

func FromRequest(r *http.Request) (User, error) {

	user := User{}

	user.Username = r.Header.Get(config.Web.RemoteUsernameHeader)
	if user.Username == "" {
		return user, fmt.Errorf("header '%s' (config.Web.RemoteUsernameHeader) is empty", config.Web.RemoteUsernameHeader)
	}

	user.RemoteUser = r.Header.Get(config.Web.RemoteUserHeader)
	if user.RemoteUser == "" {
		return user, fmt.Errorf("header '%s' (config.Web.RemoteUserHeader) is empty", config.Web.RemoteUserHeader)
	}

	group_str := r.Header.Get(config.Web.RemoteGroupsHeader)
	if user.RemoteUser == "" {
		return user, fmt.Errorf("header '%s' (config.Web.RemoteGroupsHeader) is empty", config.Web.RemoteGroupsHeader)
	}

	user.Groups = parseGroupsClaim(group_str)

	return user, nil
}

var specialGroups map[string]struct{} = map[string]struct{}{
	"ccf": {},
	"it":  {},
	"tcf": {},
}

func parseGroupsClaim(str string) []string {

	all_groups := strings.Split(str, ",")

	var groups []string
	for _, group := range all_groups {

		if strings.HasPrefix(group, "u_") || strings.HasPrefix(group, "t_") || strings.HasPrefix(group, "s_") {
			groups = append(groups, group)
		}

		if strings.HasSuffix(group, "core") {
			groups = append(groups, group)
		}

		if _, ok := specialGroups[group]; ok {
			groups = append(groups, group)
		}
	}

	return groups
}
