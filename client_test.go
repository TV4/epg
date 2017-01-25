package epg

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	r, err := c.Get(context.Background(), Sweden, Swedish, Date(2017, 1, 25), nil)
	if err != nil {
		t.Fatalf("unexpected error %#v", err)
	}

	if got, want := len(r.Days), 1; got != want {
		t.Fatalf("r.Days = %d, want %d", got, want)
	}

	if got, want := len(r.Days[0].Channels), 42; got != want {
		t.Fatalf("r.Days[0].Channels = %d, want %d", got, want)
	}
}
