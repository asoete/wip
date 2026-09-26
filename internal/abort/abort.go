package abort

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"

	"modernc.org/sqlite"
	"modernc.org/sqlite/lib"
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
	a.fatal(w, httpStatusCode, args...)
}

// Make fatal() private and proxy public methods to this function, this way,
// runtime.Caller(2) always points to the correct location

func (a *Aborter) fatal(w http.ResponseWriter, httpStatusCode int, args ...any) {

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
	_, file, linenr, ok := runtime.Caller(2)
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
		msg = "database error"
	}

	args = append([]any{msg}, args...)
	args = append(args, "error", err)

	if err == sql.ErrNoRows {
		a.fatal(w, http.StatusNotFound, args...)
		return
	}

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {

		if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			fmt.Printf("%+v", sqliteErr)
			args[0] = "UNIQUE constraint violated"
			a.fatal(w, http.StatusConflict, args...)
			return
		}
	}

	a.fatal(w, http.StatusInternalServerError, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) Error(w http.ResponseWriter, err error, userMsg string, args ...any) {

	args = append([]any{userMsg}, args...)
	args = append(args, "error", err)

	a.fatal(w, http.StatusInternalServerError, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) AuthError(w http.ResponseWriter, err error, args ...any) {

	args = append(args, "error", err)

	a.fatal(w, http.StatusForbidden, args...)
}

// ---------------------------------------------------------------------------

func (a *Aborter) InputError(w http.ResponseWriter, err error, args ...any) {

	args = append(args, "error", err)

	a.fatal(w, http.StatusBadRequest, args...)
}
