package abort

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func New(args ...any) Aborter {
	return Aborter{tags: args}
}

// ---------------------------------------------------------------------------
// FACADES
// ---------------------------------------------------------------------------

var defaultAborter Aborter = Aborter{}

func Fatal(w http.ResponseWriter, httpStatusCode int, args ...any) {
	defaultAborter.Fatal(w, httpStatusCode, args)
}

// ---------------------------------------------------------------------------
// IMPLEMENTATION
// ---------------------------------------------------------------------------

type Aborter struct {
	tags []any
}

// ---------------------------------------------------------------------------

func (a *Aborter) Fatal(w http.ResponseWriter, httpStatusCode int, args ...any) {

	var userMsg string

	// args is odd, we have a custom error message
	if len(args)%2 == 1 {
		userMsg = fmt.Sprintf("%d %s: %s", httpStatusCode, http.StatusText(httpStatusCode), args[0])
		args = args[1:]
	} else {
		userMsg = fmt.Sprintf("%d %s", httpStatusCode, http.StatusText(httpStatusCode))
	}

	// Write response to HTTP Client
	http.Error(w, userMsg, httpStatusCode)

	args = append(a.tags, args...)

	// Annotate with file and line number
	_, file, linenr, ok := runtime.Caller(1)
	if ok {
		pwd, err := os.Getwd()
		if err != nil {
			slog.Error("abort.Log(): removing pwd() from file failed", "error", err)
		} else {
			file = strings.TrimPrefix(file, pwd+"/")
			args = append(args, "file", file, "line", linenr)
		}
	}

	slog.Error(userMsg, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) Append(args ...any) {

	a.tags = append(a.tags, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) DbError(w http.ResponseWriter, err error, args ...any) {

	var msg any
	if len(args)%2 == 1 {
		msg = args[0]
		args = args[1:]
	} else {
		msg = "no such item"
	}

	args = append([]any{msg}, args...)
	args = append(args, "error", err)

	if err == sql.ErrNoRows {
		a.Fatal(w, http.StatusNotFound, args...)
		return
	}

	a.Fatal(w, http.StatusNotFound, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) Error(w http.ResponseWriter, err error, userMsg string, args ...any) {

	args = append([]any{userMsg}, args...)
	args = append(args, "error", err)

	a.Fatal(w, http.StatusInternalServerError, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) AuthError(w http.ResponseWriter, err error, args ...any) {

	args = append(args, "error", err)

	a.Fatal(w, http.StatusForbidden, args...)
}
