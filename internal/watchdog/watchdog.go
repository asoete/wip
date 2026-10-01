package watchdog

import (
	"fmt"
	"log/slog"
	"time"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/alert"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/config"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/event"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

type Watchdog struct {
	Alarm   *db.Alarm
	Channel *db.NtfyChannel
	User    *user.User
	Timer   *Timer
}

func New(alarm *db.Alarm, channel *db.NtfyChannel, user *user.User) *Watchdog {

	return &Watchdog{
		Alarm:   alarm,
		Channel: channel,
		User:    user,
	}
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) StartCountdown() error {

	watchdog.Timer = watchdog.createNewCountdownTimer(watchdog.Alarm.Deadline.Time)

	return nil
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) Update(alarm *db.Alarm) error {

	watchdog.Alarm = alarm
	watchdog.resetCountdownTimer(alarm.Deadline.Time)
	return nil
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) Cancel() error {

	watchdog.Timer.Stop()
	return nil
}

// ============================================================================

func (watchdog *Watchdog) createNewCountdownTimer(deadline time.Time) *Timer {

	callback := func() {
		watchdog.alert()
		watchdog.delayCountdownTimer(config.Ntfy.AlertDelay)
	}

	duration := time.Until(deadline)

	return &Timer{
		Timer:    time.AfterFunc(duration, callback),
		Deadline: deadline,
	}
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) resetCountdownTimer(deadline time.Time) {

	watchdog.Timer.Stop()

	watchdog.Timer = watchdog.createNewCountdownTimer(deadline)
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) delayCountdownTimer(delay time.Duration) {

	watchdog.resetCountdownTimer(time.Now().Add(delay))
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) Id() int64 {

	return watchdog.Alarm.AlarmID
}

// ----------------------------------------------------------------------------

func (watchdog *Watchdog) alert() error {
	slog.Debug("watchdog.alert()", "id", watchdog.Id(), "deadline", watchdog.Timer.Deadline.Format(time.DateTime))

	alertData := alert.Options{
		AlarmID: watchdog.Id(),
		Channel: *watchdog.Channel,
		Title:   fmt.Sprintf("[WiP] %s has missed check-out", watchdog.Alarm.User),
		Body: fmt.Sprintf(
			"**`%s`** has set an alarm to fire at *`<%s>`* with description:\n",
			watchdog.Alarm.User, watchdog.Alarm.Deadline,
		) +
			fmt.Sprintf("> %s\n\n", watchdog.Alarm.Description.String) +
			"They have missed the deadline, **please reach out** an ensure they are well.",
	}

	err := alert.Send(alertData)

	if err != nil {
		return fmt.Errorf("watchdog.alert(): sending alert for alarm(%d) failed: %w", watchdog.Alarm.AlarmID, err)
	}

	go event.Log(event.AlertFire, *watchdog.User, alertData)
	return nil
}
