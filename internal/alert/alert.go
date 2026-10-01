package alert

import (
	"fmt"
	"net/http"
	"strings"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
)

type Prio int

const (
	PrioMin     Prio = 1
	PrioLow     Prio = 2
	PrioDefault Prio = 3
	PrioHigh    Prio = 4
	PrioMax     Prio = 5
)

type Options struct {
	AlarmID int64
	Prio    Prio
	Channel db.NtfyChannel
	Title   string
	Body    string
}

func Send(a Options) error {

	// err := alarm.Hydrate()

	req, err := http.NewRequest(
		"POST",
		a.Channel.Url,
		strings.NewReader(a.Body),
	)

	if err != nil {
		return fmt.Errorf("alert.Send(): creating new HTTP request failed: %w", err)
	}

	req.Header.Set("Title", a.Title)
	req.Header.Set("Tags", "alarm_clock")
	req.Header.Set("Markdown", "yes")
	_, err = http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("alert.Send(): POSTing alert to NTFY failed: %w", err)
	}

	return nil
}
