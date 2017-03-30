/*

Package epg contains a client for the C More EPG Web API

Installation

Just go get the package:

    go get -u github.com/TV4/epg

Usage

A small usage example

      package main

      import (
      	"context"
      	"fmt"
      	"time"

      	epg "github.com/TV4/epg"
      )

      func main() {
      	var (
      		ec   = epg.NewClient()
      		ctx  = context.Background()
      		date = epg.DateAtTime(time.Now())
      	)

      	if r, err := ec.Get(ctx, epg.Sweden, epg.Swedish, date); err == nil {
      		c := r.Day().Channel(epg.TV4)

      		for _, s := range c.Schedules {
      			fmt.Println(s.CalendarDate, s.Program.Title)
      		}
      	}
      }

Documentation

http://api.cmore.se/

*/
package epg

import (
	"errors"
	"net/url"
	"strings"
)

// Country is the type used for lowercase ISO 3166-1 alpha-2 country codes
// as per https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
type Country string

const (
	// Sweden is the country code se
	Sweden Country = "se"

	// Norway is the country code no
	Norway Country = "no"

	// Denmark is the country code dk
	Denmark Country = "dk"

	// Finland is the country code fi
	Finland Country = "fi"
)

// Language is the type used for ISO 639-1 language codes
// as per https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
type Language string

const (
	// Swedish is the language code sv
	Swedish Language = "sv"

	// Norwegian is the language code no
	Norwegian Language = "no"

	// Danish is the language code da
	Danish Language = "da"

	// Finnish is the language code fi
	Finnish Language = "fi"
)

// Channel constants
//
// Data retrieved like this:
//
//     curl -H "Accept: application/xml" "https://api.cmore.se/epg/se/sv/2017-01-26/2017-02-13" | xmllint --format - |
//     grep ChannelId | awk -F '"' '{print $2 " " $4 " = \"" $2 "\""}' | sort -n | uniq | awk '{print $2 " " $3 " " $4}' | pbcopy
//
const (
	CanalExtra1            = "3"
	CanalExtra2            = "4"
	CanalExtra3            = "5"
	CanalExtraHD           = "7"
	CanalFilm1             = "8"
	CanalFilm2             = "9"
	CanalHD                = "12"
	CanalPlusHD            = "17"
	CanalPlusHitsHD        = "18"
	CanalSport3            = "22"
	CanalSportFotboll      = "25"
	CanalSportHockey       = "26"
	CanalSportSweden       = "28"
	CF4                    = "29"
	SFK                    = "32"
	SFKBoxer               = "33"
	SHD                    = "34"
	SeriesHD               = "52"
	CMoreFotbollHockeyKids = "54"
	CMoreLive2HD           = "65"
	CMoreLive3HD           = "66"
	CMoreLive4HD           = "67"
	CMoreHockeyHD          = "68"
	CMoreGolfHD            = "70"
	CMoreGolfDenmarkHD     = "71"
	SVT1                   = "74"
	SVT2                   = "75"
	TV4                    = "76"
	TV4Sport               = "78"
	Sjuan                  = "79"
	TV12                   = "80"
	TV4FaktaXL             = "81"
	TV4Fakta               = "82"
	TV4Film                = "83"
	TV4Guld                = "84"
	TV4Komedi              = "85"
	SVT24                  = "86"
	SVTKunskapskanalen     = "87"
	Barnkanalen            = "88"
	CMoreStars             = "89"
	CMoreStarsHD           = "90"
	CMoreLive5             = "91"
	CMoreLive5HD           = "92"
)

