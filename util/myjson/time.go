package myjson

import (
	"encoding/json"
	"strconv"
	"time"
)

// Time represents a long date string of the following format: "2006-01-02T15:04:05.000Z".
type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	unquoteData, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	// attempt to parse time
	if parsedTime, err := time.Parse("2006-01-02T15:04:05.000-0700", unquoteData); err == nil {
		*t = Time(parsedTime)
		return nil
	}

	// attempt to parse time again
	if parsedTime, err := time.Parse("2006-01-02T15:04:05-07:00", unquoteData); err == nil {
		*t = Time(parsedTime)
		return nil
	}

	// attempt with a different format
	if parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z", unquoteData); err == nil {
		*t = Time(parsedTime)
		return nil
	}

	// attempt with yet another format
	if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", unquoteData); err != nil {
		return err
	} else {
		*t = Time(parsedTime)
	}

	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*t).Format("2006-01-02T15:04:05.000Z"))
}
func (t Time) ToTime() time.Time {
	return time.Time(t)
}

type Date time.Time

func (d *Date) UnmarshalJSON(data []byte) error {
	unquoteData, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", unquoteData)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*d).Format("2006-01-02"))
}

func (t Date) ToTime() time.Time {
	return time.Time(t)
}

// Millis represents a Unix time in milliseconds since January 1, 1970 UTC.
type Millis time.Time

func (m *Millis) UnmarshalJSON(data []byte) error {
	d, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*m = Millis(time.UnixMilli(d))
	return nil
}

func (m Millis) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(m).UnixMilli())
}

func (t Millis) ToTime() time.Time {
	return time.Time(t)
}

// Nanos represents a Unix time in nanoseconds since January 1, 1970 UTC.
type Nanos time.Time

func (n *Nanos) UnmarshalJSON(data []byte) error {
	d, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	// Go Time package does not include a method to convert UnixNano to a time.
	timeNano := time.Unix(d/1_000_000_000, d%1_000_000_000)
	*n = Nanos(timeNano)
	return nil
}

func (n Nanos) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(n).UnixNano())
}

func (t Nanos) ToTime() time.Time {
	return time.Time(t)
}
