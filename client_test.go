package epg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, ok := NewClient().(*client)

	if !ok {
		t.Fatalf("expected *client")
	}

	if got, want := c.httpClient.Timeout, 20*time.Second; got != want {
		t.Fatalf("c.httpClient.Timeout = %s, want %s", got, want)
	}

	if got, want := c.baseURL.String(), "https://api.cmore.se"; got != want {
		t.Fatalf("c.baseURL.String() = %q, want %q", got, want)
	}

	if got, want := c.userAgent, "epg/client.go (https://github.com/TV4/epg)"; got != want {
		t.Fatalf("c.userAgent = %q, want %q", got, want)
	}
}

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

func TestHTTPClient(t *testing.T) {
	hc := &http.Client{Timeout: 5 * time.Second}

	c, ok := NewClient(HTTPClient(hc)).(*client)

	if !ok {
		t.Fatalf("expected *client")
	}

	if got, want := c.httpClient.Timeout, 5*time.Second; got != want {
		t.Fatalf("c.httpClient.Timeout = %s, want %s", got, want)
	}
}

func TestBaseURL(t *testing.T) {
	rawurl := "http://example.com/"

	c, ok := NewClient(BaseURL(rawurl)).(*client)

	if !ok {
		t.Fatalf("expected *client")
	}

	if got, want := c.baseURL.String(), rawurl; got != want {
		t.Fatalf("c.baseURL.String() = %q, want %q", got, want)
	}
}

func TestUserAgent(t *testing.T) {
	ua := "Test-Agent"

	c, ok := NewClient(UserAgent(ua)).(*client)

	if !ok {
		t.Fatalf("expected *client")
	}

	if got, want := c.userAgent, ua; got != want {
		t.Fatalf("c.userAgent = %q, want %q", got, want)
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
