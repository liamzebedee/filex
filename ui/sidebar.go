package ui

import (
	"log"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/bookmarks"
)

// Sidebar shows the places panel with bookmarks.
type Sidebar struct {
	App       *App
	ScrollWin *gtk.ScrolledWindow
	ListBox   *gtk.ListBox
	bookmarks *bookmarks.BookmarkManager
}

func NewSidebar(app *App) *Sidebar {
	s := &Sidebar{
		App:       app,
		bookmarks: bookmarks.NewBookmarkManager(),
	}

	var err error
	s.ScrollWin, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.ScrollWin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	s.ScrollWin.SetSizeRequest(170, -1)

	sc, _ := s.ScrollWin.GetStyleContext()
	sc.AddClass("sidebar")

	s.ListBox, err = gtk.ListBoxNew()
	if err != nil {
		log.Fatal(err)
	}
	s.ListBox.SetSelectionMode(gtk.SELECTION_SINGLE)

	s.ListBox.Connect("row-activated", func(lb *gtk.ListBox, row *gtk.ListBoxRow) {
		if row == nil {
			return
		}
		bm := s.getBookmarkAt(row.GetIndex())
		if bm == nil {
			return
		}
		if tab := app.ActiveTab(); tab != nil {
			tab.NavigateTo(bm.Path)
		}
	})

	// Enable drag-and-drop to add bookmarks
	s.setupDragDrop()

	s.ScrollWin.Add(s.ListBox)
	s.Render()

	return s
}

// Render rebuilds the sidebar rows from the bookmark list.
func (s *Sidebar) Render() {
	// Remove existing rows
	s.ListBox.GetChildren().Foreach(func(item interface{}) {
		if w, ok := item.(*gtk.Widget); ok {
			s.ListBox.Remove(w)
		}
	})

	// Add header (non-selectable, non-activatable)
	header, _ := gtk.LabelNew("Places")
	headerSc, _ := header.GetStyleContext()
	headerSc.AddClass("sidebar-header")
	header.SetHAlign(gtk.ALIGN_START)
	s.ListBox.Add(header)

	if headerRow := s.ListBox.GetRowAtIndex(0); headerRow != nil {
		headerRow.SetSelectable(false)
		headerRow.SetActivatable(false)
		hsc, _ := headerRow.GetStyleContext()
		hsc.AddClass("sidebar-header-row")
	}

	for _, bm := range s.bookmarks.All() {
		s.ListBox.Add(s.createBookmarkRow(bm))
	}

	s.ListBox.ShowAll()
}

func (s *Sidebar) createBookmarkRow(bm bookmarks.Bookmark) *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)

	icon, _ := gtk.ImageNewFromIconName(bm.Icon, gtk.ICON_SIZE_MENU)
	box.PackStart(icon, false, false, 0)

	label, _ := gtk.LabelNew(bm.Name)
	label.SetHAlign(gtk.ALIGN_START)
	label.SetEllipsize(3) // PANGO_ELLIPSIZE_END
	box.PackStart(label, true, true, 0)

	// Right-click to remove user bookmarks
	if bm.UserAdded {
		path := bm.Path
		box.Connect("button-press-event", func(b *gtk.Box, event *gdk.Event) bool {
			btnEvent := gdk.EventButtonNewFromEvent(event)
			if btnEvent.Button() == gdk.BUTTON_SECONDARY {
				s.showRemoveMenu(path)
				return true
			}
			return false
		})
	}

	return box
}

func (s *Sidebar) showRemoveMenu(path string) {
	menu, _ := gtk.MenuNew()

	removeItem, _ := gtk.MenuItemNewWithLabel("Remove Bookmark")
	removeItem.Connect("activate", func() {
		s.bookmarks.Remove(path)
		s.Render()
	})
	menu.Append(removeItem)
	menu.ShowAll()
	menu.PopupAtPointer((*gdk.Event)(nil))
}

func (s *Sidebar) getBookmarkAt(idx int) *bookmarks.Bookmark {
	all := s.bookmarks.All()
	bmIdx := idx - 1 // index 0 is the header row
	if bmIdx < 0 || bmIdx >= len(all) {
		return nil
	}
	return &all[bmIdx]
}

func (s *Sidebar) setupDragDrop() {
	// Accept both internal (GTK_TREE_MODEL_ROW) and external (text/uri-list) drops
	internalTarget, _ := gtk.TargetEntryNew("GTK_TREE_MODEL_ROW", gtk.TARGET_SAME_APP, 0)
	internalTarget2, _ := gtk.TargetEntryNew("FILEX_INTERNAL", gtk.TARGET_SAME_APP, 1)
	externalTarget, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 2)
	targets := []gtk.TargetEntry{*internalTarget, *internalTarget2, *externalTarget}
	s.ListBox.DragDestSet(gtk.DEST_DEFAULT_MOTION|gtk.DEST_DEFAULT_DROP, targets, gdk.ACTION_COPY|gdk.ACTION_MOVE)

	s.ListBox.Connect("drag-drop", func(widget *gtk.ListBox, ctx *gdk.DragContext, x, y int, tm uint32) bool {
		if len(globalDragPaths) == 0 {
			return false
		}
		for _, p := range globalDragPaths {
			s.AddBookmark(p)
		}
		globalDragPaths = nil
		return true
	})
}

// AddBookmark adds a path as a user bookmark and re-renders.
func (s *Sidebar) AddBookmark(path string) {
	s.bookmarks.Add(bookmarks.Bookmark{
		Name:      filepath.Base(path),
		Path:      path,
		Icon:      "folder",
		UserAdded: true,
	})
	s.Render()
	s.App.Statusbar.ShowMessage("Bookmarked " + filepath.Base(path))
}
