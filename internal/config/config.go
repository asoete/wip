package config

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"
)

var Web *WebConfig
var Ntfy *NtfyConfig

// APP.Web.Address()
type Application struct {
	Pidfile    string
	DumpConfig bool
	Web        WebConfig
	Ctl        CtlConfig
	Db         DbConfig
	Ntfy       NtfyConfig
}

// ----------------------------------------------------------------------------

func (a *Application) Init() {

	flag.BoolVar(&a.DumpConfig, "dump-config", false, "dump the active config and exit")
	flag.StringVar(&a.Pidfile, "pidfile", "", "write the program pid to this `path`")

	a.Web.Init()
	a.Ctl.Init()
	a.Db.Init()
	a.Ntfy.Init()
}

func (a *Application) Boot() {
	flag.Parse()

	a.Web.Boot()
	a.Ctl.Boot()
	a.Db.Boot()
	a.Ntfy.Boot()
}

func (a *Application) Dump() {

	fmt.Print("=== <config.Application> === \n")
	fmt.Print("  ")
	json, err := json.MarshalIndent(a, "  ", "  ")

	if err != nil {
		log.Fatalf("unable to dump config.Application to JSON: %w", err)
	}

	fmt.Printf("%s\n", json)
	fmt.Print("=== </config.Application> === \n")
}

// ============================================================================

type WebConfig struct {
	Address              string
	RemoteUserHeader     string
	RemoteUsernameHeader string
	RemoteGroupsHeader   string
} // ----------------------------------------------------------------------------

func (wc *WebConfig) Init() {
	flag.StringVar(&wc.Address, "web.address", "127.0.0.1:8080", "bind to this `address`")
	flag.StringVar(&wc.RemoteUserHeader, "web.remote-user-header", "REMOTE_USER", "this HTTP header holds the universal unique OpenIDc identifier")
	flag.StringVar(&wc.RemoteUsernameHeader, "web.remote-username-header", "OIDC_CLAIM_preferred_username", "this HTTP header holds the preferred username")
	flag.StringVar(&wc.RemoteGroupsHeader, "web.remote-groups-header", "OIDC_CLAIM_memberOf", "this HTTP header holds the comma separated groups list ")
}

// ----------------------------------------------------------------------------

func (wc *WebConfig) Boot() {

	// Expose as global :-/
	Web = wc
}

// ============================================================================

type CtlConfig struct {
	Address    string
	AuthHeader string `default:""`
	AuthToken  string
}

// ----------------------------------------------------------------------------

func (cc *CtlConfig) Init() {
	flag.StringVar(&cc.Address, "ctl.address", "127.0.0.1:8100", "bind service control to this `address`")
	flag.StringVar(&cc.AuthHeader, "ctl.auth-header", "X-WiP-Ctl-Token", "this HTTP header holds the auth token")
}

// ----------------------------------------------------------------------------

func (cc *CtlConfig) Boot() {}

// ============================================================================

type DbConfig struct {
	Dsn string
}

// ----------------------------------------------------------------------------

func (dc *DbConfig) Init() {
	flag.StringVar(&dc.Dsn, "db.dsn", ":memory:", "connect to this `dsn`")
}

// ----------------------------------------------------------------------------

func (dc *DbConfig) Boot() {}

// ============================================================================

type NtfyConfig struct {
	ServerUrlStr string
	serverUrl    *url.URL
	AlertDelay   time.Duration
}

// ----------------------------------------------------------------------------

func (nc *NtfyConfig) Init() {
	flag.StringVar(&nc.ServerUrlStr, "ntfy.server-url", "https://ntfy.sh", "NTFY.sh server `url`")

	nc.AlertDelay = time.Duration(time.Minute * 5)
	flag.Func("ntfy.alert-delay", "the `delay` between successive notification reminders", func(input string) error {
		dur, err := time.ParseDuration(input)
		nc.AlertDelay = dur
		return err
	})
}

// ----------------------------------------------------------------------------

func (nc *NtfyConfig) Boot() {

	// expose as global
	Ntfy = nc

	url, err := url.ParseRequestURI(nc.ServerUrlStr)

	if err != nil {
		log.Fatalf("config.Ntfy.Boot() failed: unable to parse url --ntfy.server-url='%s' : %s", nc.ServerUrlStr, err)
	}

	nc.serverUrl = url
}

// ----------------------------------------------------------------------------

func (nc *NtfyConfig) SecretUrlFrom(str string) (*url.URL, error) {

	hasher := sha1.New()
	_, err := hasher.Write([]byte(str))
	if err != nil {
		return &url.URL{}, err
	}

	hashValue := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	path := fmt.Sprintf("wip-%s-%.*s", str, 5, hashValue)

	return nc.UrlFrom(path), nil

}

// ----------------------------------------------------------------------------

func (nc *NtfyConfig) UrlFrom(str string) *url.URL {

	return nc.serverUrl.JoinPath(str)
}
