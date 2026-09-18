package dispatch

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/alert"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
)

type Dispatcher struct {
	alarms map[int64]*db.Alarm
	timers map[int64]*time.Timer
}

var defaultDispatcher Dispatcher = Dispatcher{
	alarms: make(map[int64]*db.Alarm),
	timers: make(map[int64]*time.Timer),
}

// ---------------------------------------------------------------------------
// FACADES
// ---------------------------------------------------------------------------

func Resume(dbqs *db.Queries) error {
	return defaultDispatcher.Resume(dbqs)
}

func Register(a db.Alarm) error {
	return defaultDispatcher.Register(a)
}

func Cancel(a db.Alarm) error {
	return defaultDispatcher.Cancel(a)
}

// ---

func GetAlarm(id int64) (*db.Alarm, bool) {
	return defaultDispatcher.GetAlarm(id)
}

func GetTimer(id int64) (*time.Timer, bool) {
	return defaultDispatcher.GetTimer(id)
}

// ---

func GetAlarms() map[int64]*db.Alarm {
	return defaultDispatcher.GetAlarms()
}

func GetTimers() map[int64]*time.Timer {
	return defaultDispatcher.GetTimers()
}

// ---------------------------------------------------------------------------
// IMPLEMENTATION
// ---------------------------------------------------------------------------

func (d *Dispatcher) Resume(dbqs *db.Queries) error {

	ctx := context.Background()
	list, err := dbqs.ListActiveAlarms(ctx)
	if err != nil {
		slog.Error("unable to boot dispatch: loading active alarms failed", "error", err)
		os.Exit(2)
	}

	for _, a := range list {
		d.Register(a)
	}

	return nil
}

// ---------------------------------------------------------------------------

func (d *Dispatcher) Register(a db.Alarm) error {

	// TODO : add RWmutex

	d.alarms[a.AlarmID] = &a

	overdueAlert, err := alert.New(alert.Options{
		Channel: "WiP-devchannel-PxLeA",
		Title:   fmt.Sprintf("[WiP] %s has missed check-out", a.User),
		Body: fmt.Sprintf(
			"**`%s`** has set an alarm to fire at *`<%s>`* with description:\n",
			a.User, a.Deadline,
		) +
			fmt.Sprintf("> %s\n\n", a.Description.String) +
			"They have missed the deadline, **please reach out** an ensure they are well.",
	})

	if err != nil {
		slog.Error("unable to create a new alert", "error", err)
		os.Exit(1)
	}

	err = d.createCountdown(a, &overdueAlert)
	if err != nil {
		return fmt.Errorf("dispatch.createCountdown() failed: %w", err)
	}

	return nil
}

func (d *Dispatcher) createCountdown(alarm db.Alarm, alertInst *alert.Alert) error {

	var ttl time.Duration
	var err error

	if len(alertInst.Sent()) > 0 {
		ttl, err = alertInst.TTR()
	} else {
		ttl, err = alarm.TTL()
		if err != nil {
			return fmt.Errorf("unable to fetch TTL from alarm: %w", err)
		}
	}

	callback := func() {
		slog.Info("countdown callback called", "alertInst", alertInst)
		alert.Send(alertInst)
		d.createCountdown(alarm, alertInst)
	}

	t := time.AfterFunc(ttl, callback)

	d.timers[alarm.AlarmID] = t

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

func (d *Dispatcher) GetAlarm(id int64) (*db.Alarm, bool) {
	alarm, found := d.alarms[id]
	return alarm, found
}

func (d *Dispatcher) GetTimer(id int64) (*time.Timer, bool) {
	timer, found := d.timers[id]
	return timer, found
}

// ---------------------------------------------------------------------------

func (d *Dispatcher) GetAlarms() map[int64]*db.Alarm {
	return d.alarms
}

func (d *Dispatcher) GetTimers() map[int64]*time.Timer {
	return d.timers
}
