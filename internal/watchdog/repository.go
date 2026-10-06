package watchdog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
)

var defaultRepository Repository = Repository{
	watchers: make(map[int64]*Watchdog),
}

func Resume() error {
	return defaultRepository.Resume()
}

func Register(watchdog *Watchdog) error {
	return defaultRepository.Register(watchdog)
}

func UpdateAlarm(alarm *db.Alarm) error {
	return defaultRepository.UpdateAlarm(alarm)
}

func CancelAlarm(alarm *db.Alarm) error {
	return defaultRepository.CancelAlarm(alarm)
}

func GetWatchdog(id int64) (*Watchdog, bool) {
	return defaultRepository.GetWatchdog(id)
}

func GetWatchers() map[int64]*Watchdog {
	return defaultRepository.GetWatchers()
}

// ===========================================================================

type Repository struct {
	mutex    sync.Mutex
	watchers map[int64]*Watchdog
}

// ---------------------------------------------------------------------------

func (repository *Repository) Resume() error {

	ctx := context.Background()

	list, err := db.ROQ.ListActiveAlarms(ctx)
	if err != nil {
		slog.Error("unable to boot repository: loading active alarms failed", "error", err)
		os.Exit(2)
	}

	for _, alarm := range list {
		slog.Debug("db load existing alarm", "alarm_id", alarm.AlarmID)

		channel, err := alarm.DbChannel()
		if err != nil {
			return fmt.Errorf("watchdog.Resume(): unable to load NtfyChannel for alarm(%d): %w", alarm.AlarmID, err)
		}
		user, err := alarm.DbUser()
		if err != nil {
			return fmt.Errorf("watchdog.Resume(): unable to load User for alarm(%d): %w", alarm.AlarmID, err)
		}

		err = repository.Register(New(&alarm, &channel, &user))
		if err != nil {
			return fmt.Errorf("repository.Resume(): unable to register alarm(%d): %w", alarm.AlarmID, err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------

func (repository *Repository) Register(watchdog *Watchdog) error {

	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.watchers[watchdog.Id()] = watchdog

	watchdog.StartCountdown()

	return nil
}

// ---------------------------------------------------------------------------

func (repository *Repository) UpdateAlarm(alarm *db.Alarm) error {

	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	watchdog, found := repository.watchers[alarm.AlarmID]
	if !found {
		fmt.Errorf("watchdog.Update(%d): requested alarm not found in list of watchdogs", alarm.AlarmID)
	}

	err := watchdog.Update(alarm)
	if err != nil {
		return fmt.Errorf("watchdog.Update(%d) failed: %w", alarm.AlarmID, err)
	}

	return nil
}

// ---------------------------------------------------------------------------

func (repository *Repository) CancelAlarm(alarm *db.Alarm) error {

	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	watchdog, found := repository.watchers[alarm.AlarmID]
	if !found {
		return fmt.Errorf("watchdog.Repository.Cancel(%d): requested alarm not found in list of watchdogs", alarm.AlarmID)
	}

	err := watchdog.Cancel()
	if err != nil {
		return fmt.Errorf("watchdog.Cancel(%d) failed: %w", alarm.AlarmID, err)
	}

	delete(repository.watchers, alarm.AlarmID)

	return nil
}

// ---------------------------------------------------------------------------

func (repository *Repository) GetWatchdog(id int64) (*Watchdog, bool) {

	w, ok := repository.watchers[id]
	return w, ok
}

func (repository *Repository) GetWatchers() map[int64]*Watchdog {

	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	return repository.watchers
}

// ---------------------------------------------------------------------------
