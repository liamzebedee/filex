package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDirectory(t *testing.T) {
	tmp := t.TempDir()

	os.Mkdir(filepath.Join(tmp, "aDir"), 0755)
	os.WriteFile(filepath.Join(tmp, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmp, ".dotfile"), []byte("secret"), 0644)

	entries, err := ListDirectory(tmp)
	if err != nil {
		t.Fatal(err)
	}

	// ListDirectory returns everything, hidden files included; filtering is
	// the job of the pure core.Visible selector.
	found := map[string]bool{}
	for _, e := range entries {
		found[e.Name] = true
		if e.Path != filepath.Join(tmp, e.Name) {
			t.Errorf("entry %q has wrong path %q", e.Name, e.Path)
		}
		if e.Name == "aDir" && !e.IsDir {
			t.Error("aDir should be marked as a directory")
		}
		if e.Name == "file.txt" && e.Size != 5 {
			t.Errorf("file.txt size = %d, want 5", e.Size)
		}
	}
	for _, want := range []string{"aDir", "file.txt", ".dotfile"} {
		if !found[want] {
			t.Errorf("ListDirectory should include %q, got %v", want, found)
		}
	}
}

func TestListDirectory_NotExist(t *testing.T) {
	_, err := ListDirectory("/nonexistent/path/xyz")
	if err == nil {
		t.Error("ListDirectory on nonexistent path should return error")
	}
}
