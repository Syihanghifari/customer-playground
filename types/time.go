package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan %T into NullTime", value)
	}
	nt.Time = t
	nt.Valid = true
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		nt.Valid = false
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	nt.Time = t
	nt.Valid = true
	return nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte(`null`), nil
	}
	return []byte(`"` + nt.Time.Format(time.RFC3339) + `"`), nil
}
