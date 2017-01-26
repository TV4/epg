package epg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGet(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	r, err := c.Get(
		context.Background(),
		Sweden,
		Swedish,
		Date(2017, 1, 25),
	)
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

func TestGetPeriod(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	r, err := c.GetPeriod(
		context.Background(),
		Denmark,
		Danish,
		Date(2017, 1, 26),
		Date(2017, 1, 27),
		url.Values{
			"genre": {"drama"},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error %#v", err)
	}

	if got, want := len(r.Days), 2; got != want {
		t.Fatalf("r.Days = %d, want %d", got, want)
	}

	if got, want := len(r.Days[1].Channels), 10; got != want {
		t.Fatalf("r.Days[0].Channels = %d, want %d", got, want)
	}

	channel := r.Days[1].Channels[4]

	if got, want := channel.Name, "CanalFilm2"; got != want {
		t.Fatalf("channel.Name = %q, want %q", got, want)
	}
}

func testServerAndClient() (*httptest.Server, Client) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")

			switch r.URL.Path {
			case "/epg/se/sv/2017-01-25":
				w.Write(swedishFullDayEPGResponseXML)
			case "/epg/dk/da/2017-01-26/2017-01-27":
				w.Write(danishTwoDaysDramaEPGResponseXML)
			default:
				w.Write(emptyEPGResponseXML)
			}
		}))

	return ts, NewClient(BaseURL(ts.URL))
}
