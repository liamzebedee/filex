package bookmarks

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestManager(t *testing.T) *BookmarkManager {
	t.Helper()
	tmp := t.TempDir()
	bm := &BookmarkManager{
		filePath: filepath.Join(tmp, "bookmarks.txt"),
	}
	// Use a minimal set of defaults for testing
	bm.defaults = []Bookmark{
		{Name: "Home", Path: "/home/test", Icon: "user-home"},
		{Name: "Root", Path: "/", Icon: "drive-harddisk"},
	}
	return bm
}

func TestAll_DefaultsOnly(t *testing.T) {
	bm := newTestManager(t)
	all := bm.All()
	if len(all) != 2 {
		t.Errorf("All() with no user bookmarks = %d items, want 2", len(all))
	}
}

func TestAdd(t *testing.T) {
	bm := newTestManager(t)

	bm.Add(Bookmark{Name: "Projects", Path: "/home/test/projects", Icon: "folder", UserAdded: true})
	all := bm.All()
	if len(all) != 3 {
		t.Errorf("All() after Add = %d items, want 3", len(all))
	}
	if all[2].Name != "Projects" {
		t.Errorf("Last bookmark name = %q, want 'Projects'", all[2].Name)
	}
}

func TestAdd_NoDuplicate(t *testing.T) {
	bm := newTestManager(t)

	bm.Add(Bookmark{Name: "Projects", Path: "/home/test/projects", Icon: "folder", UserAdded: true})
	bm.Add(Bookmark{Name: "Projects Again", Path: "/home/test/projects", Icon: "folder", UserAdded: true})

	all := bm.All()
	if len(all) != 3 {
		t.Errorf("All() after duplicate Add = %d items, want 3 (duplicate should be ignored)", len(all))
	}
}

func TestAdd_NoDuplicateDefault(t *testing.T) {
	bm := newTestManager(t)

	// Try to add a bookmark with the same path as a default
	bm.Add(Bookmark{Name: "MyHome", Path: "/home/test", Icon: "folder", UserAdded: true})

	all := bm.All()
	if len(all) != 2 {
		t.Errorf("All() after adding default-duplicate = %d items, want 2", len(all))
	}
}

func TestRemove(t *testing.T) {
	bm := newTestManager(t)

	bm.Add(Bookmark{Name: "A", Path: "/tmp/a", Icon: "folder", UserAdded: true})
	bm.Add(Bookmark{Name: "B", Path: "/tmp/b", Icon: "folder", UserAdded: true})

	bm.Remove("/tmp/a")

	all := bm.All()
	if len(all) != 3 { // 2 defaults + 1 remaining user
		t.Errorf("All() after Remove = %d items, want 3", len(all))
	}

	for _, b := range all {
		if b.Path == "/tmp/a" {
			t.Error("Removed bookmark should not appear in All()")
		}
	}
}

func TestPersistence(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "bookmarks.txt")

	// Create and save
	bm1 := &BookmarkManager{
		filePath: filePath,
		defaults: []Bookmark{{Name: "Root", Path: "/", Icon: "drive-harddisk"}},
	}
	bm1.Add(Bookmark{Name: "Saved", Path: "/tmp/saved", Icon: "folder", UserAdded: true})

	// Load in a new manager
	bm2 := &BookmarkManager{
		filePath: filePath,
		defaults: []Bookmark{{Name: "Root", Path: "/", Icon: "drive-harddisk"}},
	}
	bm2.loadUser()

	all := bm2.All()
	if len(all) != 2 {
		t.Errorf("Loaded manager All() = %d items, want 2", len(all))
	}

	found := false
	for _, b := range all {
		if b.Path == "/tmp/saved" && b.Name == "Saved" && b.UserAdded {
			found = true
		}
	}
	if !found {
		t.Error("Persisted bookmark not found after reload")
	}
}

func TestPersistence_Format(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "bookmarks.txt")

	bm := &BookmarkManager{
		filePath: filePath,
		defaults: []Bookmark{},
	}
	bm.Add(Bookmark{Name: "My Folder", Path: "/home/user/my folder", Icon: "folder", UserAdded: true})

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if content != "/home/user/my folder|My Folder\n" {
		t.Errorf("Bookmark file content = %q, want '/home/user/my folder|My Folder\\n'", content)
	}
}
