package db

import (
	"fmt"
	"log/slog"
	"time"
)

func (a Alarm) IsCancelled() bool {

	return a.CancelledAt.Valid
}

func (a Alarm) TTL() (time.Duration, error) {

	if a.IsCancelled() {
		return time.Duration(0), nil
	}

	deadline_t, err := a.DeadlineTime()

	// TODO: better error handling /propagation
	if err != nil {
		return time.Duration(0), err
	}

	return time.Until(deadline_t), nil
}

func (a Alarm) CreatedAtTime() (time.Time, error) {

	created_at_t, err := time.Parse(time.DateTime, a.CreatedAt)

	if err != nil {
		return time.Now(), fmt.Errorf("converting db.Alarm.CreatedAt to time.Time failed: %w", err)
	}

	return created_at_t, nil
}

func (a Alarm) DeadlineTime() (time.Time, error) {

	deadline_t, err := time.Parse(time.DateTime, a.Deadline)

	if err != nil {
		return time.Now(), fmt.Errorf("converting db.Alarm.Deadline to time.Time failed: %w", err)
	}

	return deadline_t, nil
}

func (a Alarm) CancelledAtTime() (time.Time, error) {

	cancelled_at_t, err := time.Parse(time.DateTime, a.CancelledAt.String)

	if err != nil {
		return time.Now(), fmt.Errorf("converting db.Alarm.CancelledAt to time.Time failed: %w", err)
	}

	return cancelled_at_t, nil
}

func (a Alarm) Overdue() (bool, error) {

	ttl, err := a.TTL()
	if err != nil {
		slog.Error("db.Alarm.Overdue() failed: retrieving TTL failed", "error", err)
		return false, fmt.Errorf("db.Alarm.Overdue() failed: retrieving TTL failed: %w", err)
	}

	return ttl <= 0 && !a.IsCancelled(), nil
}
