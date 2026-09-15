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

func Cancel(a db.Alarm) error {
	return defaultDispatcher.Cancel(a)
}

// ---

func GetAlarm(id int64) (db.Alarm, bool) {
	return defaultDispatcher.GetAlarm(id)
}

func GetTimer(id int64) (*time.Timer, bool) {
	return defaultDispatcher.GetTimer(id)
}

// ---------------------------------------------------------------------------

func (d *Dispatcher) Register(a db.Alarm) error {

	// TODO : add RWmutex

	d.alarms[a.AlarmID] = a

	f := func() {
		slog.Warn("ALARM CALLBACK called", "db.Alarm", a)
		sendAlert(a)
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

// ---------------------------------------------------------------------------

func (d *Dispatcher) Cancel(a db.Alarm) error {

	// TODO : add RWmutex

	ntfyTimer, found := d.timers[a.AlarmID]

	if !found {
		slog.Warn("dispatch.Cancel(): not timer found for alarm", "alarm_id", a.AlarmID, "db.Alarm", a)
	} else {
		if ok := ntfyTimer.Stop(); !ok {
			slog.Warn(
				"dispatch.Cancel(): unable to stop AfterFunc Timer: no timer associated with alarm",
				"alarm_id", a.AlarmID,
				"db.Alarm", a,
			)
		}
	}

	delete(d.timers, a.AlarmID)
	delete(d.alarms, a.AlarmID)

	return nil
}

// ---------------------------------------------------------------------------

func (d *Dispatcher) GetAlarm(id int64) (db.Alarm, bool) {
	alarm, found := d.alarms[id]
	return alarm, found
}

func (d *Dispatcher) GetTimer(id int64) (*time.Timer, bool) {
	timer, found := d.timers[id]
	return timer, found
}
