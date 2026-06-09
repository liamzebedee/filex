package core

import "testing"

func names(entries []FileEntry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Name
	}
	return out
}

func assertOrder(t *testing.T, entries []FileEntry, want ...string) {
	t.Helper()
	got := names(entries)
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestSortEntries_ByName(t *testing.T) {
	entries := []FileEntry{
		{Name: "charlie"},
		{Name: "alpha"},
		{Name: "bravo"},
		{Name: "DirA", IsDir: true},
	}
	SortEntries(entries, SortByName, true)
	assertOrder(t, entries, "DirA", "alpha", "bravo", "charlie")
}

func TestSortEntries_ByName_Descending(t *testing.T) {
	entries := []FileEntry{
		{Name: "alpha"},
		{Name: "charlie"},
		{Name: "bravo"},
	}
	SortEntries(entries, SortByName, false)
	assertOrder(t, entries, "charlie", "bravo", "alpha")
}

func TestSortEntries_DirsFirstEvenDescending(t *testing.T) {
	entries := []FileEntry{
		{Name: "zfile"},
		{Name: "adir", IsDir: true},
	}
	SortEntries(entries, SortByName, false)
	if !entries[0].IsDir {
		t.Error("directories should come first regardless of direction")
	}
}

func TestSortEntries_BySize(t *testing.T) {
	entries := []FileEntry{
		{Name: "big", Size: 1000},
		{Name: "small", Size: 10},
		{Name: "medium", Size: 500},
	}
	SortEntries(entries, SortBySize, true)
	assertOrder(t, entries, "small", "medium", "big")
}

func TestSortEntries_ByDate(t *testing.T) {
	entries := []FileEntry{
		{Name: "new", ModTime: 300},
		{Name: "old", ModTime: 100},
		{Name: "mid", ModTime: 200},
	}
	SortEntries(entries, SortByDate, true)
	assertOrder(t, entries, "old", "mid", "new")
}

func TestSortEntries_TiesBreakByName(t *testing.T) {
	entries := []FileEntry{
		{Name: "b", Size: 10},
		{Name: "a", Size: 10},
	}
	SortEntries(entries, SortBySize, true)
	assertOrder(t, entries, "a", "b")
}

func TestVisible_HiddenFilter(t *testing.T) {
	entries := []FileEntry{
		{Name: ".hidden"},
		{Name: "shown"},
	}
	s := NewTabState("/")

	got := Visible(entries, s)
	assertOrder(t, got, "shown")

	got = Visible(entries, s.WithHidden(true))
	assertOrder(t, got, ".hidden", "shown")
}

func TestVisible_QueryFilter(t *testing.T) {
	entries := []FileEntry{
		{Name: "README.md"},
		{Name: "main.go"},
		{Name: "readme.txt"},
		{Name: "config.yaml"},
	}
	s := NewTabState("/")

	got := Visible(entries, s.WithQuery("readme"))
	assertOrder(t, got, "README.md", "readme.txt")

	got = Visible(entries, s.WithQuery(""))
	if len(got) != 4 {
		t.Errorf("empty query should return all entries, got %d", len(got))
	}

	got = Visible(entries, s.WithQuery("nonexistent"))
	if len(got) != 0 {
		t.Errorf("non-matching query should return nothing, got %d", len(got))
	}
}

func TestVisible_DoesNotMutateInput(t *testing.T) {
	entries := []FileEntry{
		{Name: "b"},
		{Name: "a"},
	}
	Visible(entries, NewTabState("/"))
	if entries[0].Name != "b" {
		t.Error("Visible must not reorder the input slice")
	}
}
