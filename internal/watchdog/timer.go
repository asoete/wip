package watchdog

import (
	"time"
)

type Timer struct {
	*time.Timer
	Deadline time.Time
}

// ----------------------------------------------------------------------------

func (timer *Timer) Stop() bool {
	return timer.Timer.Stop()
}
