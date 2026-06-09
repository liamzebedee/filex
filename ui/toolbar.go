package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/core"
)

// Toolbar contains navigation buttons, breadcrumb/path bar, search, and
// view toggle for one tab. Its widgets are controlled: Render writes
// state into them, and the syncing flag stops those writes from echoing
// back as user events.
type Toolbar struct {
	Tab *Tab
	Box *gtk.Box

	BackBtn    *gtk.Button
	ForwardBtn *gtk.Button

	PathStack     *gtk.Stack
	BreadcrumbBox *gtk.Box
	PathEntry     *gtk.Entry

	SearchEntry *gtk.SearchEntry

	ListBtn *gtk.ToggleButton
	IconBtn *gtk.ToggleButton

	syncing   bool   // true while Render writes widget state
	crumbPath string // path the breadcrumb was last built for
}

func NewToolbar(tab *Tab) *Toolbar {
	tb := &Toolbar{Tab: tab}
	var err error

	tb.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 3)
	if err != nil {
		log.Fatal(err)
	}
	sc, _ := tb.Box.GetStyleContext()
	sc.AddClass("toolbar-box")

	// Back / Forward — compact nav buttons
	tb.BackBtn = navButton("go-previous-symbolic", "Back")
	tb.ForwardBtn = navButton("go-next-symbolic", "Forward")

	tb.BackBtn.Connect("clicked", func() { tab.GoBack() })
	tb.ForwardBtn.Connect("clicked", func() { tab.GoForward() })

	tb.Box.PackStart(tb.BackBtn, false, false, 0)
	tb.Box.PackStart(tb.ForwardBtn, false, false, 0)

	// Path Stack: breadcrumb vs entry
	tb.PathStack, _ = gtk.StackNew()
	tb.PathStack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_NONE)
	tb.PathStack.SetTransitionDuration(0)

	// Breadcrumb bar wrapped in EventBox — click empty space to edit path
	tb.BreadcrumbBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	bcSc, _ := tb.BreadcrumbBox.GetStyleContext()
	bcSc.AddClass("breadcrumb-bar")

	breadcrumbEvent, _ := gtk.EventBoxNew()
	breadcrumbEvent.Add(tb.BreadcrumbBox)
	breadcrumbEvent.Connect("button-press-event", func(eb *gtk.EventBox, event *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(event)
		if btnEvent.Button() == gdk.BUTTON_PRIMARY {
			tb.ShowPathEntry()
			return true
		}
		return false
	})

	// Path text entry (hidden by default)
	tb.PathEntry, _ = gtk.EntryNew()
	entrySc, _ := tb.PathEntry.GetStyleContext()
	entrySc.AddClass("path-entry")
	tb.PathEntry.Connect("activate", func() {
		text, _ := tb.PathEntry.GetText()
		text = strings.TrimSpace(text)
		tb.ShowBreadcrumb()
		if text == "" {
			return
		}
		if strings.HasPrefix(text, "~") {
			home, _ := os.UserHomeDir()
			text = filepath.Join(home, text[1:])
		}
		// Unreadable or invalid paths are rejected by commit with a
		// statusbar message; no need to validate here.
		tab.NavigateTo(text)
	})
	tb.PathEntry.Connect("key-press-event", func(entry *gtk.Entry, event *gdk.Event) bool {
		keyEvent := gdk.EventKeyNewFromEvent(event)
		if keyEvent.KeyVal() == gdk.KEY_Escape {
			tb.ShowBreadcrumb()
			return true
		}
		return false
	})
	tb.PathEntry.Connect("focus-out-event", func(entry *gtk.Entry, event *gdk.Event) bool {
		tb.ShowBreadcrumb()
		return false
	})

	tb.PathStack.AddNamed(breadcrumbEvent, "breadcrumb")
	tb.PathStack.AddNamed(tb.PathEntry, "entry")
	tb.PathStack.SetVisibleChildName("breadcrumb")

	tb.Box.PackStart(tb.PathStack, true, true, 0)

	// Right-side icon buttons: search, view toggles
	tb.SearchEntry, _ = gtk.SearchEntryNew()
	searchSc, _ := tb.SearchEntry.GetStyleContext()
	searchSc.AddClass("search-entry")
	tb.SearchEntry.SetPlaceholderText("Search…")
	tb.SearchEntry.Connect("search-changed", func() {
		if tb.syncing {
			return
		}
		query, _ := tb.SearchEntry.GetText()
		tab.SetQuery(query)
	})
	tb.SearchEntry.Connect("key-press-event", func(e *gtk.SearchEntry, event *gdk.Event) bool {
		if gdk.EventKeyNewFromEvent(event).KeyVal() == gdk.KEY_Escape {
			tab.SetQuery("") // Render clears the entry text
			return true
		}
		return false
	})
	tb.Box.PackStart(tb.SearchEntry, false, false, 0)

	// View toggle
	viewBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	viewSc, _ := viewBox.GetStyleContext()
	viewSc.AddClass("view-toggle")

	// Initial active state is established by the first Render.
	tb.ListBtn, _ = gtk.ToggleButtonNew()
	listImg, _ := gtk.ImageNewFromIconName("view-list-symbolic", gtk.ICON_SIZE_MENU)
	tb.ListBtn.SetImage(listImg)
	tb.ListBtn.SetTooltipText("List view")

	tb.IconBtn, _ = gtk.ToggleButtonNew()
	iconImg, _ := gtk.ImageNewFromIconName("view-grid-symbolic", gtk.ICON_SIZE_MENU)
	tb.IconBtn.SetImage(iconImg)
	tb.IconBtn.SetTooltipText("Icon view")

	tb.ListBtn.Connect("toggled", func() { tb.onViewToggle(tb.ListBtn, core.ListMode) })
	tb.IconBtn.Connect("toggled", func() { tb.onViewToggle(tb.IconBtn, core.IconMode) })

	viewBox.PackStart(tb.ListBtn, false, false, 0)
	viewBox.PackStart(tb.IconBtn, false, false, 0)
	tb.Box.PackStart(viewBox, false, false, 0)

	return tb
}

