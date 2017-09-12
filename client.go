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

// Client for the EPG Web API
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	userAgent  string
}

// NewClient creates an EPG Client
func NewClient(options ...func(*Client)) *Client {
	c := &Client{
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
func HTTPClient(hc *http.Client) func(*Client) {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// BaseURL changes the *client base URL based on the provided rawurl
func BaseURL(rawurl string) func(*Client) {
	return func(c *Client) {
		if u, err := url.Parse(rawurl); err == nil {
			c.baseURL = u
		}
	}
}

// UserAgent changes the User-Agent used in requests sent by the *client
func UserAgent(ua string) func(*Client) {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// Date formats a year, month, day into the format yyyy-mm-dd
func Date(year int, month time.Month, day int) string {
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// DateAtTime returns the date string for the provided time.Time
func DateAtTime(t time.Time) string {
	return Date(t.Date())
}

// Get retrieves a response for the given country, language and date
func (c *Client) Get(ctx context.Context, country Country, language Language, date string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getPath(country, language, date), c.query(attributes))
}

// GetPeriod retrieves the response for the period fromDate until toDate
func (c *Client) GetPeriod(ctx context.Context, country Country, language Language, fromDate, toDate string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getPeriodPath(country, language, fromDate, toDate), c.query(attributes))
}

// GetChannelGroup retrieves the channel group in the period fromDate until toDate
func (c *Client) GetChannelGroup(ctx context.Context, country Country, language Language, fromDate, toDate, channelGroup string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getChannelGroupPath(country, language, fromDate, toDate, channelGroup), c.query(attributes))
}

// GetChannel retrieves a channel in the period fromDate until toDate
func (c *Client) GetChannel(ctx context.Context, country Country, language Language, fromDate, toDate, channelID string, attributes ...url.Values) (*Response, error) {
	return c.get(ctx, c.getChannelPath(country, language, fromDate, toDate, channelID), c.query(attributes))
}

func (c *Client) getPath(country Country, language Language, date string) string {
	return fmt.Sprintf("/epg/%s/%s/%s", country, language, date)
}

func (c *Client) getPeriodPath(country Country, language Language, fromDate, toDate string) string {
	return fmt.Sprintf("/epg/%s/%s/%s/%s", country, language, fromDate, toDate)
}

func (c *Client) getChannelGroupPath(country Country, language Language, fromDate, toDate, channelGroup string) string {
	return fmt.Sprintf("/epg/%s/%s/%s/%s/%s", country, language, fromDate, toDate, channelGroup)
}

func (c *Client) getChannelPath(country Country, language Language, fromDate, toDate, channelID string) string {
	return fmt.Sprintf("/epg/%s/%s/%s/%s/%s", country, language, fromDate, toDate, channelID)
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (*Response, error) {
	req, err := c.request(ctx, path, query)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	r, err := c.decodeResponse(resp)
	if err != nil {
		return nil, err
	}

	r.Meta = &Meta{
		"path":  path,
		"query": query,
	}

	return r, nil
}

func (c *Client) query(attributes []url.Values) url.Values {
	if len(attributes) > 0 {
		return attributes[0]
	}

	return url.Values{}
}

func (c *Client) request(ctx context.Context, path string, query url.Values) (*http.Request, error) {
	rawurl := path

	if len(query) > 0 {
		rawurl += "?" + query.Encode()
	}

	rel, err := url.Parse(rawurl)
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

func (c *Client) decodeResponse(resp *http.Response) (*Response, error) {
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
