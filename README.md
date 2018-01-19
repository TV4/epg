# epg

[![Build Status](https://travis-ci.org/TV4/epg.svg?branch=master)](https://travis-ci.org/TV4/epg)
[![Go Report Card](https://goreportcard.com/badge/github.com/TV4/epg)](https://goreportcard.com/report/github.com/TV4/epg)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/TV4/epg)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/TV4/epg#license-mit)

Go client for the [C More EPG Web API](https://api.cmore.se/).

## Status

Still under active development. Expect breaking changes.

## Installation

    go get -u github.com/TV4/epg

## Usage examples

**The TV4 schedule for today**

```go
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
			fmt.Println(s.CalendarDate.Format("15:04"), s.Program.Title)
		}
	}
}
```

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

**Primetime movies in Sweden 2017-05-24**

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
		epg.Date(2017, 5, 24),
		url.Values{
			"filter": {"primetimemovies"},
		},
	)
	if err != nil {
		return
	}

	d := r.Day()

	programs := map[string]epg.Program{}

	for _, c := range d.Channels {
		for _, s := range c.Schedules {
			programs[s.Program.ID] = s.Program
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

<https://api.cmore.se/>

## License (MIT)

Copyright © 2017-2018 TV4

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the "Software"),
> to deal in the Software without restriction, including without limitation
> the rights to use, copy, modify, merge, publish, distribute, sublicense,
> and/or sell copies of the Software, and to permit persons to whom the
> Software is furnished to do so, subject to the following conditions:
>
> The above copyright notice and this permission notice shall be included
> in all copies or substantial portions of the Software.
>
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
> OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
> IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
> DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
> TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
> OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
