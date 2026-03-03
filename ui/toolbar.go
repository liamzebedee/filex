package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
)

// Toolbar contains navigation buttons, breadcrumb/path bar, search, and view toggle.
type Toolbar struct {
	App *App
	Box *gtk.Box

	BackBtn    *gtk.Button
	ForwardBtn *gtk.Button
	UpBtn      *gtk.Button

	PathStack     *gtk.Stack
	BreadcrumbBox *gtk.Box
	PathEntry     *gtk.Entry

	SearchEntry *gtk.SearchEntry

	ListBtn *gtk.ToggleButton
	IconBtn *gtk.ToggleButton
}

func NewToolbar(app *App) *Toolbar {
	tb := &Toolbar{App: app}
	var err error

	tb.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)
	if err != nil {
		log.Fatal(err)
	}
	sc, _ := tb.Box.GetStyleContext()
	sc.AddClass("toolbar-box")

	// Navigation buttons
	tb.BackBtn = toolButton("go-previous-symbolic", "Back")
	tb.ForwardBtn = toolButton("go-next-symbolic", "Forward")
	tb.UpBtn = toolButton("go-up-symbolic", "Up")

	tb.BackBtn.Connect("clicked", func() {
		if tab := app.ActiveTab(); tab != nil {
			tab.GoBack()
		}
	})
	tb.ForwardBtn.Connect("clicked", func() {
		if tab := app.ActiveTab(); tab != nil {
			tab.GoForward()
		}
	})
	tb.UpBtn.Connect("clicked", func() {
		if tab := app.ActiveTab(); tab != nil {
			tab.GoUp()
		}
	})

	navBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	navBox.PackStart(tb.BackBtn, false, false, 0)
	navBox.PackStart(tb.ForwardBtn, false, false, 0)
	navBox.PackStart(tb.UpBtn, false, false, 0)
	tb.Box.PackStart(navBox, false, false, 0)

	// Path Stack: breadcrumb vs entry
	tb.PathStack, _ = gtk.StackNew()
	tb.PathStack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_CROSSFADE)
	tb.PathStack.SetTransitionDuration(100)

	// Breadcrumb bar
	tb.BreadcrumbBox, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	bcSc, _ := tb.BreadcrumbBox.GetStyleContext()
	bcSc.AddClass("breadcrumb-bar")

	// Make breadcrumb box clickable to switch to entry mode
	breadcrumbEvent, _ := gtk.EventBoxNew()
	breadcrumbEvent.Add(tb.BreadcrumbBox)
	breadcrumbEvent.Connect("button-press-event", func(eb *gtk.EventBox, event *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(event)
		if btnEvent.Button() == gdk.BUTTON_PRIMARY && btnEvent.Type() == gdk.EVENT_DOUBLE_BUTTON_PRESS {
			tb.ShowPathEntry()
			return true
		}
		return false
	})

	// Path entry
	tb.PathEntry, _ = gtk.EntryNew()
	entrySc, _ := tb.PathEntry.GetStyleContext()
	entrySc.AddClass("path-entry")
	tb.PathEntry.Connect("activate", func() {
		text, _ := tb.PathEntry.GetText()
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		// Expand ~
		if strings.HasPrefix(text, "~") {
			home, _ := os.UserHomeDir()
			text = filepath.Join(home, text[1:])
		}
		info, err := os.Stat(text)
		if err == nil && info.IsDir() {
			if tab := app.ActiveTab(); tab != nil {
				tab.NavigateAndPush(text)
			}
		}
		tb.ShowBreadcrumb()
	})
	tb.PathEntry.Connect("key-press-event", func(entry *gtk.Entry, event *gdk.Event) bool {
		keyEvent := gdk.EventKeyNewFromEvent(event)
		if keyEvent.KeyVal() == gdk.KEY_Escape {
			tb.ShowBreadcrumb()
			return true
		}
		return false
	})

	tb.PathStack.AddNamed(breadcrumbEvent, "breadcrumb")
	tb.PathStack.AddNamed(tb.PathEntry, "entry")
	tb.PathStack.SetVisibleChildName("breadcrumb")

	tb.Box.PackStart(tb.PathStack, true, true, 4)

	// Search entry
	tb.SearchEntry, _ = gtk.SearchEntryNew()
	searchSc, _ := tb.SearchEntry.GetStyleContext()
	searchSc.AddClass("search-entry")
	tb.SearchEntry.SetPlaceholderText("Search…")
	tb.SearchEntry.Connect("search-changed", func() {
		query, _ := tb.SearchEntry.GetText()
		if tab := app.ActiveTab(); tab != nil {
			tab.FileView.SearchFilter(query)
		}
	})
	tb.Box.PackStart(tb.SearchEntry, false, false, 0)

	// View toggle
	viewBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	viewSc, _ := viewBox.GetStyleContext()
	viewSc.AddClass("view-toggle")

	tb.ListBtn, _ = gtk.ToggleButtonNew()
	listImg, _ := gtk.ImageNewFromIconName("view-list-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	tb.ListBtn.SetImage(listImg)
	tb.ListBtn.SetActive(true)
	tb.ListBtn.SetTooltipText("List view")

	tb.IconBtn, _ = gtk.ToggleButtonNew()
	iconImg, _ := gtk.ImageNewFromIconName("view-grid-symbolic", gtk.ICON_SIZE_SMALL_TOOLBAR)
	tb.IconBtn.SetImage(iconImg)
	tb.IconBtn.SetTooltipText("Icon view")

	tb.ListBtn.Connect("toggled", func() {
		if tb.ListBtn.GetActive() {
			tb.IconBtn.SetActive(false)
			if tab := app.ActiveTab(); tab != nil {
				tab.SetViewMode(ListMode)
			}
		}
	})
	tb.IconBtn.Connect("toggled", func() {
		if tb.IconBtn.GetActive() {
			tb.ListBtn.SetActive(false)
			if tab := app.ActiveTab(); tab != nil {
				tab.SetViewMode(IconMode)
			}
		}
	})

	viewBox.PackStart(tb.ListBtn, false, false, 0)
	viewBox.PackStart(tb.IconBtn, false, false, 0)
	tb.Box.PackStart(viewBox, false, false, 0)

	return tb
}

