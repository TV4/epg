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
