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
	pt, err := time.ParseInLocation("2006-01-02T15:04:05", attr.Value, Stockholm)
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
