package dispatch

import (
	"fmt"
	"log/slog"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
)

type Dispatcher struct {
	alarms map[int64]db.Alarm
	timers map[int64]*time.Timer
}

var defaultDispatcher Dispatcher = Dispatcher{
	alarms: make(map[int64]db.Alarm),
	timers: make(map[int64]*time.Timer),
}

func Register(a db.Alarm) error {

	return defaultDispatcher.Register(a)
}

func (d *Dispatcher) Register(a db.Alarm) error {

	d.alarms[a.AlarmID] = a

	f := func() {
		slog.Warn("ALARM CALLBACK called", "db.Alarm", a)
	}

	ttl, err := a.TTL()

	if err != nil {
		slog.Error("unable to parse Deadline into Duration", "a.Deadline", a.Deadline, "a", a)
		return fmt.Errorf("unable to parse Deadline into Duration: %w", err)
	}

	t := time.AfterFunc(ttl, f)

	d.timers[a.AlarmID] = t

	return nil
}
