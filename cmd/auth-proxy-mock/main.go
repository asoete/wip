package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var listenAddr string
var forwardToAddr string

var remoteUser string
var ssoSub string

func init() {

	// options
	flag.StringVar(&listenAddr, "listen-on", "127.0.0.1:8888", "listen on this address for requests to forward")
	flag.StringVar(&forwardToAddr, "forward-to", "http://127.0.0.1:8080/", "forward requests to this address")

	// headers
	flag.StringVar(&remoteUser, "u", "", "shorthand: set REMOTE_USER http header")
	flag.StringVar(&remoteUser, "remote-user", "", "set REMOTE_USER http header")

	flag.StringVar(&ssoSub, "s", "", "shorthand: set SSO_SUB http header")
	flag.StringVar(&ssoSub, "sso-sub", "", "set SSO_SUB http header")
}

func main() {

	flag.Parse()

	forwardToUrl, err := url.Parse(forwardToAddr)

	if err != nil {
		log.Fatal(err)
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(forwardToUrl)

			// Headers sent to the target server.
			if remoteUser != "" {
				r.Out.Header.Set("REMOTE_USER", remoteUser)
			}

			if ssoSub != "" {
				r.Out.Header.Set("SSO_SUB", ssoSub)
			}
		},
	}

	log.Printf("APM (Auth Proxy Mock) listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))
}