// UpdateForTab refreshes toolbar state for the given tab.
func (tb *Toolbar) UpdateForTab(tab *Tab) {
	tb.BackBtn.SetSensitive(tab.CanGoBack())
	tb.ForwardBtn.SetSensitive(tab.CanGoForward())

	parent := filepath.Dir(tab.Path)
	tb.UpBtn.SetSensitive(parent != tab.Path)

	tb.updateBreadcrumb(tab.Path)

	// Update view toggle
	tb.ListBtn.SetActive(tab.ViewMode == ListMode)
	tb.IconBtn.SetActive(tab.ViewMode == IconMode)
}

// ShowPathEntry switches to the path entry for editing.
func (tb *Toolbar) ShowPathEntry() {
	if tab := tb.App.ActiveTab(); tab != nil {
		tb.PathEntry.SetText(tab.Path)
	}
	tb.PathStack.SetVisibleChildName("entry")
	tb.PathEntry.GrabFocus()
	tb.PathEntry.SelectRegion(0, -1)
}

// ShowBreadcrumb switches back to the breadcrumb display.
func (tb *Toolbar) ShowBreadcrumb() {
	tb.PathStack.SetVisibleChildName("breadcrumb")
}

func (tb *Toolbar) updateBreadcrumb(path string) {
	// Remove old children
	tb.BreadcrumbBox.GetChildren().Foreach(func(item interface{}) {
		if w, ok := item.(*gtk.Widget); ok {
			tb.BreadcrumbBox.Remove(w)
		}
	})

	parts := strings.Split(path, "/")
	accumulated := "/"

	for i, part := range parts {
		if part == "" && i > 0 {
			continue
		}

		if i > 0 {
			accumulated = filepath.Join(accumulated, part)
		}
		displayName := part
		if displayName == "" {
			displayName = "/"
		}

		// Separator
		if i > 0 {
			sep, _ := gtk.LabelNew("›")
			sepSc, _ := sep.GetStyleContext()
			sepSc.AddClass("breadcrumb-sep")
			tb.BreadcrumbBox.PackStart(sep, false, false, 0)
		}

		// Button
		btn, _ := gtk.ButtonNewWithLabel(displayName)
		btnSc, _ := btn.GetStyleContext()
		btnSc.AddClass("breadcrumb-btn")
		btn.SetRelief(gtk.RELIEF_NONE)

		targetPath := accumulated
		btn.Connect("clicked", func() {
			if tab := tb.App.ActiveTab(); tab != nil {
				tab.NavigateAndPush(targetPath)
			}
		})
		tb.BreadcrumbBox.PackStart(btn, false, false, 0)
	}

	tb.BreadcrumbBox.ShowAll()
}

func toolButton(iconName, tooltip string) *gtk.Button {
	btn, _ := gtk.ButtonNew()
	img, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_SMALL_TOOLBAR)
	btn.SetImage(img)
	btn.SetTooltipText(tooltip)
	btn.SetRelief(gtk.RELIEF_NONE)
	return btn
}

// SearchFilter on FileView filters displayed results
func (fv *FileView) SearchFilter(query string) {
	if query == "" {
		fv.Refresh()
		return
	}
	entries := fileops.FilterEntries(fv.Entries, query)
	fv.populateStore(entries)
}
