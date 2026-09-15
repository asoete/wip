package db

import (
	"fmt"
	"log/slog"
	"time"
)

func (a Alarm) TTL() (time.Duration, error) {

	deadline_t, err := time.Parse(time.DateTime, a.Deadline)

	// TODO: better error handling /propagation
	if err != nil {
		return time.Duration(0), err
	}

	return time.Until(deadline_t), nil
}
func (a Alarm) Overdue() (bool, error) {

	ttl, err := a.TTL()
	if err != nil {
		slog.Error("db.Alarm.Overdue() failed: retrieving TTL failed", "error", err)
		return false, fmt.Errorf("db.Alarm.Overdue() failed: retrieving TTL failed: %w", err)
	}

	return ttl <= 0 && !a.CancelledAt.Valid, nil
}
