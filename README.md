# epg

Go client for the [C More EPG Web API](http://api.cmore.se/).

## Usage examples

**Drama on CMoreStarsHD**

```go
package main

import (
	"context"
	"encoding/json"
	"net/url"
	"os"

	epg "github.com/TV4/epg"
)

func main() {
	c := epg.NewClient()

	if r, err := c.GetChannel(
		context.Background(),
		epg.Sweden,
		epg.Swedish,
		epg.Date(2017, 1, 26),
		epg.Date(2017, 1, 28),
		epg.CMoreStarsHD,
		url.Values{
			"genre": {"Drama"},
		},
	); err == nil {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")
		enc.Encode(r)
	}
}
```

**Primetime movies in Sweden 2017-01-29**

```go
package main

import (
	"context"
	"encoding/json"
	"net/url"
	"os"

	epg "github.com/TV4/epg"
)

func main() {
	r, err := epg.NewClient().Get(
		context.Background(),
		epg.Sweden,
		epg.Swedish,
		epg.Date(2017, 1, 29),
		url.Values{
			"filter": {"primetimemovies"},
		},
	)
	if err != nil {
		return
	}

	programs := map[string]epg.Program{}

	for _, d := range r.Days {
		for _, c := range d.Channels {
			for _, s := range c.Schedules {
				programs[s.Program.ProgramID] = s
			}
		}
	}

	var movies []movie

	for _, p := range programs {
		var cover string

		for _, m := range p.Images {
			if m.Category == "Cover" {
				cover = m.URL("164").String()
			}
		}

		movies = append(movies, movie{
			Title: p.Title,
			Facts: p.SynopsisFacts,
			Cover: cover,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	enc.Encode(movies)
}

type movie struct {
	Title string `json:"title"`
	Facts string `json:"facts,omitempty"`
	Cover string `json:"cover,omitempty"`
}
```

## API documentation

<http://api.cmore.se/>
