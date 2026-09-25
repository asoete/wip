package alert

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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
	Prio    Prio
	Channel db.NtfyChannel
	Title   string
	Body    string
}

type Alert struct {
	Options
	channel db.NtfyChannel
	sent    []time.Time
}

func New(o Options) (Alert, error) {

	a := Alert{}

	a.Title = o.Title
	a.Body = o.Body

	a.channel = o.Channel

	a.Prio = PrioDefault
	if o.Prio != 0 {
		a.Prio = o.Prio
	}

	return a, nil
}

func Send(a *Alert) error {

	// err := alarm.Hydrate()

	req, err := http.NewRequest(
		"POST",
		a.channel.Url,
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

	a.sent = append(a.sent, time.Now())

	return nil
}

// ---------------------------------------------------------------------------

func (a *Alert) Sent() []time.Time {

	return a.sent
}

// ---------------------------------------------------------------------------

// Time To Reminder
func (a *Alert) TTR() (time.Duration, error) {

	if len(a.sent) == 0 {
		return time.Duration(0), fmt.Errorf("we should not be calling TTR (Time To next Reminder) when the alert was never sent before")
	}

	return time.Until(a.sent[len(a.sent)-1].Add(time.Minute * 5)), nil
}
