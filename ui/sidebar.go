package ui

import (
	"log"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/bookmarks"
	"filex/i18n"
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
		idx := row.GetIndex()
		bm := s.getBookmarkAt(idx)
		if bm == nil {
			return
		}
		if tab := app.ActiveTab(); tab != nil {
			tab.NavigateAndPush(bm.Path)
		}
	})

	// Enable drag-and-drop to add bookmarks
	s.setupDragDrop()

	s.ScrollWin.Add(s.ListBox)
	s.Populate()

	return s
}

func (s *Sidebar) Populate() {
	// Remove existing rows
	s.ListBox.GetChildren().Foreach(func(item interface{}) {
		if w, ok := item.(*gtk.Widget); ok {
			s.ListBox.Remove(w)
		}
	})

	// Add header
	header, _ := gtk.LabelNew(i18n.T("Places"))
	headerSc, _ := header.GetStyleContext()
	headerSc.AddClass("sidebar-header")
	header.SetHAlign(gtk.ALIGN_START)
	s.ListBox.Add(header)

	// Add bookmarks
	for _, bm := range s.bookmarks.All() {
		row := s.createBookmarkRow(bm)
		s.ListBox.Add(row)
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
		box.Connect("button-press-event", func(b *gtk.Box, event *gdk.Event) bool {
			btnEvent := gdk.EventButtonNewFromEvent(event)
			if btnEvent.Button() == gdk.BUTTON_SECONDARY {
				s.showRemoveMenu(bm.Path, btnEvent)
				return true
			}
			return false
		})
	}

	return box
}

func (s *Sidebar) showRemoveMenu(path string, event *gdk.EventButton) {
	menu, _ := gtk.MenuNew()

	removeItem, _ := gtk.MenuItemNewWithLabel(i18n.T("Remove Bookmark"))
	removeItem.Connect("activate", func() {
		s.bookmarks.Remove(path)
		s.Populate()
	})
	menu.Append(removeItem)
	menu.ShowAll()
	menu.PopupAtPointer((*gdk.Event)(nil))
}

func (s *Sidebar) getBookmarkAt(idx int) *bookmarks.Bookmark {
	// Index 0 is the header
	all := s.bookmarks.All()
	bmIdx := idx - 1 // account for header
	if bmIdx < 0 || bmIdx >= len(all) {
		return nil
	}
	bm := all[bmIdx]
	return &bm
}

func (s *Sidebar) setupDragDrop() {
	// Accept both internal (GTK_TREE_MODEL_ROW) and external (text/uri-list) drops
	internalTarget, _ := gtk.TargetEntryNew("GTK_TREE_MODEL_ROW", gtk.TARGET_SAME_APP, 0)
	internalTarget2, _ := gtk.TargetEntryNew("FILEX_INTERNAL", gtk.TARGET_SAME_APP, 1)
	externalTarget, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 2)
	targets := []gtk.TargetEntry{*internalTarget, *internalTarget2, *externalTarget}
	s.ListBox.DragDestSet(gtk.DEST_DEFAULT_MOTION|gtk.DEST_DEFAULT_DROP, targets, gdk.ACTION_COPY|gdk.ACTION_MOVE)

	s.ListBox.Connect("drag-drop", func(widget *gtk.ListBox, ctx *gdk.DragContext, x, y int, tm uint32) bool {
		log.Printf("[DnD] sidebar drag-drop: globalDragPaths=%d", len(globalDragPaths))

		if len(globalDragPaths) > 0 {
			log.Printf("[DnD]   sidebar adding bookmarks from globalDragPaths: %v", globalDragPaths)
			for _, p := range globalDragPaths {
				name := filepath.Base(p)
				s.bookmarks.Add(bookmarks.Bookmark{
					Name:      name,
					Path:      p,
					Icon:      "folder",
					UserAdded: true,
				})
			}
			s.Populate()
			return true
		}
		return false
	})
}

// AddBookmark adds a path as a sidebar bookmark.
func (s *Sidebar) AddBookmark(path string) {
	name := filepath.Base(path)
	s.bookmarks.Add(bookmarks.Bookmark{
		Name:      name,
		Path:      path,
		Icon:      "folder",
		UserAdded: true,
	})
	s.Populate()
}
