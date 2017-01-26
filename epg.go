package epg

import (
	"errors"
	"net/url"
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

var (
	// ErrNotFound means that the resource could not be found
	ErrNotFound = errors.New("not found")

	// ErrUnknown means that an unexpected error occurred
	ErrUnknown = errors.New("unknown error")
)

// Response data from the EPG API
type Response struct {
	Days      []Day  `xml:"Day"`
	FromDate  string `xml:"FromDate,attr" json:"from_date,omitempty"`
	UntilDate string `xml:"UntilDate,attr" json:"until_date,omitempty"`
	Meta      *Meta  `xml:"-" json:"meta,omitempty"`
}

// Meta is a type used for request/response metadata
type Meta map[string]interface{}

// Day is an EPG day
type Day struct {
	BroadcastDate string    `xml:"BroadcastDate,attr" json:"broadcast_date"`
	Channels      []Channel `xml:"Channel" json:"channels"`
}

// Channel is a TV channel in the EPG
type Channel struct {
	ChannelID   int        `xml:"ChannelId,attr" json:"channel_id"`
	Name        string     `xml:"Name,attr" json:"name"`
	Title       string     `xml:"Title,attr" json:"title"`
	LogoID      string     `xml:"LogoId,attr" json:"logo_id"`
	LogoDarkID  string     `xml:"LogoDarkId,attr" json:"logo_dark_id"`
	LogoLightID string     `xml:"LogoLightID,attr" json:"logo_light_id"`
	IsHD        bool       `xml:"IsHd,attr" json:"hd"`
	Schedules   []Schedule `xml:"Schedule" json:"schedules"`
}

// Schedule is the TV program schedule of a channel in the EPG
type Schedule struct {
	ScheduleID        string  `xml:"ScheduleId,attr" json:"schedule_id"`
	NextStart         string  `xml:NextStart",attr" json:"next_start"`
	CalendarDate      string  `xml:"CalendarDate,attr" json:"calendar_date"`
	IsPremiere        bool    `xml:"IsPremiere,attr" json:"premiere"`
	IsDubbed          bool    `xml:"IsDubbed,attr" json:"dubbed"`
	Type              string  `xml:"Type,attr" json:"type"`
	AlsoAvailableInHD bool    `xml:"AlsoAvailableInHD,attr" json:"also_available_in_hd"`
	AlsoAvailableIn3D bool    `xml:"AlsoAvailableIn3D,attr" json:"also_available_in_3d"`
	Is3D              bool    `xml:"Is3D" json:"is_3d"`
	IsPPV             bool    `xml:"IsPPV" json:"is_ppv"`
	PlayAssetID       string  `xml:"PlayAssetId1,attr" json:"play_asset_id"`
	Program           Program `xml:"Program" json:"program"`
}

// Program is the program that is scheduled in the EPG
type Program struct {
	ProgramID                string  `xml:"ProgramId,attr" json:"program_id"`
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
