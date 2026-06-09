// Package core holds the pure data structures, algorithms and state
// transitions of the file manager. Nothing here touches GTK or the
// filesystem: every function maps values to values, so the package is
// fully unit-testable without a display.
//
// The model is React-like: TabState is the single source of truth for a
// tab, transitions return new states instead of mutating, and the UI is
// rendered as a function of (state, entries) via the Visible selector.
package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileEntry describes a single directory entry. It is a plain comparable
// value, safe to copy and share.
type FileEntry struct {
	Name    string
	Path    string
	Size    int64
	ModTime int64 // unix seconds
	IsDir   bool
	Mode    os.FileMode
}

// SortKey selects which field entries are ordered by.
type SortKey int

const (
	SortByName SortKey = iota
	SortBySize
	SortByDate
	SortByType
)

// Visible derives the entries to display from the full directory listing
// and the tab state: hidden-file filter, then search filter, then sort.
// The input slice is never modified.
func Visible(entries []FileEntry, s TabState) []FileEntry {
	visible := make([]FileEntry, 0, len(entries))
	query := strings.ToLower(s.Query)
	for _, e := range entries {
		if !s.ShowHidden && strings.HasPrefix(e.Name, ".") {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(e.Name), query) {
			continue
		}
		visible = append(visible, e)
	}
	SortEntries(visible, s.SortKey, s.SortAsc)
	return visible
}

// SortEntries orders entries in place by the given key. Directories always
// sort before files regardless of direction; ties fall back to name so the
// order is deterministic.
func SortEntries(entries []FileEntry, key SortKey, asc bool) {
	sort.SliceStable(entries, func(i, j int) bool {
		a, b := entries[i], entries[j]
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		if !asc {
			a, b = b, a
		}
		return entryLess(a, b, key)
	})
}

func entryLess(a, b FileEntry, key SortKey) bool {
	switch key {
	case SortBySize:
		if a.Size != b.Size {
			return a.Size < b.Size
		}
	case SortByDate:
		if a.ModTime != b.ModTime {
			return a.ModTime < b.ModTime
		}
	case SortByType:
		extA := strings.ToLower(filepath.Ext(a.Name))
		extB := strings.ToLower(filepath.Ext(b.Name))
		if extA != extB {
			return extA < extB
		}
	}
	return strings.ToLower(a.Name) < strings.ToLower(b.Name)
}
