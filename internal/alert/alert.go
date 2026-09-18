package alert

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	// "vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
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
	Channel string
	Title   string
	Body    string
}

type Alert struct {
	Options
	channel string
	url     *url.URL
	sent    []time.Time
}

func New(o Options) (Alert, error) {

	a := Alert{}

	err := a.SetChannel(o.Channel)
	if err != nil {
		return a, fmt.Errorf("alert.SetChannel() failed: %w", err)
	}

	a.Title = o.Title
	a.Body = o.Body

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
		a.url.String(),
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

func (a *Alert) SetChannel(c string) error {

	ntfy_server := "https://msg.irc.ugent.be:1443"

	if c == "" {
		return fmt.Errorf("unable to create url.URL from ntfy_server(%s) + channel(%s): channel can not be empty", ntfy_server, c)
	}

	url_str := ntfy_server + "/" + c

	url, err := url.Parse(url_str)
	if err != nil {

		return fmt.Errorf("unable to create url.URL from ntfy_server(%s) + channel(%s): %w", ntfy_server, c, err)
	}

	a.channel = c
	a.url = url

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
