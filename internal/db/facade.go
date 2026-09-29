package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
)

var roDB *sql.DB
var rwDB *sql.DB

var ROQ *Queries
var rwq *Queries

// ---------------------------------------------------------------------------

func BootRO(dsn string) error {

	// options: https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc#Driver.Open
	options := url.Values{}
	options.Set("mode", "ro")
	options.Set("_busy_timeout", "5000")

	dsn = fmt.Sprintf("%s?%s", dsn, options.Encode())

	slog.Debug("db.BootRO(): dns with options", "dsn", dsn)

	var err error
	roDB, err = sql.Open("sqlite", dsn)

	if err != nil {
		return fmt.Errorf("sql.Open() failed: %w", err)
	}

	ROQ = New(roDB)

	return nil
}

func BootRW(dsn string) error {

	// options: https://pkg.go.dev/modernc.org/sqlite?utm_source=godoc#Driver.Open
	// good settings: https://kerkour.com/sqlite-for-servers
	options := url.Values{}
	options.Set("_journal_mode", "WAL")
	options.Set("_synchronous", "NORMAL")
	options.Set("_busy_timeout", "5000")

	dsn = fmt.Sprintf("%s?%s", dsn, options.Encode())

	slog.Debug("db.BootRW(): dns with options", "dsn", dsn)

	var err error
	rwDB, err = sql.Open("sqlite", dsn)

	if err != nil {
		return fmt.Errorf("sql.Open() failed: %w", err)
	}

	rwq = New(rwDB)

	return nil
}

// ---------------------------------------------------------------------------

func StartTx() (*sql.Tx, *Queries, error) {

	tx, err := rwDB.Begin()
	if err != nil {
		return tx, &Queries{}, fmt.Errorf("db.StartTx() failed: %w", err)
	}

	return tx, rwq.WithTx(tx), nil
}
