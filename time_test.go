package epg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"testing"
	"time"
)

func TestUnmarshalXMLAttr(t *testing.T) {
	for _, tt := range []struct {
		attr xml.Attr
		want time.Time
		err  error
	}{
		{xml.Attr{}, time.Time{}, nil},
		{xml.Attr{Value: "2017-01-02"}, time.Date(2017, 1, 2, 0, 0, 0, 0, Stockholm), nil},
		{xml.Attr{Value: "2017-01-02T14:28:56"}, time.Date(2017, 1, 2, 14, 28, 56, 0, Stockholm), nil},
		{xml.Attr{Value: "2017-01-02T14:28:56+02:00"}, time.Date(2017, 1, 2, 13, 28, 56, 0, Stockholm), nil},
		{xml.Attr{Value: "not-a-date"}, time.Time{}, errors.New(
			`parsing time "not-a-date" as "2006-01-02": cannot parse "not-a-date" as "2006"`,
		)},
	} {
		et := &Time{}

		err := et.UnmarshalXMLAttr(tt.attr)
		if err != nil {
			if tt.err == nil {
				t.Fatalf("et.UnmarshalXMLAttr(%#v) = %v, want nil", tt.attr, err)
			}

			if got, want := err.Error(), tt.err.Error(); got != want {
				t.Fatalf("err.Error() = %v, want %v", got, want)
			}
		}

		if !et.Time.Equal(tt.want) {
			t.Fatalf("%v != %v", et.Time, tt.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var b bytes.Buffer

		w := bufio.NewWriter(&b)

		v := struct {
			Time *Time `json:"time"`
		}{nil}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w.Flush()

		if got, want := b.String(), `{"time":null}`+"\n"; got != want {
			t.Fatalf("b.String() = %q, want %q", got, want)
		}
	})

	t.Run("zero", func(t *testing.T) {
		var b bytes.Buffer

		w := bufio.NewWriter(&b)

		v := struct {
			Time *Time `json:"time"`
		}{&Time{}}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w.Flush()

		if got, want := b.String(), `{"time":null}`+"\n"; got != want {
			t.Fatalf("b.String() = %q, want %q", got, want)
		}
	})

	t.Run("CET", func(t *testing.T) {
		var b bytes.Buffer

		w := bufio.NewWriter(&b)

		v := struct {
			Time *Time `json:"time"`
		}{&Time{
			time.Date(2017, time.January, 22, 16, 49, 0, 0, Stockholm),
		}}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w.Flush()

		if got, want := b.String(), `{"time":"2017-01-22T16:49:00+01:00"}`+"\n"; got != want {
			t.Fatalf("b.String() = %q, want %q", got, want)
		}
	})

	t.Run("CEST", func(t *testing.T) {
		var b bytes.Buffer

		w := bufio.NewWriter(&b)

		v := struct {
			Time *Time `json:"time"`
		}{&Time{
			time.Date(2017, time.May, 22, 16, 49, 0, 0, Stockholm),
		}}

		if err := json.NewEncoder(w).Encode(v); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		w.Flush()

		if got, want := b.String(), `{"time":"2017-05-22T16:49:00+02:00"}`+"\n"; got != want {
			t.Fatalf("b.String() = %q, want %q", got, want)
		}
	})
}
