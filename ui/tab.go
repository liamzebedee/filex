package ui

import (
	"log"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
)

// tabRegistry maps widget native pointers to Tab structs.
var tabRegistry = make(map[uintptr]*Tab)

// Tab represents a single file browser tab.
type Tab struct {
	App      *App
	Toolbar  *Toolbar
	FileView *FileView
	Box      *gtk.Box   // container widget added to notebook
	TabLabel *gtk.Label // the label in the tab header

	Path         string
	History      []string
	HistoryIndex int
	ShowHidden   bool
	ViewMode     ViewMode // ListMode or IconMode
	SortColumn   int
	SortAsc      bool
}

type ViewMode int

const (
	ListMode ViewMode = iota
	IconMode
)

func NewTab(app *App, path string) *Tab {
	tab := &Tab{
		App:        app,
		Path:       path,
		History:    []string{path},
		ShowHidden: false,
		ViewMode:   ListMode,
		SortColumn: fileops.SortByName,
		SortAsc:    true,
	}

	var err error
	tab.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Per-tab toolbar: back/fwd, breadcrumb, search, view toggle
	tab.Toolbar = NewToolbar(app)
	tab.Box.PackStart(tab.Toolbar.Box, false, false, 0)

	// File view
	tab.FileView = NewFileView(tab)
	tab.Box.PackStart(tab.FileView.ScrollWin, true, true, 0)

	// Create tab label with close button
	tabLabelBox := tab.createTabLabel()
	tab.Box.ShowAll()

	pageNum := app.Notebook.AppendPage(tab.Box, tabLabelBox)
	app.Notebook.SetCurrentPage(pageNum)
	app.Notebook.SetTabReorderable(tab.Box, true)

	// Register tab by its widget's native pointer
	tabRegistry[tab.Box.ToWidget().Native()] = tab

	tab.Navigate(path)

	return tab
}

func (t *Tab) createTabLabel() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)

	t.TabLabel, _ = gtk.LabelNew(filepath.Base(t.Path))
	t.TabLabel.SetWidthChars(12)
	t.TabLabel.SetMaxWidthChars(20)
	t.TabLabel.SetEllipsize(3) // PANGO_ELLIPSIZE_END
	t.TabLabel.SetHAlign(gtk.ALIGN_START)
	box.PackStart(t.TabLabel, true, true, 0)

	closeBtn, _ := gtk.ButtonNew()
	closeImg, _ := gtk.ImageNewFromIconName("window-close-symbolic", gtk.ICON_SIZE_MENU)
	closeBtn.SetImage(closeImg)
	closeBtn.SetRelief(gtk.RELIEF_NONE)
	closeBtn.Connect("clicked", func() {
		t.Close()
	})
	box.PackStart(closeBtn, false, false, 0)
	box.ShowAll()

	return box
}

// updateTabLabel updates the tab label text to reflect the current directory.
func (t *Tab) updateTabLabel() {
	if t.TabLabel == nil {
		return
	}
	name := filepath.Base(t.Path)
	if t.Path == "/" {
		name = "/"
	}
	t.TabLabel.SetText(name)
}

// Navigate changes the tab's directory and refreshes the view.
func (t *Tab) Navigate(path string) {
	t.setBusyCursor(true)
	t.Path = path
	t.FileView.Refresh()
	t.updateTabLabel()
	if t.Toolbar != nil {
		t.Toolbar.UpdateForTab(t)
	}
	if t.App.Statusbar != nil {
		t.App.Statusbar.Update(t)
	}
	t.setBusyCursor(false)
}

// setBusyCursor sets or clears the watch cursor on the window.
func (t *Tab) setBusyCursor(busy bool) {
	win, err := t.App.Window.GetWindow()
	if err != nil || win == nil {
		return
	}
	if busy {
		display, err := gdk.DisplayGetDefault()
		if err != nil {
			return
		}
		cursor, err := gdk.CursorNewFromName(display, "wait")
		if err == nil {
			win.SetCursor(cursor)
		}
	} else {
		win.SetCursor(nil)
	}
}

// NavigateAndPush navigates and pushes to history.
func (t *Tab) NavigateAndPush(path string) {
	if t.HistoryIndex < len(t.History)-1 {
		t.History = t.History[:t.HistoryIndex+1]
	}
	t.History = append(t.History, path)
	t.HistoryIndex = len(t.History) - 1
	t.Navigate(path)
}

func (t *Tab) GoBack() {
	if t.HistoryIndex > 0 {
		t.HistoryIndex--
		t.Navigate(t.History[t.HistoryIndex])
	}
}

func (t *Tab) GoForward() {
	if t.HistoryIndex < len(t.History)-1 {
		t.HistoryIndex++
		t.Navigate(t.History[t.HistoryIndex])
	}
}

func (t *Tab) GoUp() {
	parent := filepath.Dir(t.Path)
	if parent != t.Path {
		t.NavigateAndPush(parent)
	}
}

func (t *Tab) CanGoBack() bool {
	return t.HistoryIndex > 0
}

func (t *Tab) CanGoForward() bool {
	return t.HistoryIndex < len(t.History)-1
}

func (t *Tab) ToggleHidden() {
	t.ShowHidden = !t.ShowHidden
	t.FileView.Refresh()
	if t.App.Statusbar != nil {
		t.App.Statusbar.Update(t)
	}
}

func (t *Tab) SetViewMode(mode ViewMode) {
	if t.ViewMode == mode {
		return
	}
	t.ViewMode = mode
	t.FileView.SwitchView(mode)
}

func (t *Tab) Close() {
	pageNum := t.App.Notebook.PageNum(t.Box)
	if pageNum < 0 {
		return
	}
	if t.App.Notebook.GetNPages() <= 1 {
		return
	}
	delete(tabRegistry, t.Box.ToWidget().Native())
	t.App.Notebook.RemovePage(pageNum)
}
