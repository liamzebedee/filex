package util

import (
	"testing"
	"time"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{0, "0 bytes"},
		{1, "1 bytes"},
		{512, "512 bytes"},
		{1023, "1023 bytes"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
		{1610612736, "1.5 GB"},
		{1099511627776, "1.0 TB"},
	}
	for _, tt := range tests {
		got := FormatSize(tt.size)
		if got != tt.want {
			t.Errorf("FormatSize(%d) = %q, want %q", tt.size, got, tt.want)
		}
	}
}

func TestFormatDate_Today(t *testing.T) {
	now := time.Now()
	result := FormatDate(now)
	if len(result) < 5 || result[:5] != "Today" {
		t.Errorf("FormatDate(now) = %q, expected to start with 'Today'", result)
	}
}

func TestFormatDate_ThisYear(t *testing.T) {
	now := time.Now()
	// Pick a date this year but not today
	d := time.Date(now.Year(), 1, 1, 10, 30, 0, 0, time.Local)
	if d.YearDay() == now.YearDay() {
		d = d.AddDate(0, 0, 1)
	}
	result := FormatDate(d)
	if result[:5] == "Today" {
		t.Errorf("FormatDate for different day this year should not say Today: %q", result)
	}
	// Should NOT contain a year
	if len(result) > 12 {
		t.Errorf("FormatDate this year should not include year: %q", result)
	}
}

func TestFormatDate_PastYear(t *testing.T) {
	d := time.Date(2020, 6, 15, 12, 0, 0, 0, time.Local)
	result := FormatDate(d)
	if result != "Jun 15 2020" {
		t.Errorf("FormatDate(2020-06-15) = %q, want 'Jun 15 2020'", result)
	}
}

func TestMimeFor(t *testing.T) {
	tests := []struct {
		name  string
		isDir bool
		want  string
	}{
		{"folder", true, "inode/directory"},
		{"file.txt", false, "text/plain"},
		{"photo.png", false, "image/png"},
		{"photo.PNG", false, "image/png"}, // extension match is case-insensitive
		{"video.mp4", false, "video/mp4"},
		{"archive.zip", false, "application/zip"},
		{"code.go", false, "text/x-go"},
		{"noextension", false, "application/octet-stream"},
		{".hidden", false, "application/octet-stream"},
		{"file.unknown", false, "application/octet-stream"},
		{"document.pdf", false, "application/pdf"},
		{"song.mp3", false, "audio/mpeg"},
	}
	for _, tt := range tests {
		got := MimeFor(tt.name, tt.isDir)
		if got != tt.want {
			t.Errorf("MimeFor(%q, %v) = %q, want %q", tt.name, tt.isDir, got, tt.want)
		}
	}
}
