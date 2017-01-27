package epg

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client for EPG Web API
type Client interface {
	Get(ctx context.Context, country Country, language Language, date string, attributes ...url.Values) (*Response, error)
	GetPeriod(ctx context.Context, country Country, language Language, fromDate, toDate string, attributes ...url.Values) (*Response, error)
	GetChannelGroup(ctx context.Context, country Country, language Language, fromDate, toDate, channelGroup string, attributes ...url.Values) (*Response, error)
	GetChannel(ctx context.Context, country Country, language Language, fromDate, toDate, channelID string, attributes ...url.Values) (*Response, error)
}

type client struct {
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string
}

// NewClient creates an EPG Client
func NewClient(options ...func(*client)) Client {
	c := &client{
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "api.cmore.se",
		},
		userAgent: "epg/client.go (https://github.com/TV4/epg)",
	}

	for _, f := range options {
		f(c)
	}

	return c
}

// HTTPClient changes the *client HTTP client to the provided *http.Client
func HTTPClient(hc *http.Client) func(*client) {
	return func(c *client) {
		c.httpClient = hc
	}
}

// BaseURL changes the *client base URL based on the provided rawurl
func BaseURL(rawurl string) func(*client) {
	return func(c *client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// UserAgent changes the User-Agent used in requests sent by the *client
func UserAgent(ua string) func(*client) {
	return func(c *client) {
		c.userAgent = ua
	}
}

// Date formats a year, month, day into the format yyyy-mm-dd
func Date(year int, month time.Month, day int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func (c *client) Get(ctx context.Context, country Country, language Language, date string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getPath(country, language, date), c.query(attributes))
}

func (c *client) GetPeriod(ctx context.Context, country Country, language Language, fromDate, toDate string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getPeriodPath(country, language, fromDate, toDate), c.query(attributes))
}

func (c *client) GetChannelGroup(ctx context.Context, country Country, language Language, fromDate, toDate, channelGroup string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getChannelGroupPath(country, language, fromDate, toDate, channelGroup), c.query(attributes))
}

func (c *client) GetChannel(ctx context.Context, country Country, language Language, fromDate, toDate, channelID string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getChannelGroupPath(country, language, fromDate, toDate, channelID), c.query(attributes))
}

func (c *client) getPath(country Country, language Language, date string) string {
	return fmt.Sprintf("/epg/%s/%s/%s", country, language, date)
}

func (c *client) getPeriodPath(country Country, language Language, fromDate, toDate string) string {
	return fmt.Sprintf("/epg/%s/%s/%s/%s", country, language, fromDate, toDate)
}

func (c *client) getChannelGroupPath(country Country, language Language, fromDate, toDate, channelGroup string) string {
	return fmt.Sprintf("/epg/%s/%s/%s/%s/%s", country, language, fromDate, toDate, channelGroup)
}

func (c *client) getChannelPath(country Country, language Language, fromDate, toDate, channelID string) string {
	return fmt.Sprintf("/epg/%s/%s/%s/%s/%s", country, language, fromDate, toDate, channelID)
}

func (c *client) get(ctx context.Context, path string, query url.Values) (*Response, error) {
	req, err := c.request(ctx, path, query)
	if err != nil {
		return nil, err
	}

	r, err := c.do(req)
	if err != nil {
		return nil, err
	}

	r.Meta = &Meta{
		"path":  path,
		"query": query,
	}

	return r, nil
}

func (c *client) query(attributes []url.Values) url.Values {
	if len(attributes) > 0 {
		return attributes[0]
	}

	return url.Values{}
}

func (c *client) request(ctx context.Context, path string, query url.Values) (*http.Request, error) {
	rel, err := url.Parse(path + "?" + query.Encode())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", c.baseURL.ResolveReference(rel).String(), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Add("Accept", "application/xml")
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

func (c *client) do(req *http.Request) (*Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.CopyN(ioutil.Discard, resp.Body, 64)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, ErrNotFound
		default:
			return nil, ErrUnknown
		}
	}

	var r Response

	if err := xml.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
