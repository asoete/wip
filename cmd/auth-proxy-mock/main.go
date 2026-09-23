package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var listenAddr string
var forwardToAddr string

type strlist []string
var xheaders strlist

func (l *strlist) Set( v string ) error {
	*l = append(*l, v)
	return nil
}

func (l *strlist) String() string {
	return fmt.Sprintf("%s", *l)
}

func init() {

	// options
	flag.StringVar(&listenAddr, "listen-on", "127.0.0.1:8888", "listen on this address for requests to forward")
	flag.StringVar(&forwardToAddr, "forward-to", "http://127.0.0.1:8080/", "forward requests to this address")

	// headers
	flag.Var(&xheaders, "header", "specify additional HTTP headers (`name:value`)")
}

func main() {

	flag.Parse()

	forwardToUrl, err := url.Parse(forwardToAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("== DEBUG: dump config =====================================\n")

	configmap := map[string]string{
		"--listen-on": listenAddr,
		"--forward-to": forwardToAddr,
	}

	for name, value := range configmap {
		fmt.Printf("%20s = %s\n", name, value)
	}

	for _, input := range xheaders {
		parts := strings.SplitN(input, ":", 2)
		if len(parts) < 2 {
			slog.Warn("HTTP --header has no value!", "name", parts[0], "value", "")
			parts = append(parts, "")
		}
		fmt.Printf("%20s = %s\n", parts[0], parts[1])
	}

	fmt.Printf("===========================================================\n")

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(forwardToUrl)

			for _, input := range xheaders {
				parts := strings.SplitN(input, ":", 2)
				if len(parts) < 2 {
					parts = append(parts, "")
				}
				r.Out.Header.Set(parts[0], parts[1])
			}
		},
	}

	log.Printf("APM (Auth Proxy Mock) listening on %s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, proxy))
}
