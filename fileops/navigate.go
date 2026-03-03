package fileops

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Sort column constants
const (
	SortByName = iota
	SortBySize
	SortByDate
	SortByType
)

// FileEntry holds info about a single file for display.
type FileEntry struct {
	Name    string
	Path    string
	Size    int64
	ModTime int64
	IsDir   bool
	Mode    os.FileMode
	Mime    string
}

// ListDirectory reads and returns the entries in a directory.
func ListDirectory(dirPath string, showHidden bool) ([]FileEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var result []FileEntry
	for _, e := range entries {
		name := e.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		result = append(result, FileEntry{
			Name:    name,
			Path:    filepath.Join(dirPath, name),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
			IsDir:   info.IsDir(),
			Mode:    info.Mode(),
		})
	}

	return result, nil
}

// SortEntries sorts file entries. Directories always come first.
func SortEntries(entries []FileEntry, sortCol int, ascending bool) {
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		// Directories first
		if a.IsDir != b.IsDir {
			return a.IsDir
		}

		var less bool
		switch sortCol {
		case SortByName:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case SortBySize:
			less = a.Size < b.Size
		case SortByDate:
			less = a.ModTime < b.ModTime
		case SortByType:
			extA := filepath.Ext(a.Name)
			extB := filepath.Ext(b.Name)
			if extA == extB {
				less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
			} else {
				less = strings.ToLower(extA) < strings.ToLower(extB)
			}
		default:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		}

		if !ascending {
			return !less
		}
		return less
	})
}

// FilterEntries filters entries by a search query (case-insensitive name match).
func FilterEntries(entries []FileEntry, query string) []FileEntry {
	if query == "" {
		return entries
	}
	q := strings.ToLower(query)
	var result []FileEntry
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Name), q) {
			result = append(result, e)
		}
	}
	return result
}
