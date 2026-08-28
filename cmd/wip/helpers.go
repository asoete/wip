package main

import (
	"fmt"
	"log"
	"net/http"
)

func abort(w http.ResponseWriter, code int, format string, a ...any) {

	format = fmt.Sprintf("%d %s: %s", code, http.StatusText(code), format)
	http.Error(w, fmt.Sprintf(format, a...), code)
	log.Printf(format, a...)
}
