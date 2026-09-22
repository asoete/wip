package sqlite

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) Scan(value any) error {
	var parsed time.Time

	switch v := value.(type) {
	case time.Time:
		parsed = v
	case string:
		var err error
		parsed, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
	case []byte:
		var err error
		parsed, err = time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot scan %T into LocalTime", value)
	}

	t.Time = parsed.In(time.Local)
	return nil
}

func (t Time) Value() (driver.Value, error) {
	return t.Time.UTC().Format(time.DateTime), nil
}

// ----------------------------------------------------------------------------

type NullTime struct {
	time.Time
	Valid bool
}

func (t *NullTime) Scan(value any) error {
	var parsed time.Time
	var valid bool = true

	switch v := value.(type) {
	case nil:
		valid = false
	case time.Time:
		parsed = v
	case string:
		var err error
		parsed, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
	case []byte:
		var err error
		parsed, err = time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cannot scan %T into LocalTime", value)
	}

	t.Time = parsed.In(time.Local)
	t.Valid = valid
	return nil
}

func (t NullTime) Value() (driver.Value, error) {
	return t.Time.UTC().Format(time.DateTime), nil
}
