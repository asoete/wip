package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/db"
	"vsc.irc.ugent.be/itsupport/work-in-peace/internal/user"
)

// ============================================================================
// FACADE
// ============================================================================

var defaultLogger *logger = &logger{}

func SetDefault(l *logger) {
	defaultLogger = l
}

func Log(ev EventType, u user.User, data any) {
	defaultLogger.Log(ev, u, data)
}

// ============================================================================
// IMPLEMENTATION
// ============================================================================

func New(dbi *db.Queries) *logger {

	return &logger{}
}

type logger struct{}

type EventType int

const (
	AlarmCreate EventType = iota
	AlarmModify
	AlarmCancel
	AlertFire
	AlertCancel
	AppStart
	AppStop
	ChannelCreate
	ChannelModify
	ChannelDelete
	UserLogin
)

var EventName = map[EventType]string{
	AlarmCreate:   "alarm:create",
	AlarmModify:   "alarm:modify",
	AlarmCancel:   "alarm:cancel",
	AlertFire:     "alert:fire",
	AlertCancel:   "alert:cancel",
	AppStart:      "app:start",
	AppStop:       "app:stop",
	ChannelCreate: "channel:create",
	ChannelModify: "channel:modify",
	ChannelDelete: "channel:delete",
	UserLogin:     "user:login",
}

func (l *logger) Log(event_type EventType, u user.User, data any) {

	// -- handle eventType to eventName

	var event_name string
	event_name, ok := EventName[event_type]
	if !ok {
		event_name = fmt.Sprintf("event_type:%d", event_type)
		slog.Error("unknown event type: (%d), using fallback", "EventType", event_type, "fallback", event_name)
	}

	// -- Obtain sub ids

	subid, subidValid := atoSubid(data)

	// -- handle data

	var str string
	bts, err := json.Marshal(data)

	if err != nil {
		slog.Error("event.Log(): unable to json.Marshal data", "error", err)
		str = fmt.Sprintf("%+v", data)
	} else {
		str = string(bts)
	}

	eventParams := db.InsertEventParams{
		User: u.Username,
		Type: event_name,
		Subid: sql.NullInt64{
			Int64: subid,
			Valid: subidValid,
		},
		Data: sql.NullString{
			String: str,
			Valid:  str != "",
		},
	}

	ctx := context.Background()
	tx, dbQ, err := db.StartTx()
	if err != nil {
		slog.Error("event.Log(): start db transaction failed", "error", err, "event_details", eventParams)
		return
	}
	defer tx.Rollback()

	err = dbQ.InsertEvent(ctx, eventParams)

	err = tx.Commit()
	if err != nil {
		slog.Error("event.Log(): db transaction failed", "error", err, "event_details", eventParams)
	}
}

func atoSubid(subject any) (int64, bool) {

	switch x := subject.(type) {
	case db.Alarm:
		return x.AlarmID, true
	case db.NtfyChannel:
		return x.ChannelID, true
	}

	return 0, false
}