// onViewToggle dispatches a view-mode change; commit re-renders, which
// snaps both toggle buttons to the (possibly unchanged) state — including
// when the user clicks the already-active toggle off.
func (tb *Toolbar) onViewToggle(btn *gtk.ToggleButton, mode core.ViewMode) {
	if tb.syncing {
		return
	}
	if btn.GetActive() {
		tb.Tab.SetViewMode(mode)
	} else {
		tb.Tab.Render()
	}
}

// Render writes the tab state into the toolbar widgets.
func (tb *Toolbar) Render(s core.TabState) {
	tb.syncing = true
	defer func() { tb.syncing = false }()

	tb.BackBtn.SetSensitive(s.History.CanBack())
	tb.ForwardBtn.SetSensitive(s.History.CanForward())
	tb.ListBtn.SetActive(s.ViewMode == core.ListMode)
	tb.IconBtn.SetActive(s.ViewMode == core.IconMode)

	if cur, _ := tb.SearchEntry.GetText(); cur != s.Query {
		tb.SearchEntry.SetText(s.Query)
	}

	if tb.crumbPath != s.Path() {
		tb.renderBreadcrumb(s.Path())
		tb.crumbPath = s.Path()
	}
}

// ShowPathEntry switches to the path entry for editing.
func (tb *Toolbar) ShowPathEntry() {
	tb.PathEntry.SetText(tb.Tab.State.Path())
	tb.PathStack.SetVisibleChildName("entry")
	tb.PathEntry.GrabFocus()
	tb.PathEntry.SelectRegion(0, -1)
}

// ShowBreadcrumb switches back to the breadcrumb display.
func (tb *Toolbar) ShowBreadcrumb() {
	tb.PathStack.SetVisibleChildName("breadcrumb")
}

func (tb *Toolbar) renderBreadcrumb(path string) {
	// Clear old children
	tb.BreadcrumbBox.GetChildren().Foreach(func(item interface{}) {
		if w, ok := item.(*gtk.Widget); ok {
			tb.BreadcrumbBox.Remove(w)
		}
	})

	// Root crumb: a home icon navigating to "/"
	rootBtn, _ := gtk.ButtonNew()
	img, _ := gtk.ImageNewFromIconName("go-home-symbolic", gtk.ICON_SIZE_MENU)
	rootBtn.SetImage(img)
	tb.addCrumb(rootBtn, "/")

	accumulated := "/"
	for _, part := range strings.Split(path, "/") {
		if part == "" {
			continue
		}
		accumulated = filepath.Join(accumulated, part)

		sep, _ := gtk.LabelNew("/")
		sepSc, _ := sep.GetStyleContext()
		sepSc.AddClass("breadcrumb-sep")
		tb.BreadcrumbBox.PackStart(sep, false, false, 0)

		btn, _ := gtk.ButtonNewWithLabel(part)
		tb.addCrumb(btn, accumulated)
	}

	tb.BreadcrumbBox.ShowAll()
}

// addCrumb styles a breadcrumb button and wires it to navigate to target.
func (tb *Toolbar) addCrumb(btn *gtk.Button, target string) {
	btnSc, _ := btn.GetStyleContext()
	btnSc.AddClass("breadcrumb-btn")
	btn.SetRelief(gtk.RELIEF_NONE)
	btn.Connect("clicked", func() {
		tb.Tab.NavigateTo(target)
	})
	tb.BreadcrumbBox.PackStart(btn, false, false, 0)
}

func navButton(iconName, tooltip string) *gtk.Button {
	btn, _ := gtk.ButtonNew()
	img, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_MENU)
	btn.SetImage(img)
	btn.SetTooltipText(tooltip)
	btn.SetRelief(gtk.RELIEF_NONE)
	sc, _ := btn.GetStyleContext()
	sc.AddClass("nav-btn")
	return btn
}
