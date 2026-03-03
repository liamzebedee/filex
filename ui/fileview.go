package ui

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
	"filex/i18n"
	"filex/util"
)

// ListStore columns
const (
	colIcon    = 0 // string (icon name)
	colName    = 1 // string
	colSize    = 2 // string (formatted)
	colDate    = 3 // string (formatted)
	colPath    = 4 // string (full path)
	colIsDir   = 5 // bool
	colSizeRaw = 6 // int64 (for sorting)
	colDateRaw = 7 // int64 (for sorting)
)

// globalDragPaths holds paths being dragged across any FileView.
// This is global because the source and destination FileView are different instances.
var globalDragPaths []string

// FileView manages both list and icon views with a shared ListStore.
type FileView struct {
	Tab       *Tab
	ScrollWin *gtk.ScrolledWindow
	Stack     *gtk.Stack
	TreeView  *gtk.TreeView
	IconView  *gtk.IconView
	Store     *gtk.ListStore
	Entries   []fileops.FileEntry
}

func NewFileView(tab *Tab) *FileView {
	fv := &FileView{Tab: tab}

	var err error
	fv.ScrollWin, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fv.ScrollWin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

	// Shared ListStore: icon, name, size, date, path, isDir, sizeRaw, dateRaw
	fv.Store, err = gtk.ListStoreNew(
		glib.TYPE_STRING,  // icon name
		glib.TYPE_STRING,  // name
		glib.TYPE_STRING,  // size (formatted)
		glib.TYPE_STRING,  // date (formatted)
		glib.TYPE_STRING,  // full path
		glib.TYPE_BOOLEAN, // isDir
		glib.TYPE_INT64,   // size raw
		glib.TYPE_INT64,   // date raw
	)
	if err != nil {
		log.Fatal(err)
	}

	// Stack to switch between list and icon views
	fv.Stack, err = gtk.StackNew()
	if err != nil {
		log.Fatal(err)
	}
	fv.Stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_CROSSFADE)
	fv.Stack.SetTransitionDuration(150)

	fv.buildTreeView()
	fv.buildIconView()

	fv.Stack.AddNamed(fv.TreeView, "list")
	fv.Stack.AddNamed(fv.IconView, "icon")

	if tab.ViewMode == IconMode {
		fv.Stack.SetVisibleChildName("icon")
	} else {
		fv.Stack.SetVisibleChildName("list")
	}

	fv.ScrollWin.Add(fv.Stack)

	return fv
}

