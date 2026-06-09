package core

// History is a browser-style navigation history: a stack of paths plus a
// cursor. The current location is always Stack[Index] — it is derived,
// never stored anywhere else, so it cannot drift out of sync. Methods have
// value semantics: they return a new History and never mutate the
// receiver's backing array.
type History struct {
	Stack []string
	Index int
}

func NewHistory(start string) History {
	return History{Stack: []string{start}}
}

// Current returns the path the history points at.
func (h History) Current() string { return h.Stack[h.Index] }

// Push drops any forward entries and appends path. Pushing the current
// path is a no-op, so repeated navigation to the same place does not grow
// the stack.
func (h History) Push(path string) History {
	if path == h.Current() {
		return h
	}
	stack := append(append([]string{}, h.Stack[:h.Index+1]...), path)
	return History{Stack: stack, Index: len(stack) - 1}
}

func (h History) Back() History {
	if h.CanBack() {
		h.Index--
	}
	return h
}

func (h History) Forward() History {
	if h.CanForward() {
		h.Index++
	}
	return h
}

func (h History) CanBack() bool    { return h.Index > 0 }
func (h History) CanForward() bool { return h.Index < len(h.Stack)-1 }
