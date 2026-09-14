package dispatch

import (
	"fmt"
	"net/http"
	"strings"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
)

func sendAlert(a db.Alarm) error {

	// err := alarm.Hydrate()

	req, err := http.NewRequest(
		"POST",
		"https://ntfy.sh/WiP-testing-channel_OfTD36Z4uARQeJsQB76pNQMoloPv",
		strings.NewReader(a.Description.String),
	)

	if err != nil {
		return fmt.Errorf("dispatch.sendAlert(): creating new alert failed: %w", err)
	}
	req.Header.Set("Title", fmt.Sprintf("WiP(%s): Deadline expired", a.User))
	req.Header.Set("Tags", "alarm_clock")
	http.DefaultClient.Do(req)

	return nil
}
