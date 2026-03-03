package util

import (
	"os"
	"strings"
	"testing"
	"time"
)

func setEnglishLocale(t *testing.T) {
	t.Helper()
	t.Setenv("FILEX_LANG", "en")
}

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
	}
	for _, tt := range tests {
		got := FormatSize(tt.size)
		if got != tt.want {
			t.Errorf("FormatSize(%d) = %q, want %q", tt.size, got, tt.want)
		}
	}
}

func TestFormatDate_Today(t *testing.T) {
	setEnglishLocale(t)
	now := time.Now()
	result := FormatDate(now)
	if len(result) < 5 || result[:5] != "Today" {
		t.Errorf("FormatDate(now) = %q, expected to start with 'Today'", result)
	}
}

func TestFormatDate_ThisYear(t *testing.T) {
	setEnglishLocale(t)
	now := time.Now()
	// Pick a date this year but not today
	d := time.Date(now.Year(), 1, 1, 10, 30, 0, 0, time.Local)
	if d.YearDay() == now.YearDay() {
		d = d.AddDate(0, 0, 1)
	}
	result := FormatDate(d)
	if result[:3] != "Jan" && len(result) > 0 {
		// It should be "Mon DD HH:MM" format, not "Today"
		if result[:5] == "Today" {
			t.Errorf("FormatDate for different day this year should not say Today: %q", result)
		}
	}
	// Should NOT contain a year
	if len(result) > 12 {
		t.Errorf("FormatDate this year should not include year: %q", result)
	}
}

func TestFormatDate_PastYear(t *testing.T) {
	setEnglishLocale(t)
	d := time.Date(2020, 6, 15, 12, 0, 0, 0, time.Local)
	result := FormatDate(d)
	if result != "Jun 15 2020" {
		t.Errorf("FormatDate(2020-06-15) = %q, want 'Jun 15 2020'", result)
	}
}

func TestFormatDate_TodayChinese(t *testing.T) {
	t.Setenv("FILEX_LANG", "zh")
	now := time.Now()
	result := FormatDate(now)
	if !strings.HasPrefix(result, "今天 ") {
		t.Errorf("FormatDate(now) in zh = %q, expected prefix %q", result, "今天 ")
	}
}

func TestMimeFromName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"file.txt", "text/plain"},
		{"photo.png", "image/png"},
		{"photo.PNG", "application/octet-stream"}, // case sensitive
		{"video.mp4", "video/mp4"},
		{"archive.zip", "application/zip"},
		{"code.go", "text/x-go"},
		{"noextension", "application/octet-stream"},
		{".hidden", "application/octet-stream"},
		{"file.unknown", "application/octet-stream"},
		{"document.pdf", "application/pdf"},
		{"song.mp3", "audio/mpeg"},
	}
	for _, tt := range tests {
		got := MimeFromName(tt.name)
		if got != tt.want {
			t.Errorf("MimeFromName(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

type mockFileInfo struct {
	name  string
	isDir bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestDetectMimeType_Dir(t *testing.T) {
	info := mockFileInfo{name: "folder", isDir: true}
	got := DetectMimeType(info)
	if got != "inode/directory" {
		t.Errorf("DetectMimeType(dir) = %q, want 'inode/directory'", got)
	}
}

func TestDetectMimeType_File(t *testing.T) {
	info := mockFileInfo{name: "test.go", isDir: false}
	got := DetectMimeType(info)
	if got != "text/x-go" {
		t.Errorf("DetectMimeType(test.go) = %q, want 'text/x-go'", got)
	}
}
