package epg

import "testing"

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
