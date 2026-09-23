package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

// APP.Web.Address()
type Application struct {
	Pidfile    string
	DumpConfig bool
	Web        WebConfig
	Ctl        CtlConfig
	Db         DbConfig
}

// ----------------------------------------------------------------------------

func (a *Application) Init() {

	flag.BoolVar(&a.DumpConfig, "dump-config", false, "dump the active config and exit")
	flag.StringVar(&a.Pidfile, "pidfile", "", "write the program pid to this `path`")

	a.Web.Init()
	a.Ctl.Init()
	a.Db.Init()
}

func (a *Application) Boot() {
	flag.Parse()

	a.Web.Boot()
	a.Ctl.Boot()
	a.Db.Boot()
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
	Address string
}

// ----------------------------------------------------------------------------

func (wc *WebConfig) Init() {
	flag.StringVar(&wc.Address, "web.address", "127.0.0.1:8080", "bind to this `address`")
}

// ----------------------------------------------------------------------------

func (wc *WebConfig) Boot() {}

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
	flag.StringVar(&dc.Dsn, "db.dsn", "sqlite::memory:", "connect to this `dsn`")
}

// ----------------------------------------------------------------------------

func (dc *DbConfig) Boot() {}