func (fv *FileView) buildTreeView() {
	var err error
	fv.TreeView, err = gtk.TreeViewNewWithModel(fv.Store)
	if err != nil {
		log.Fatal(err)
	}
	fv.TreeView.SetHeadersVisible(true)
	fv.TreeView.SetEnableSearch(false)
	fv.TreeView.SetRubberBanding(true)

	sc, _ := fv.TreeView.GetStyleContext()
	sc.AddClass("file-view")

	sel, _ := fv.TreeView.GetSelection()
	sel.SetMode(gtk.SELECTION_MULTIPLE)

	// Column: Icon + Name
	nameCol, _ := gtk.TreeViewColumnNew()
	nameCol.SetTitle(i18n.T("Name"))
	nameCol.SetExpand(true)
	nameCol.SetResizable(true)
	nameCol.SetMinWidth(200)
	nameCol.SetSortColumnID(colName)

	iconRenderer, _ := gtk.CellRendererPixbufNew()
	nameCol.PackStart(iconRenderer, false)
	nameCol.AddAttribute(iconRenderer, "icon-name", colIcon)

	textRenderer, _ := gtk.CellRendererTextNew()
	nameCol.PackStart(textRenderer, true)
	nameCol.AddAttribute(textRenderer, "text", colName)

	fv.TreeView.AppendColumn(nameCol)

	// Column: Size
	sizeCol, _ := gtk.TreeViewColumnNew()
	sizeCol.SetTitle(i18n.T("Size"))
	sizeCol.SetResizable(true)
	sizeCol.SetMinWidth(80)
	sizeCol.SetSortColumnID(colSizeRaw)
	sizeRenderer, _ := gtk.CellRendererTextNew()
	sizeCol.PackStart(sizeRenderer, true)
	sizeCol.AddAttribute(sizeRenderer, "text", colSize)
	fv.TreeView.AppendColumn(sizeCol)

	// Column: Modified
	dateCol, _ := gtk.TreeViewColumnNew()
	dateCol.SetTitle(i18n.T("Modified"))
	dateCol.SetResizable(true)
	dateCol.SetMinWidth(120)
	dateCol.SetSortColumnID(colDateRaw)
	dateRenderer, _ := gtk.CellRendererTextNew()
	dateCol.PackStart(dateRenderer, true)
	dateCol.AddAttribute(dateRenderer, "text", colDate)
	fv.TreeView.AppendColumn(dateCol)

	// Double-click to open
	fv.TreeView.Connect("row-activated", func(tv *gtk.TreeView, path *gtk.TreePath, col *gtk.TreeViewColumn) {
		fv.activateRow(path)
	})

	// Right-click context menu
	fv.TreeView.Connect("button-press-event", func(tv *gtk.TreeView, event *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(event)
		if btnEvent.Button() == gdk.BUTTON_SECONDARY {
			fv.showContextMenu(btnEvent)
			return true
		}
		return false
	})

	// DnD: Use EnableModelDragSource/Dest so GTK handles row-level drag and drop.
	// We use a dummy target and never touch SelectionData (it's broken in gotk3).
	// Instead we use globalDragPaths to pass data between drag-begin and drag-drop.
	rowTarget, _ := gtk.TargetEntryNew("GTK_TREE_MODEL_ROW", gtk.TARGET_SAME_WIDGET, 0)
	crossTarget, _ := gtk.TargetEntryNew("GTK_TREE_MODEL_ROW", gtk.TARGET_SAME_APP, 1)
	targets := []gtk.TargetEntry{*rowTarget, *crossTarget}

	fv.TreeView.EnableModelDragSource(gdk.BUTTON1_MASK, targets, gdk.ACTION_MOVE)
	fv.TreeView.EnableModelDragDest(targets, gdk.ACTION_MOVE)

	fv.TreeView.Connect("drag-begin", func(tv *gtk.TreeView, ctx *gdk.DragContext) {
		globalDragPaths = fv.SelectedPaths()
		log.Printf("[DnD] TreeView drag-begin: %d paths: %v", len(globalDragPaths), globalDragPaths)
	})

	// drag-drop: the actual drop event. We handle it entirely via globalDragPaths.
	fv.TreeView.Connect("drag-drop", func(tv *gtk.TreeView, ctx *gdk.DragContext, x, y int, tm uint32) bool {
		log.Printf("[DnD] TreeView drag-drop at (%d,%d), globalDragPaths=%d", x, y, len(globalDragPaths))
		if len(globalDragPaths) == 0 {
			return false
		}

		// Determine drop target: which row are we over?
		destDir := fv.Tab.Path // default: drop into current directory
		treePath, pos, ok := fv.TreeView.GetDestRowAtPos(x, y)
		if ok && treePath != nil {
			log.Printf("[DnD]   drop on row %s, pos=%v", treePath.String(), pos)
			iter, err := fv.Store.GetIter(treePath)
			if err == nil {
				isDir, _ := getBoolFromStore(fv.Store, iter, colIsDir)
				if isDir && (pos == gtk.TREE_VIEW_DROP_INTO_OR_BEFORE || pos == gtk.TREE_VIEW_DROP_INTO_OR_AFTER) {
					rowPath, _ := getStringFromStore(fv.Store, iter, colPath)
					destDir = rowPath
					log.Printf("[DnD]   dropping INTO folder: %s", destDir)
				}
			}
		} else {
			log.Printf("[DnD]   drop on empty area, using current dir: %s", destDir)
		}

		sources := make([]string, 0, len(globalDragPaths))
		for _, p := range globalDragPaths {
			if filepath.Dir(p) != destDir {
				sources = append(sources, p)
			}
		}
		globalDragPaths = nil

		log.Printf("[DnD]   moving %d files -> %s", len(sources), destDir)
		if len(sources) > 0 {
			go func() {
				err := fileops.PasteFiles(sources, destDir, true)
				log.Printf("[DnD]   PasteFiles result: err=%v", err)
				glib_idle_add(func() {
					fv.Refresh()
					fv.refreshAllTabs()
				})
			}()
		}

		return true
	})

	fv.TreeView.Connect("drag-end", func(tv *gtk.TreeView, ctx *gdk.DragContext) {
		log.Printf("[DnD] TreeView drag-end")
		globalDragPaths = nil
	})
}

