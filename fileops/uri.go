package fileops

import "strings"

// ParseURIList parses a text/uri-list string into individual URIs.
// Per RFC 2483: lines starting with # are comments, URIs are \r\n separated.
func ParseURIList(data string) []string {
	var uris []string
	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimRight(line, "\r")
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		uris = append(uris, line)
	}
	return uris
}

// URIToPath converts a file:// URI to a local path.
func URIToPath(uri string) string {
	if len(uri) > 7 && uri[:7] == "file://" {
		return uri[7:]
	}
	return uri
}
