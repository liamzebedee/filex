package core

import "testing"

func TestHistory_PushBackForward(t *testing.T) {
	h := NewHistory("/a")
	h = h.Push("/b")
	h = h.Push("/c")

	if h.Current() != "/c" {
		t.Fatalf("Current = %q, want /c", h.Current())
	}
	if !h.CanBack() || h.CanForward() {
		t.Fatal("at the top of the stack: CanBack should be true, CanForward false")
	}

	h = h.Back()
	if h.Current() != "/b" || !h.CanForward() {
		t.Fatalf("after Back: Current = %q, CanForward = %v", h.Current(), h.CanForward())
	}

	h = h.Forward()
	if h.Current() != "/c" {
		t.Fatalf("after Forward: Current = %q, want /c", h.Current())
	}
}

func TestHistory_PushTruncatesForward(t *testing.T) {
	h := NewHistory("/a").Push("/b").Push("/c").Back().Back() // at /a
	h = h.Push("/d")

	if h.Current() != "/d" || h.CanForward() {
		t.Fatal("Push must drop forward entries")
	}
	h = h.Back()
	if h.Current() != "/a" {
		t.Fatalf("after Back: Current = %q, want /a", h.Current())
	}
}

func TestHistory_PushCurrentIsNoOp(t *testing.T) {
	h := NewHistory("/a").Push("/a")
	if len(h.Stack) != 1 {
		t.Fatalf("pushing the current path must not grow the stack, got %v", h.Stack)
	}
}

func TestHistory_BackForwardAtBoundsAreNoOps(t *testing.T) {
	h := NewHistory("/a").Back()
	if h.Current() != "/a" {
		t.Fatal("Back at the start must be a no-op")
	}
	h = h.Forward()
	if h.Current() != "/a" {
		t.Fatal("Forward at the end must be a no-op")
	}
}

func TestHistory_PushDoesNotAliasBackingArray(t *testing.T) {
	h := NewHistory("/a").Push("/b").Push("/c").Back() // at /b, /c still forward
	h1 := h.Push("/x")
	h2 := h.Push("/y")
	if h1.Current() != "/x" || h2.Current() != "/y" {
		t.Fatalf("branched histories must be independent: %q, %q", h1.Current(), h2.Current())
	}
}

func TestTabState_PathDerivedFromHistory(t *testing.T) {
	s := NewTabState("/home")
	if s.Path() != "/home" {
		t.Fatalf("Path = %q, want /home", s.Path())
	}
	s = s.Navigate("/tmp")
	if s.Path() != "/tmp" || !s.History.CanBack() {
		t.Fatal("Navigate must push and move the current path")
	}
}

func TestTabState_LocationChangeClearsQuery(t *testing.T) {
	s := NewTabState("/a").WithQuery("foo")

	if s.Navigate("/b").Query != "" {
		t.Error("Navigate must clear the query")
	}

	s2 := NewTabState("/a").Navigate("/b").WithQuery("foo")
	if s2.Back().Query != "" {
		t.Error("Back must clear the query")
	}
	s3 := s2.Back().WithQuery("bar")
	if s3.Forward().Query != "" {
		t.Error("Forward must clear the query")
	}
}

func TestTabState_Up(t *testing.T) {
	s := NewTabState("/a/b")
	s = s.Up()
	if s.Path() != "/a" {
		t.Fatalf("Up: Path = %q, want /a", s.Path())
	}
	root := NewTabState("/").Up()
	if root.Path() != "/" || len(root.History.Stack) != 1 {
		t.Fatal("Up at root must be a no-op")
	}
}

func TestTabState_WithSortTogglesDirection(t *testing.T) {
	s := NewTabState("/")
	s = s.WithSort(SortBySize)
	if s.SortKey != SortBySize || !s.SortAsc {
		t.Fatal("selecting a new key must sort ascending")
	}
	s = s.WithSort(SortBySize)
	if s.SortAsc {
		t.Fatal("selecting the active key must flip direction")
	}
}
