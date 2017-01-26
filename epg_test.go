package epg

import "testing"

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
