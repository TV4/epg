package epg

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

// Stockholm is the Time Zone in Sweden
var Stockholm *time.Location

func init() {
	if location, err := time.LoadLocation("Europe/Stockholm"); err == nil {
		Stockholm = location
	}
}

type Time struct {
	time.Time
}

func (t *Time) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "0001-01-01T00:00:00+01:00":
		attr.Value = "0001-01-01T00:00:00+00:00"
	case "9999-12-31T23:59:59+01:00":
		attr.Value = "9999-12-31T23:59:59+00:00"
	}

	var format string

	switch len(attr.Value) {
	case 25, 20:
		pt, err := time.Parse(time.RFC3339, attr.Value)
		if err != nil {
			return err
		}

		*t = Time{pt}

		return nil
	case 19:
		format = "2006-01-02T15:04:05"
	case 10:
		format = "2006-01-02"
	case 0:
		return nil
	}

	pt, err := time.ParseInLocation(format, attr.Value, Stockholm)
	if err != nil {
		return err
	}

	*t = Time{pt}

	return nil
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	return json.Marshal(t.Time)
}
