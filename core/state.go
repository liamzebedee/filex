package core

import "path/filepath"

// ViewMode selects how a tab displays its entries.
type ViewMode int

const (
	ListMode ViewMode = iota
	IconMode
)

// TabState is the complete UI state of one tab. It is a pure value: every
// transition returns a new TabState, and the widgets are rendered as a
// function of this state plus the directory listing.
type TabState struct {
	History    History
	ShowHidden bool
	ViewMode   ViewMode
	SortKey    SortKey
	SortAsc    bool
	Query      string
}

func NewTabState(path string) TabState {
	return TabState{History: NewHistory(path), SortAsc: true}
}

// Path returns the tab's current directory, derived from the history.
func (s TabState) Path() string { return s.History.Current() }

// Navigate pushes a new location. Changing location clears the search.
func (s TabState) Navigate(path string) TabState {
	s.History = s.History.Push(path)
	s.Query = ""
	return s
}

// Back moves to the previous location; at the start it is a no-op.
func (s TabState) Back() TabState {
	s.History = s.History.Back()
	s.Query = ""
	return s
}

// Forward moves to the next location; at the end it is a no-op.
func (s TabState) Forward() TabState {
	s.History = s.History.Forward()
	s.Query = ""
	return s
}

// Up navigates to the parent directory; at the root it is a no-op.
func (s TabState) Up() TabState {
	if parent := filepath.Dir(s.Path()); parent != s.Path() {
		return s.Navigate(parent)
	}
	return s
}

func (s TabState) WithHidden(show bool) TabState    { s.ShowHidden = show; return s }
func (s TabState) WithViewMode(m ViewMode) TabState { s.ViewMode = m; return s }
func (s TabState) WithQuery(q string) TabState      { s.Query = q; return s }

// WithSort selects the sort key; selecting the active key again flips the
// direction, matching list-header click behavior.
func (s TabState) WithSort(key SortKey) TabState {
	if s.SortKey == key {
		s.SortAsc = !s.SortAsc
	} else {
		s.SortKey = key
		s.SortAsc = true
	}
	return s
}
