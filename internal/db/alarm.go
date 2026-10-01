package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

func (a Alarm) IsCancelled() bool {

	return a.CancelledAt.Valid
}

// ----------------------------------------------------------------------------

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

// ----------------------------------------------------------------------------

func (a Alarm) CreatedAtTime() (time.Time, error) {

	return a.CreatedAt.Time, nil
}

// ----------------------------------------------------------------------------

func (a Alarm) DeadlineTime() (time.Time, error) {

	return a.Deadline.Time, nil
}

// ----------------------------------------------------------------------------

func (a Alarm) CancelledAtTime() (time.Time, error) {

	return a.CancelledAt.Time, nil
}

// ----------------------------------------------------------------------------

func (a Alarm) Overdue() (bool, error) {

	ttl, err := a.TTL()
	if err != nil {
		slog.Error("db.Alarm.Overdue() failed: retrieving TTL failed", "error", err)
		return false, fmt.Errorf("db.Alarm.Overdue() failed: retrieving TTL failed: %w", err)
	}

	return ttl <= 0 && !a.IsCancelled(), nil
}

// ----------------------------------------------------------------------------

func (a Alarm) DbChannel() (NtfyChannel, error) {

	ctx := context.Background()

	channel, err := ROQ.SelectChannel(ctx, a.Channel)
	if err != nil {
		return NtfyChannel{}, fmt.Errorf("(db) error fetching channel(%s): %w", a.Channel, err)
	}

	return channel, nil
}

// ----------------------------------------------------------------------------

func (a Alarm) DbUser() (user.User, error) {

	return user.User{Username: a.User}, nil
}
