package main

import (
	"net/http"

	embedder "vsc.irc.ugent.be/itsupport/work-in-peace/assets"
)

func init() {

	http.Handle("/fonts/", http.FileServerFS(embedder.Fonts))
	http.Handle("/images/", http.FileServerFS(embedder.Images))
}
