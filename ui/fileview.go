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
	nameCol.SetTitle("Name")
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
	sizeCol.SetTitle("Size")
	sizeCol.SetResizable(true)
	sizeCol.SetMinWidth(80)
	sizeCol.SetSortColumnID(colSizeRaw)
	sizeRenderer, _ := gtk.CellRendererTextNew()
	sizeCol.PackStart(sizeRenderer, true)
	sizeCol.AddAttribute(sizeRenderer, "text", colSize)
	fv.TreeView.AppendColumn(sizeCol)

	// Column: Modified
	dateCol, _ := gtk.TreeViewColumnNew()
	dateCol.SetTitle("Modified")
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

	// Enable drag source for dragging files to sidebar bookmarks
	fv.setupDragSource(fv.TreeView.Widget)
}

func (fv *FileView) buildIconView() {
	var err error
	fv.IconView, err = gtk.IconViewNewWithModel(fv.Store)
	if err != nil {
		log.Fatal(err)
	}
	// Use column indices for icon view rendering
	// colIcon is a string column with icon names — we need pixbuf column
	// Since we store icon names (not pixbufs), we use the text column and set markup
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

	// Enable drag source for dragging files to sidebar bookmarks
	fv.setupDragSource(fv.IconView.Widget)
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

// Helper to get a string value from the ListStore
func getStringFromStore(store *gtk.ListStore, iter *gtk.TreeIter, col int) (string, error) {
	val, err := store.GetValue(iter, col)
	if err != nil {
		return "", err
	}
	return val.GetString()
}

// Helper to get a bool value from the ListStore
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

// setupDragSource configures a widget to provide URI data when dragged.
func (fv *FileView) setupDragSource(w gtk.Widget) {
	target, _ := gtk.TargetEntryNew("text/uri-list", gtk.TARGET_OTHER_APP, 0)
	w.DragSourceSet(gdk.BUTTON1_MASK, []gtk.TargetEntry{*target}, gdk.ACTION_COPY)

	w.Connect("drag-data-get", func(
		widget *gtk.Widget,
		ctx *gdk.DragContext,
		selData *gtk.SelectionData,
		info uint,
		time uint32,
	) {
		paths := fv.SelectedPaths()
		if len(paths) == 0 {
			return
		}
		var uris []string
		for _, p := range paths {
			uris = append(uris, "file://"+p)
		}
		selData.SetURIs(uris)
	})
}

// GetSelectedName returns the name of the first selected file
func (fv *FileView) GetSelectedName() string {
	p := fv.SelectedPath()
	if p == "" {
		return ""
	}
	return filepath.Base(p)
}