func (fv *FileView) buildIconView() {
	var err error
	fv.IconView, err = gtk.IconViewNewWithModel(fv.Store)
	if err != nil {
		log.Fatal(err)
	}
	fv.IconView.SetTextColumn(colName)
	fv.IconView.SetSelectionMode(gtk.SELECTION_MULTIPLE)
	fv.IconView.SetItemWidth(80)
	fv.IconView.SetColumnSpacing(8)
	fv.IconView.SetRowSpacing(4)

	sc, _ := fv.IconView.GetStyleContext()
	sc.AddClass("file-view")

	// Double-click to open
	fv.IconView.Connect("item-activated", func(iv *gtk.IconView, path *gtk.TreePath) {
		fv.activateRow(path)
	})

	// Right-click
	fv.IconView.Connect("button-press-event", func(iv *gtk.IconView, event *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(event)
		if btnEvent.Button() == gdk.BUTTON_SECONDARY {
			fv.showContextMenu(btnEvent)
			return true
		}
		return false
	})

	// IconView DnD: use widget-level drag source + dest on the same widget
	ivTarget, _ := gtk.TargetEntryNew("FILEX_INTERNAL", gtk.TARGET_SAME_APP, 0)
	ivTargets := []gtk.TargetEntry{*ivTarget}
	fv.IconView.ToWidget().DragSourceSet(gdk.BUTTON1_MASK, ivTargets, gdk.ACTION_MOVE)
	fv.IconView.ToWidget().DragDestSet(gtk.DEST_DEFAULT_MOTION|gtk.DEST_DEFAULT_DROP, ivTargets, gdk.ACTION_MOVE)

	fv.IconView.Connect("drag-begin", func(iv *gtk.IconView, ctx *gdk.DragContext) {
		globalDragPaths = fv.SelectedPaths()
		log.Printf("[DnD] IconView drag-begin: %d paths: %v", len(globalDragPaths), globalDragPaths)
	})

	fv.IconView.Connect("drag-drop", func(iv *gtk.IconView, ctx *gdk.DragContext, x, y int, tm uint32) bool {
		log.Printf("[DnD] IconView drag-drop at (%d,%d), globalDragPaths=%d", x, y, len(globalDragPaths))
		if len(globalDragPaths) == 0 {
			return false
		}
		destDir := fv.Tab.Path
		// Check if dropped on a folder icon
		treePath := fv.IconView.GetPathAtPos(x, y)
		if treePath != nil {
			iter, err := fv.Store.GetIter(treePath)
			if err == nil {
				isDir, _ := getBoolFromStore(fv.Store, iter, colIsDir)
				if isDir {
					rowPath, _ := getStringFromStore(fv.Store, iter, colPath)
					destDir = rowPath
					log.Printf("[DnD]   IconView drop on folder: %s", destDir)
				}
			}
		}

		sources := make([]string, 0, len(globalDragPaths))
		for _, p := range globalDragPaths {
			if filepath.Dir(p) != destDir {
				sources = append(sources, p)
			}
		}
		globalDragPaths = nil

		log.Printf("[DnD]   IconView moving %d files -> %s", len(sources), destDir)
		if len(sources) > 0 {
			go func() {
				err := fileops.PasteFiles(sources, destDir, true)
				log.Printf("[DnD]   PasteFiles result: err=%v", err)
				glib_idle_add(func() {
					fv.Refresh()
					fv.refreshAllTabs()
				})
			}()
		}
		return true
	})

	fv.IconView.Connect("drag-end", func(iv *gtk.IconView, ctx *gdk.DragContext) {
		log.Printf("[DnD] IconView drag-end")
		globalDragPaths = nil
	})
}

// refreshAllTabs refreshes all open tabs (useful after file move operations).
func (fv *FileView) refreshAllTabs() {
	app := fv.Tab.App
	for _, tab := range tabRegistry {
		if tab != fv.Tab {
			tab.FileView.Refresh()
		}
	}
	if app.Statusbar != nil {
		app.Statusbar.Update(fv.Tab)
	}
}

func (fv *FileView) activateRow(path *gtk.TreePath) {
	iter, err := fv.Store.GetIter(path)
	if err != nil {
		return
	}
	filePath, _ := getStringFromStore(fv.Store, iter, colPath)
	isDir, _ := getBoolFromStore(fv.Store, iter, colIsDir)

	if isDir {
		fv.Tab.NavigateAndPush(filePath)
	} else {
		go exec.Command("xdg-open", filePath).Start()
	}
}

