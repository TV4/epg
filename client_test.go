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

func TestDate(t *testing.T) {
	for _, tt := range []struct {
		year  int
		month time.Month
		day   int
		want  string
	}{
		{1, time.February, 3, "0001-02-03"},
		{2009, time.November, 10, "2009-11-10"},
		{2017, time.January, 26, "2017-01-26"},
	} {
		t.Run(tt.want, func(t *testing.T) {
			if got := Date(tt.year, tt.month, tt.day); got != tt.want {
				t.Fatalf("Date(%d, %d, %d) = %q, want %q", tt.year, tt.month, tt.day, got, tt.want)
			}
		})
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
		t.Fatalf("len(r.Days) = %d, want %d", got, want)
	}

	if got, want := len(r.Days[1].Channels), 10; got != want {
		t.Fatalf("len(r.Days[0].Channels) = %d, want %d", got, want)
	}

	channel := r.Days[1].Channels[4]

	if got, want := channel.Name, "CanalFilm2"; got != want {
		t.Fatalf("channel.Name = %q, want %q", got, want)
	}
}

func TestGetChannelGroup(t *testing.T) {
	ts, c := testServerAndClient()
	defer ts.Close()

	r, err := c.GetChannelGroup(
		context.Background(),
		Sweden,
		Swedish,
		Date(2017, 1, 27),
		Date(2017, 1, 27),
		"27",
		url.Values{
			"filter": {"livesports"},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error %#v", err)
	}

	if got, want := len(r.Days), 1; got != want {
		t.Fatalf("r.Days = %d, want %d", got, want)
	}

	d := r.Day()

	if got, want := len(d.Channels), 9; got != want {
		t.Fatalf("len(d.Channels) = %d, want %d", got, want)
	}

	channel := d.Channel(CanalSportSweden)

	if got, want := channel.LogoID, "ec7d2da1-5b0d-4135-ac54-32149414c557"; got != want {
		t.Fatalf("channel.LogoID = %q, want %q", got, want)
	}
}

func TestRequest(t *testing.T) {
	ua := "Test-Request-Agent"

	c, ok := NewClient(UserAgent(ua)).(*client)

	if !ok {
		t.Fatalf("expected *client")
	}

	r, err := c.request(context.Background(), "/foo", url.Values{"bar": {"baz"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := r.Header.Get("User-Agent"), ua; got != want {
		t.Fatalf("r.Header.Get(\"User-Agent\") = %q, want %q", got, want)
	}

	if got, want := r.URL.String(), "https://api.cmore.se/foo?bar=baz"; got != want {
		t.Fatalf("r.URL.String() = %q, want %q", got, want)
	}
}

func testServerAndClient() (*httptest.Server, Client) {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/xml; charset=utf-8")

			switch r.URL.String() {
			case "/epg/se/sv/2017-01-25":
				w.Write(swedishFullDayEPGResponseXML)
			case "/epg/se/sv/2017-01-27/2017-01-27/27?filter=livesports":
				w.Write(swedishLiveSportsEPGResponseXML)
			case "/epg/dk/da/2017-01-26/2017-01-27?genre=drama":
				w.Write(danishTwoDaysDramaEPGResponseXML)
			default:
				w.Write(emptyEPGResponseXML)
			}
		}))

	return ts, NewClient(BaseURL(ts.URL))
}
