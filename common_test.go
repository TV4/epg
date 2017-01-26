package epg

import (
	"net/http"
	"net/http/httptest"
)

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
