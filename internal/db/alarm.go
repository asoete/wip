package db

import (
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
