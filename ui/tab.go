package ui

import (
	"log"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"

	"filex/core"
	"filex/fileops"
)

// tabRegistry maps notebook page widgets (by native pointer) to their Tab.
var tabRegistry = make(map[uintptr]*Tab)

// Tab owns one notebook page. State is its single source of truth and
// Entries caches the current directory listing; everything on screen is
// rendered from those two values. Every user action funnels through
// commit: pure state transition → reload if the location changed → Render.
type Tab struct {
	App      *App
	Toolbar  *Toolbar
	FileView *FileView
	Box      *gtk.Box   // container widget added to notebook
	TabLabel *gtk.Label // the label in the tab header

	State   core.TabState
	Entries []core.FileEntry
}

func NewTab(app *App, path string) *Tab {
	tab := &Tab{App: app, State: core.NewTabState(path)}

	var err error
	tab.Box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Per-tab toolbar: back/fwd, breadcrumb, search, view toggle
	tab.Toolbar = NewToolbar(tab)
	tab.Box.PackStart(tab.Toolbar.Box, false, false, 0)

	// File view
	tab.FileView = NewFileView(tab)
	tab.Box.PackStart(tab.FileView.ScrollWin, true, true, 0)

	tabLabelBox := tab.createTabLabel()
	tab.Box.ShowAll()

	pageNum := app.Notebook.AppendPage(tab.Box, tabLabelBox)
	app.Notebook.SetCurrentPage(pageNum)
	app.Notebook.SetTabReorderable(tab.Box, true)

	// Register tab by its widget's native pointer
	tabRegistry[tab.Box.ToWidget().Native()] = tab

	tab.Refresh()

	return tab
}

// commit makes next the tab's state, loading the directory listing when
// the location changed, and re-renders. Navigation into an unreadable
// directory is rejected: the old state stays and the error is reported.
func (t *Tab) commit(next core.TabState) {
	if t.Entries == nil || next.Path() != t.State.Path() {
		entries, err := fileops.ListDirectory(next.Path())
		if err != nil {
			log.Printf("filex: %v", err)
			t.App.Statusbar.ShowMessage("Cannot open " + next.Path())
			return
		}
		t.Entries = entries
	}
	t.State = next
	t.Render()
}

// Refresh re-reads the current directory from disk and re-renders.
func (t *Tab) Refresh() {
	t.Entries = nil
	t.commit(t.State)
}

// State transitions — thin wrappers over the pure core transitions.

func (t *Tab) NavigateTo(path string)      { t.commit(t.State.Navigate(path)) }
func (t *Tab) GoBack()                     { t.commit(t.State.Back()) }
func (t *Tab) GoForward()                  { t.commit(t.State.Forward()) }
func (t *Tab) GoUp()                       { t.commit(t.State.Up()) }
func (t *Tab) ToggleHidden()               { t.commit(t.State.WithHidden(!t.State.ShowHidden)) }
func (t *Tab) SetViewMode(m core.ViewMode) { t.commit(t.State.WithViewMode(m)) }
func (t *Tab) SetQuery(q string)           { t.commit(t.State.WithQuery(q)) }
func (t *Tab) SetSort(key core.SortKey)    { t.commit(t.State.WithSort(key)) }

// visible derives the entries currently on screen.
func (t *Tab) visible() []core.FileEntry {
	return core.Visible(t.Entries, t.State)
}

// Render syncs every widget owned by the tab to the current state.
func (t *Tab) Render() {
	visible := t.visible()
	t.FileView.Render(visible, t.State)
	t.Toolbar.Render(t.State)
	t.TabLabel.SetText(tabTitle(t.State.Path()))
	if t.App.ActiveTab() == t {
		t.App.Statusbar.Render(t.State.Path(), len(visible))
	}
}

func tabTitle(path string) string {
	if path == "/" {
		return "/"
	}
	return filepath.Base(path)
}

func (t *Tab) createTabLabel() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 4)

	t.TabLabel, _ = gtk.LabelNew(tabTitle(t.State.Path()))
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
