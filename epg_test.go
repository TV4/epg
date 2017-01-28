package epg

import (
	"fmt"
	"testing"
)

func TestChannelID(t *testing.T) {
	for _, tt := range []struct {
		name string
		want string
	}{
		{"UnknownChannel", ""},
		{"TV4", TV4},
		{"TV12", TV12},
		{"SVT1", SVT1},
		{"CanalSportSweden", CanalSportSweden},
		{"CMoreFotbollHockeyKids", CMoreFotbollHockeyKids},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChannelID(tt.name); got != tt.want {
				t.Fatalf("ChannelID(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestResponseDay(t *testing.T) {
	d1 := Day{
		BroadcastDate: "2017-01-01T00:00:00",
		Channels: []Channel{
			{ID: "1"},
		},
	}

	d2 := Day{
		BroadcastDate: "2017-01-02T00:00:00",
		Channels: []Channel{
			{ID: "1"},
			{ID: "2"},
		},
	}

	for _, tt := range []struct {
		r     *Response
		dates []string
		day   Day
	}{
		{&Response{}, nil, Day{}},
		{&Response{Days: []Day{d1}}, nil, d1},
		{&Response{Days: []Day{d1, d2}}, nil, d1},
		{&Response{Days: []Day{d1, d2}}, []string{Date(2017, 1, 1)}, d1},
		{&Response{Days: []Day{d1, d2}}, []string{Date(2017, 1, 2)}, d2},
		{&Response{Days: []Day{d1}}, []string{Date(2017, 1, 2)}, Day{}},
	} {
		day := tt.r.Day(tt.dates...)

		if got, want := len(day.Channels), len(tt.day.Channels); got != want {
			t.Fatalf("len(day.Channels) = %d, want %d", got, want)
		}

		if got, want := day.BroadcastDate, tt.day.BroadcastDate; got != want {
			t.Fatalf("day.BroadcastDate = %q, want %q", got, want)
		}
	}
}

func TestDayChannel(t *testing.T) {
	for _, tt := range []struct {
		day  Day
		id   string
		want string
	}{
		{Day{}, "", ""},
		{Day{Channels: []Channel{{ID: "1"}}}, "1", "1"},
		{Day{Channels: []Channel{{ID: "1"}}}, "2", ""},
		{Day{Channels: []Channel{{ID: "1"}, {ID: "2"}}}, "2", "2"},
	} {
		t.Run(tt.id, func(t *testing.T) {
			c := tt.day.Channel(tt.id)

			if got, want := c.ID, tt.want; got != want {
				t.Fatalf("c.ID = %q, want %q", got, want)
			}
		})
	}
}

func TestImageURL(t *testing.T) {
	for _, tt := range []struct {
		image  Image
		format string
		want   string
	}{
		{Image{ID: "123"}, "456", "https://img-cdn-cmore.b17g.services/123/456.img"},
		{Image{ID: "456"}, "789", "https://img-cdn-cmore.b17g.services/456/789.img"},
	} {
		t.Run(tt.image.ID, func(t *testing.T) {
			if got := tt.image.URL(tt.format).String(); got != tt.want {
				t.Fatalf("tt.image.URL(%q).String() = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

func TestNames(t *testing.T) {
	for _, tt := range []struct {
		program Program
		count   int
		last    string
	}{
		{Program{ID: "38253", Actors: "August Diehl,Sara Hjort Ditlevsen, Jo Adrian Haavind"}, 3, "Jo Adrian Haavind"},
	} {
		t.Run(tt.program.ID, func(t *testing.T) {
			names := Names(tt.program.Actors)

			if got, want := len(names), tt.count; got != want {
				t.Fatalf("len(%#v) = %d, want %d", names, got, want)
			}

			if got, want := names[len(names)-1], tt.last; got != want {
				t.Fatalf("names[len(names)-1] = %q, want %q", got, want)
			}
		})
	}
}

func ExampleNames() {
	fmt.Printf("%#v\n", Names("August Diehl,Sara Hjort Ditlevsen, Jo Adrian Haavind"))
	// Output: []string{"August Diehl", "Sara Hjort Ditlevsen", "Jo Adrian Haavind"}
}