func (fv *FileView) showContextMenu(event *gdk.EventButton) {
	ShowContextMenu(fv.Tab, event)
}

// Refresh reloads directory contents into the store.
func (fv *FileView) Refresh() {
	entries, err := fileops.ListDirectory(fv.Tab.Path, fv.Tab.ShowHidden)
	if err != nil {
		log.Printf("Error reading directory %s: %v", fv.Tab.Path, err)
		return
	}

	fileops.SortEntries(entries, fv.Tab.SortColumn, fv.Tab.SortAsc)
	fv.Entries = entries
	fv.populateStore(entries)
}

// SwitchView toggles between list and icon view modes.
func (fv *FileView) SwitchView(mode ViewMode) {
	if mode == IconMode {
		fv.Stack.SetVisibleChildName("icon")
	} else {
		fv.Stack.SetVisibleChildName("list")
	}
}

// SelectedPaths returns paths of all selected files.
func (fv *FileView) SelectedPaths() []string {
	var paths []string

	if fv.Tab.ViewMode == ListMode {
		sel, err := fv.TreeView.GetSelection()
		if err != nil {
			return nil
		}
		rows := sel.GetSelectedRows(fv.Store)
		rows.Foreach(func(item interface{}) {
			path := item.(*gtk.TreePath)
			iter, err := fv.Store.GetIter(path)
			if err != nil {
				return
			}
			p, _ := getStringFromStore(fv.Store, iter, colPath)
			paths = append(paths, p)
		})
	} else {
		selectedItems := fv.IconView.GetSelectedItems()
		selectedItems.Foreach(func(item interface{}) {
			path := item.(*gtk.TreePath)
			iter, err := fv.Store.GetIter(path)
			if err != nil {
				return
			}
			p, _ := getStringFromStore(fv.Store, iter, colPath)
			paths = append(paths, p)
		})
	}

	return paths
}

// SelectedPath returns the single selected path, or "" if none/multiple.
func (fv *FileView) SelectedPath() string {
	paths := fv.SelectedPaths()
	if len(paths) == 1 {
		return paths[0]
	}
	return ""
}

func getStringFromStore(store *gtk.ListStore, iter *gtk.TreeIter, col int) (string, error) {
	val, err := store.GetValue(iter, col)
	if err != nil {
		return "", err
	}
	return val.GetString()
}

func getBoolFromStore(store *gtk.ListStore, iter *gtk.TreeIter, col int) (bool, error) {
	val, err := store.GetValue(iter, col)
	if err != nil {
		return false, err
	}
	v, err := val.GoValue()
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, nil
	}
	return b, nil
}

// fakeFileInfo adapts FileEntry to os.FileInfo for mime detection.
type fakeFileInfo struct {
	entry fileops.FileEntry
}

func (f fakeFileInfo) Name() string       { return f.entry.Name }
func (f fakeFileInfo) IsDir() bool        { return f.entry.IsDir }
func (f fakeFileInfo) Size() int64        { return f.entry.Size }
func (f fakeFileInfo) Mode() os.FileMode  { return f.entry.Mode }
func (f fakeFileInfo) ModTime() time.Time { return time.Unix(f.entry.ModTime, 0) }
func (f fakeFileInfo) Sys() interface{}   { return nil }

// populateStore fills the ListStore with the given entries.
func (fv *FileView) populateStore(entries []fileops.FileEntry) {
	fv.Store.Clear()
	for _, e := range entries {
		iter := fv.Store.Append()
		mime := util.DetectMimeType(fakeFileInfo{e})
		icon := util.IconForMime(mime)
		sizeStr := ""
		if !e.IsDir {
			sizeStr = util.FormatSize(e.Size)
		}
		dateStr := util.FormatDate(time.Unix(e.ModTime, 0))
		fv.Store.Set(iter,
			[]int{colIcon, colName, colSize, colDate, colPath, colIsDir, colSizeRaw, colDateRaw},
			[]interface{}{icon, e.Name, sizeStr, dateStr, e.Path, e.IsDir, e.Size, e.ModTime},
		)
	}
}

// GetSelectedName returns the name of the first selected file
func (fv *FileView) GetSelectedName() string {
	p := fv.SelectedPath()
	if p == "" {
		return ""
	}
	return filepath.Base(p)
}
