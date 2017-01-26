package epg

// Channel constants
//
// Data retrieved like this:
//
// curl -H "Accept: application/xml" "https://api.cmore.se/epg/se/sv/2017-01-26/2017-02-13" | xmllint --format - |
// grep ChannelId | awk -F '"' '{print $2 " " $4 " = \"" $2 "\""}' | sort -n | uniq | awk '{print $2 " " $3 " " $4}' | pbcopy
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
