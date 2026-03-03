package util

import "strings"

// IconForMime returns a GTK icon name for a given mime type.
func IconForMime(mime string) string {
	// Direct matches
	iconMap := map[string]string{
		"inode/directory":              "folder",
		"application/pdf":              "application-pdf",
		"application/zip":              "application-x-archive",
		"application/x-tar":            "application-x-archive",
		"application/gzip":             "application-x-archive",
		"application/x-bzip2":          "application-x-archive",
		"application/x-xz":             "application-x-archive",
		"application/x-7z-compressed":  "application-x-archive",
		"application/x-rar-compressed": "application-x-archive",
		"application/x-deb":            "application-x-deb",
		"application/x-rpm":            "application-x-rpm",
		"application/x-iso9660-image":  "application-x-cd-image",
		"application/x-executable":     "application-x-executable",
		"application/x-sharedlib":      "application-x-sharedlib",
		"application/x-shellscript":    "text-x-script",
		"application/json":             "text-x-generic",
		"application/xml":              "text-xml",
	}

	if icon, ok := iconMap[mime]; ok {
		return icon
	}

	// Category matches
	parts := strings.SplitN(mime, "/", 2)
	if len(parts) == 2 {
		switch parts[0] {
		case "text":
			return "text-x-generic"
		case "image":
			return "image-x-generic"
		case "audio":
			return "audio-x-generic"
		case "video":
			return "video-x-generic"
		}
	}

	return "text-x-generic"
}
