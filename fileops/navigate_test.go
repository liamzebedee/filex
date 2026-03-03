package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDirectory(t *testing.T) {
	tmp := t.TempDir()

	// Create test files and dirs
	os.Mkdir(filepath.Join(tmp, "aDir"), 0755)
	os.Mkdir(filepath.Join(tmp, ".hidden"), 0755)
	os.WriteFile(filepath.Join(tmp, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmp, ".dotfile"), []byte("secret"), 0644)

	// Without hidden
	entries, err := ListDirectory(tmp, false)
	if err != nil {
		t.Fatal(err)
	}
	names := entryNames(entries)
	if contains(names, ".hidden") || contains(names, ".dotfile") {
		t.Errorf("ListDirectory(showHidden=false) should not include hidden files, got %v", names)
	}
	if !contains(names, "aDir") || !contains(names, "file.txt") {
		t.Errorf("ListDirectory(showHidden=false) should include visible files, got %v", names)
	}

	// With hidden
	entries, err = ListDirectory(tmp, true)
	if err != nil {
		t.Fatal(err)
	}
	names = entryNames(entries)
	if !contains(names, ".hidden") || !contains(names, ".dotfile") {
		t.Errorf("ListDirectory(showHidden=true) should include hidden files, got %v", names)
	}
}

func TestListDirectory_NotExist(t *testing.T) {
	_, err := ListDirectory("/nonexistent/path/xyz", false)
	if err == nil {
		t.Error("ListDirectory on nonexistent path should return error")
	}
}

func TestSortEntries_ByName(t *testing.T) {
	entries := []FileEntry{
		{Name: "charlie", IsDir: false},
		{Name: "alpha", IsDir: false},
		{Name: "bravo", IsDir: false},
		{Name: "DirA", IsDir: true},
	}

	SortEntries(entries, SortByName, true)

	// Directories first
	if !entries[0].IsDir {
		t.Error("Directories should come first")
	}
	// Then alphabetical
	if entries[1].Name != "alpha" {
		t.Errorf("Expected 'alpha' second, got %q", entries[1].Name)
	}
	if entries[2].Name != "bravo" {
		t.Errorf("Expected 'bravo' third, got %q", entries[2].Name)
	}
	if entries[3].Name != "charlie" {
		t.Errorf("Expected 'charlie' fourth, got %q", entries[3].Name)
	}
}

func TestSortEntries_ByName_Descending(t *testing.T) {
	entries := []FileEntry{
		{Name: "alpha", IsDir: false},
		{Name: "charlie", IsDir: false},
		{Name: "bravo", IsDir: false},
	}

	SortEntries(entries, SortByName, false)

	if entries[0].Name != "charlie" {
		t.Errorf("Descending: expected 'charlie' first, got %q", entries[0].Name)
	}
}

func TestSortEntries_BySize(t *testing.T) {
	entries := []FileEntry{
		{Name: "big", Size: 1000, IsDir: false},
		{Name: "small", Size: 10, IsDir: false},
		{Name: "medium", Size: 500, IsDir: false},
	}

	SortEntries(entries, SortBySize, true)

	if entries[0].Name != "small" {
		t.Errorf("Expected 'small' first by size, got %q", entries[0].Name)
	}
	if entries[2].Name != "big" {
		t.Errorf("Expected 'big' last by size, got %q", entries[2].Name)
	}
}

func TestSortEntries_ByDate(t *testing.T) {
	entries := []FileEntry{
		{Name: "new", ModTime: 300, IsDir: false},
		{Name: "old", ModTime: 100, IsDir: false},
		{Name: "mid", ModTime: 200, IsDir: false},
	}

	SortEntries(entries, SortByDate, true)

	if entries[0].Name != "old" {
		t.Errorf("Expected 'old' first by date, got %q", entries[0].Name)
	}
}

func TestSortEntries_DirsAlwaysFirst(t *testing.T) {
	entries := []FileEntry{
		{Name: "zfile", IsDir: false},
		{Name: "adir", IsDir: true},
	}

	SortEntries(entries, SortByName, true)

	if !entries[0].IsDir {
		t.Error("Directory should come first regardless of name")
	}
}

func TestFilterEntries(t *testing.T) {
	entries := []FileEntry{
		{Name: "README.md"},
		{Name: "main.go"},
		{Name: "readme.txt"},
		{Name: "config.yaml"},
	}

	// Case-insensitive match
	result := FilterEntries(entries, "readme")
	if len(result) != 2 {
		t.Errorf("FilterEntries('readme') expected 2 results, got %d", len(result))
	}

	// Empty query returns all
	result = FilterEntries(entries, "")
	if len(result) != 4 {
		t.Errorf("FilterEntries('') expected 4 results, got %d", len(result))
	}

	// No match
	result = FilterEntries(entries, "nonexistent")
	if len(result) != 0 {
		t.Errorf("FilterEntries('nonexistent') expected 0 results, got %d", len(result))
	}
}

func entryNames(entries []FileEntry) []string {
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	return names
}

func contains(s []string, item string) bool {
	for _, v := range s {
		if v == item {
			return true
		}
	}
	return false
}