var channels = map[string]string{
	"CanalExtra1":            CanalExtra1,
	"CanalExtra2":            CanalExtra2,
	"CanalExtra3":            CanalExtra3,
	"CanalExtraHD":           CanalExtraHD,
	"CanalFilm1":             CanalFilm1,
	"CanalFilm2":             CanalFilm2,
	"CanalHD":                CanalHD,
	"CanalPlusHD":            CanalPlusHD,
	"CanalPlusHitsHD":        CanalPlusHitsHD,
	"CanalSport3":            CanalSport3,
	"CanalSportFotboll":      CanalSportFotboll,
	"CanalSportHockey":       CanalSportHockey,
	"CanalSportSweden":       CanalSportSweden,
	"CF4":                    CF4,
	"SFK":                    SFK,
	"SFKBoxer":               SFKBoxer,
	"SHD":                    SHD,
	"SeriesHD":               SeriesHD,
	"CMoreFotbollHockeyKids": CMoreFotbollHockeyKids,
	"CMoreLive2HD":           CMoreLive2HD,
	"CMoreLive3HD":           CMoreLive3HD,
	"CMoreLive4HD":           CMoreLive4HD,
	"CMoreHockeyHD":          CMoreHockeyHD,
	"CMoreGolfHD":            CMoreGolfHD,
	"CMoreGolfDenmarkHD":     CMoreGolfDenmarkHD,
	"SVT1":                   SVT1,
	"SVT2":                   SVT2,
	"TV4":                    TV4,
	"TV4Sport":               TV4Sport,
	"Sjuan":                  Sjuan,
	"TV12":                   TV12,
	"TV4FaktaXL":             TV4FaktaXL,
	"TV4Fakta":               TV4Fakta,
	"TV4Film":                TV4Film,
	"TV4Guld":                TV4Guld,
	"TV4Komedi":              TV4Komedi,
	"SVT24":                  SVT24,
	"SVTKunskapskanalen":     SVTKunskapskanalen,
	"Barnkanalen":            Barnkanalen,
	"CMoreStars":             CMoreStars,
	"CMoreStarsHD":           CMoreStarsHD,
	"CMoreLive5":             CMoreLive5,
	"CMoreLive5HD":           CMoreLive5HD,
}

// ChannelID returns the channel ID based on provided channel name
func ChannelID(name string) string {
	return channels[name]
}

var (
	// ErrNotFound means that the resource could not be found
	ErrNotFound = errors.New("not found")

	// ErrUnknown means that an unexpected error occurred
	ErrUnknown = errors.New("unknown error")
)

// Response data from the EPG API
type Response struct {
	Days      []Day  `xml:"Day,omitempty" json:"days,omitempty"`
	FromDate  string `xml:"FromDate,attr" json:"from_date,omitempty"`
	UntilDate string `xml:"UntilDate,attr" json:"until_date,omitempty"`
	Meta      *Meta  `xml:"-" json:"meta,omitempty"`
}

// Day returns the first day in the response, or the (optional) provided date.
// Returns empty Day if not found
func (r *Response) Day(dates ...string) Day {
	if len(r.Days) == 0 {
		return Day{}
	}

	if len(dates) == 0 {
		return r.Days[0]
	}

	for _, d := range r.Days {
		if strings.HasPrefix(d.BroadcastDate, dates[0]) {
			return d
		}
	}

	return Day{}
}

// Meta is a type used for request/response metadata
type Meta map[string]interface{}

// Day is an EPG day
type Day struct {
	BroadcastDate string    `xml:"BroadcastDate,attr" json:"broadcast_date"`
	Channels      []Channel `xml:"Channel" json:"channels,omitempty"`
}

// Channel returns the channel with the given id.
// Returns empty Channel if not found
func (d Day) Channel(id string) Channel {
	for _, c := range d.Channels {
		if c.ID == id {
			return c
		}
	}

	return Channel{}
}

// Channel is a TV channel in the EPG
type Channel struct {
	ID          string     `xml:"ChannelId,attr" json:"channel_id"`
	Name        string     `xml:"Name,attr" json:"name"`
	Title       string     `xml:"Title,attr" json:"title"`
	LogoID      string     `xml:"LogoId,attr" json:"logo_id"`
	LogoDarkID  string     `xml:"LogoDarkId,attr" json:"logo_dark_id"`
	LogoLightID string     `xml:"LogoLightId,attr" json:"logo_light_id"`
	IsHD        bool       `xml:"IsHd,attr" json:"hd"`
	Schedules   []Schedule `xml:"Schedule" json:"schedules,omitempty"`
}

