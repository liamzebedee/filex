package util

import "testing"

func TestIconForMime(t *testing.T) {
	tests := []struct {
		mime string
		want string
	}{
		{"inode/directory", "folder"},
		{"application/pdf", "application-pdf"},
		{"application/zip", "application-x-archive"},
		{"application/x-tar", "application-x-archive"},
		{"application/gzip", "application-x-archive"},
		{"application/json", "text-x-generic"},
		{"text/plain", "text-x-generic"},
		{"text/x-go", "text-x-generic"},
		{"image/png", "image-x-generic"},
		{"image/jpeg", "image-x-generic"},
		{"audio/mpeg", "audio-x-generic"},
		{"video/mp4", "video-x-generic"},
		{"application/octet-stream", "text-x-generic"}, // fallback
		{"something/weird", "text-x-generic"},          // fallback
	}
	for _, tt := range tests {
		got := IconForMime(tt.mime)
		if got != tt.want {
			t.Errorf("IconForMime(%q) = %q, want %q", tt.mime, got, tt.want)
		}
	}
}
