package aws

import (
	"encoding/json"
	"time"
)

type TimeMillis int64

func (t *TimeMillis) UnmarshalJSON(b []byte) error {
	var timeStr string
	if err := json.Unmarshal(b, &timeStr); err != nil {
		return err
	}
	tt, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return err
	}
	*t = TimeMillis(tt.UnixNano() / 1e6) // Convert nanoseconds to milliseconds
	return nil
}

type UnixMillis int64

// UnmarshalJSON converts an ISO 8601 formatted time string to Unix milliseconds.
func (u *UnixMillis) UnmarshalJSON(data []byte) error {
	// Parse the string to time.Time
	t, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	if err != nil {
		return err
	}
	// Convert to Unix milliseconds and assign to UnixMillis
	*u = UnixMillis(t.UnixNano() / int64(time.Millisecond))
	return nil
}

type JSONString string

func (s *JSONString) UnmarshalJSON(data []byte) error {
	*s = JSONString(data)
	return nil
}