// Schedule is the TV program schedule of a channel in the EPG
type Schedule struct {
	ID                string  `xml:"ScheduleId,attr" json:"schedule_id"`
	NextStart         string  `xml:"NextStart,attr" json:"next_start"`
	CalendarDate      string  `xml:"CalendarDate,attr" json:"calendar_date"`
	IsPremiere        bool    `xml:"IsPremiere,attr" json:"premiere"`
	IsDubbed          bool    `xml:"IsDubbed,attr" json:"dubbed"`
	Type              string  `xml:"Type,attr" json:"type"`
	AlsoAvailableInHD bool    `xml:"AlsoAvailableInHD,attr" json:"also_available_in_hd"`
	AlsoAvailableIn3D bool    `xml:"AlsoAvailableIn3D,attr" json:"also_available_in_3d"`
	Is3D              bool    `xml:"Is3D,attr" json:"is_3d"`
	IsPPV             bool    `xml:"IsPPV,attr" json:"is_ppv"`
	PlayAssetID       string  `xml:"PlayAssetId1,attr" json:"play_asset_id"`
	Program           Program `xml:"Program" json:"program"`
}

// Program is the program that is scheduled in the EPG
type Program struct {
	ID                       string  `xml:"ProgramId,attr" json:"program_id"`
	Title                    string  `xml:"Title,attr" json:"title"`
	OriginalTitle            string  `xml:"OriginalTitle,attr" json:"original_title"`
	Genre                    string  `xml:"Genre,attr" json:"genre"`
	GenreKey                 string  `xml:"GenreKey,attr" json:"genre_key"`
	FirstCalendarDate        string  `xml:"FirstCalendarDate,attr" json:"first_calendar_date"`
	LastCalendarDate         string  `xml:"LastCalendarDate,attr" json:"last_calendar_date"`
	VodStart                 string  `xml:"VodStart,attr" json:"vod_start"`
	VodEnd                   string  `xml:"VodEnd,attr" json:"vod_end"`
	Duration                 int     `xml:"Duration,attr" json:"duration"`
	ContentSourceID          string  `xml:"ContentSourceId,attr" json:"content_source_id"`
	ProductionYear           int     `xml:"ProductionYear,attr" json:"production_year"`
	Rating                   string  `xml:"Rating,attr" json:"rating"`
	Actors                   string  `xml:"Actors,attr" json:"actors"`
	Directors                string  `xml:"Directors,attr" json:"directors"`
	Class                    string  `xml:"Class,attr" json:"class"`
	Type                     string  `xml:"Type,attr" json:"type"`
	Category                 string  `xml:"Category,attr" json:"category"`
	IsDubbedVersionAvailable bool    `xml:"IsDubbedVersionAvailable,attr" json:"dubbed_version_available"`
	VOD                      bool    `xml:"Vod,attr" json:"vod"`
	OTTBlackout              bool    `xml:"OTTBlackout,attr" json:"ott_blackout"`
	IsDubbed                 bool    `xml:"IsDubbed,attr" json:"dubbed"`
	Images                   []Image `xml:"Resources>Image" json:"images"`
	SeriesID                 string  `xml:"SeriesId,attr" json:"series_id"`
	SeasonNumber             int     `xml:"SeasonNumber,attr" json:"season_number"`
	EpisodeNumber            int     `xml:"EpisodeNumber,attr" json:"episode_number"`
	NumberOfEpisodes         int     `xml:"NumberOfEpisodes,attr" json:"number_of_episodes"`
	SynopsisExtraShort       string  `xml:"Synopsis>ExtraShort" json:"extra_short"`
	SynopsisShort            string  `xml:"Synopsis>Short" json:"short"`
	SynopsisMedium           string  `xml:"Synopsis>Medium" json:"medium"`
	SynopsisLong             string  `xml:"Synopsis>Long" json:"long"`
	SynopsisFacts            string  `xml:"Synopsis>Facts" json:"facts"`
}

// ImageBaseURL is the base URL for images
var ImageBaseURL = &url.URL{Scheme: "https", Host: "img-cdn-cmore.b17g.services"}

// Image is a typed identifier for an image that can be retrieved at
// https://img-cdn-cmore.b17g.services/:id/:format.img
//
// (format 164 can be used to retrieve the full size image)
//
type Image struct {
	ID       string `xml:"Id,attr" json:"id"`
	Category string `xml:"Category,attr" json:"category"`
}

// URL returns an *url.URL based on the ImageBaseURL, image ID and provided format
func (m Image) URL(format string) *url.URL {
	return ImageBaseURL.ResolveReference(&url.URL{Path: "/" + m.ID + "/" + format + ".img"})
}

// Names takes a string of comma separated names, splits them into a slice, trims any space around each name
func Names(s string) []string {
	var names []string

	for _, n := range strings.Split(s, ",") {
		names = append(names, strings.TrimSpace(n))
	}

	return names
}
