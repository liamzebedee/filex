// Package util holds small pure formatting and classification helpers.
package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// FormatSize returns a human-readable file size string.
func FormatSize(size int64) string {
	const (
		KB = int64(1024)
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case size >= TB:
		return fmt.Sprintf("%.1f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}

// FormatDate returns a human-readable date string.
func FormatDate(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("Today 15:04")
	}
	if t.Year() == now.Year() {
		return t.Format("Jan 02 15:04")
	}
	return t.Format("Jan 02 2006")
}

// MimeFor guesses a mime type from a file name and kind. The extension
// lookup is case-insensitive.
func MimeFor(name string, isDir bool) string {
	if isDir {
		return "inode/directory"
	}
	if mime, ok := extMime[strings.ToLower(filepath.Ext(name))]; ok {
		return mime
	}
	return "application/octet-stream"
}

var extMime = map[string]string{
	".txt":  "text/plain",
	".md":   "text/markdown",
	".go":   "text/x-go",
	".py":   "text/x-python",
	".js":   "text/javascript",
	".ts":   "text/typescript",
	".html": "text/html",
	".css":  "text/css",
	".json": "application/json",
	".xml":  "application/xml",
	".yaml": "application/yaml",
	".yml":  "application/yaml",
	".toml": "application/toml",
	".sh":   "application/x-shellscript",
	".bash": "application/x-shellscript",
	".c":    "text/x-c",
	".h":    "text/x-c",
	".cpp":  "text/x-c++",
	".rs":   "text/x-rust",
	".java": "text/x-java",
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".ico":  "image/x-icon",
	".mp3":  "audio/mpeg",
	".wav":  "audio/wav",
	".flac": "audio/flac",
	".ogg":  "audio/ogg",
	".mp4":  "video/mp4",
	".mkv":  "video/x-matroska",
	".avi":  "video/x-msvideo",
	".webm": "video/webm",
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".zip":  "application/zip",
	".tar":  "application/x-tar",
	".gz":   "application/gzip",
	".bz2":  "application/x-bzip2",
	".xz":   "application/x-xz",
	".7z":   "application/x-7z-compressed",
	".rar":  "application/x-rar-compressed",
	".deb":  "application/x-deb",
	".rpm":  "application/x-rpm",
	".iso":  "application/x-iso9660-image",
	".exe":  "application/x-executable",
	".so":   "application/x-sharedlib",
}
