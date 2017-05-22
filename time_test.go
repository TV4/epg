package epg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

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
