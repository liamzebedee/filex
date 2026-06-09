package ui

import (
	"log"
	"path/filepath"
	"slices"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"filex/core"
	"filex/fileops"
	"filex/util"
)

// ListStore columns. The store is a write-only render target: it holds
// display strings only, and reads go through entryAt into the rendered
// slice instead.
const (
	colIcon = iota // icon name (list view)
	colName
	colSize   // formatted
	colDate   // formatted
	colPixbuf // icon/thumbnail pixbuf (icon view)
)

// globalDragPaths holds paths being dragged across any FileView. It is
// global because the source and destination FileView are different
// instances (gotk3's SelectionData is broken, so data can't ride along
// with the drag itself).
var globalDragPaths []string

// FileView renders a tab's visible entries as both a list and an icon
// grid over one shared ListStore. rendered mirrors the store row-for-row
// and is the slice all reads (selection, activation, drops) resolve
// against.
type FileView struct {
	Tab       *Tab
	ScrollWin *gtk.ScrolledWindow
	Stack     *gtk.Stack
	TreeView  *gtk.TreeView
	IconView  *gtk.IconView
	Store     *gtk.ListStore

	rendered []core.FileEntry
	sortCols map[core.SortKey]*gtk.TreeViewColumn
}

func NewFileView(tab *Tab) *FileView {
	fv := &FileView{Tab: tab, sortCols: make(map[core.SortKey]*gtk.TreeViewColumn)}

	var err error
	fv.ScrollWin, err = gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	fv.ScrollWin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

	fv.Store, err = gtk.ListStoreNew(
		glib.TYPE_STRING,    // icon name
		glib.TYPE_STRING,    // name
		glib.TYPE_STRING,    // size (formatted)
		glib.TYPE_STRING,    // date (formatted)
		gdk.PixbufGetType(), // icon/thumbnail pixbuf
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
	fv.ScrollWin.Add(fv.Stack)

	return fv
}

// Render syncs the views to the visible entries. The store is repopulated
// only when the entries actually changed; view mode and sort indicators
// are cheap and sync unconditionally.
func (fv *FileView) Render(visible []core.FileEntry, s core.TabState) {
	if !slices.Equal(fv.rendered, visible) {
		fv.populateStore(visible)
		fv.rendered = visible
	}

	if s.ViewMode == core.IconMode {
		fv.Stack.SetVisibleChildName("icon")
	} else {
		fv.Stack.SetVisibleChildName("list")
	}

	for key, col := range fv.sortCols {
		col.SetSortIndicator(key == s.SortKey)
		if key == s.SortKey {
			if s.SortAsc {
				col.SetSortOrder(gtk.SORT_ASCENDING)
			} else {
				col.SetSortOrder(gtk.SORT_DESCENDING)
			}
		}
	}
}

// populateStore rewrites the store from the given entries. It is the only
// writer; nothing ever reads the store back.
func (fv *FileView) populateStore(entries []core.FileEntry) {
	fv.Store.Clear()
	thumbBudget := thumbsPerRender
	for _, e := range entries {
		sizeStr := ""
		if !e.IsDir {
			sizeStr = util.FormatSize(e.Size)
		}
		iter := fv.Store.Append()
		fv.Store.Set(iter,
			[]int{colIcon, colName, colSize, colDate},
			[]interface{}{
				util.IconForMime(util.MimeFor(e.Name, e.IsDir)),
				e.Name,
				sizeStr,
				util.FormatDate(time.Unix(e.ModTime, 0)),
			},
		)
		if pb := gridPixbufFor(e, &thumbBudget); pb != nil {
			fv.Store.SetValue(iter, colPixbuf, pb)
		}
	}
}

// entryAt maps a tree path back to its entry through the rendered slice.
func (fv *FileView) entryAt(tp *gtk.TreePath) (core.FileEntry, bool) {
	idx := tp.GetIndices()
	if len(idx) == 0 || idx[0] < 0 || idx[0] >= len(fv.rendered) {
		return core.FileEntry{}, false
	}
	return fv.rendered[idx[0]], true
}

// activateRow opens the entry at the given tree path: directories
// navigate, files open with their default application.
func (fv *FileView) activateRow(tp *gtk.TreePath) {
	e, ok := fv.entryAt(tp)
	if !ok {
		return
	}
	if e.IsDir {
		fv.Tab.NavigateTo(e.Path)
	} else {
		fileops.OpenFile(e.Path)
	}
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
	nameCol := fv.newSortColumn("Name", core.SortByName)
	nameCol.SetExpand(true)
	nameCol.SetMinWidth(200)

	iconRenderer, _ := gtk.CellRendererPixbufNew()
	nameCol.PackStart(iconRenderer, false)
	nameCol.AddAttribute(iconRenderer, "icon-name", colIcon)

	textRenderer, _ := gtk.CellRendererTextNew()
	nameCol.PackStart(textRenderer, true)
	nameCol.AddAttribute(textRenderer, "text", colName)

	fv.TreeView.AppendColumn(nameCol)

	// Column: Size
	sizeCol := fv.newSortColumn("Size", core.SortBySize)
	sizeCol.SetMinWidth(80)
	sizeRenderer, _ := gtk.CellRendererTextNew()
	sizeCol.PackStart(sizeRenderer, true)
	sizeCol.AddAttribute(sizeRenderer, "text", colSize)
	fv.TreeView.AppendColumn(sizeCol)

	// Column: Modified
	dateCol := fv.newSortColumn("Modified", core.SortByDate)
	dateCol.SetMinWidth(120)
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
			ShowContextMenu(fv.Tab, btnEvent)
			return true
		}
		return false
	})

	// DnD: Use EnableModelDragSource/Dest so GTK handles row-level drag and
	// drop. The target is a dummy — data travels via globalDragPaths.
	rowTarget, _ := gtk.TargetEntryNew("GTK_TREE_MODEL_ROW", gtk.TARGET_SAME_WIDGET, 0)
	crossTarget, _ := gtk.TargetEntryNew("GTK_TREE_MODEL_ROW", gtk.TARGET_SAME_APP, 1)
	targets := []gtk.TargetEntry{*rowTarget, *crossTarget}

	fv.TreeView.EnableModelDragSource(gdk.BUTTON1_MASK, targets, gdk.ACTION_MOVE)
	fv.TreeView.EnableModelDragDest(targets, gdk.ACTION_MOVE)

	fv.TreeView.Connect("drag-begin", func(tv *gtk.TreeView, ctx *gdk.DragContext) {
		globalDragPaths = fv.SelectedPaths()
	})

	fv.TreeView.Connect("drag-drop", func(tv *gtk.TreeView, ctx *gdk.DragContext, x, y int, tm uint32) bool {
		// Drop into the folder row under the cursor, else the current dir.
		destDir := fv.Tab.State.Path()
		if tp, pos, ok := fv.TreeView.GetDestRowAtPos(x, y); ok && tp != nil {
			onto := pos == gtk.TREE_VIEW_DROP_INTO_OR_BEFORE || pos == gtk.TREE_VIEW_DROP_INTO_OR_AFTER
			if e, ok := fv.entryAt(tp); ok && e.IsDir && onto {
				destDir = e.Path
			}
		}
		return fv.dropDragged(destDir)
	})

	fv.TreeView.Connect("drag-end", func(tv *gtk.TreeView, ctx *gdk.DragContext) {
		globalDragPaths = nil
	})
}

// newSortColumn creates a clickable column that dispatches the given sort
// key; the indicator is rendered from state in Render.
func (fv *FileView) newSortColumn(title string, key core.SortKey) *gtk.TreeViewColumn {
	col, _ := gtk.TreeViewColumnNew()
	col.SetTitle(title)
	col.SetResizable(true)
	col.SetClickable(true)
	col.Connect("clicked", func() {
		fv.Tab.SetSort(key)
	})
	fv.sortCols[key] = col
	return col
}

func (fv *FileView) buildIconView() {
	var err error
	fv.IconView, err = gtk.IconViewNewWithModel(fv.Store)
	if err != nil {
		log.Fatal(err)
	}
	fv.IconView.SetTextColumn(colName)
	fv.IconView.SetPixbufColumn(colPixbuf)
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

	// Handle double-click before DnD can intercept it, and right-click for
	// the context menu.
	fv.IconView.Connect("button-press-event", func(iv *gtk.IconView, event *gdk.Event) bool {
		btnEvent := gdk.EventButtonNewFromEvent(event)
		if btnEvent.Button() == gdk.BUTTON_SECONDARY {
			ShowContextMenu(fv.Tab, btnEvent)
			return true
		}
		if btnEvent.Button() == gdk.BUTTON_PRIMARY && btnEvent.Type() == gdk.EVENT_DOUBLE_BUTTON_PRESS {
			if tp := fv.IconView.GetPathAtPos(int(btnEvent.X()), int(btnEvent.Y())); tp != nil {
				fv.activateRow(tp)
				return true
			}
		}
		return false
	})

	// IconView DnD: widget-level drag source + dest on the same widget
	ivTarget, _ := gtk.TargetEntryNew("FILEX_INTERNAL", gtk.TARGET_SAME_APP, 0)
	ivTargets := []gtk.TargetEntry{*ivTarget}
	fv.IconView.ToWidget().DragSourceSet(gdk.BUTTON1_MASK, ivTargets, gdk.ACTION_MOVE)
	fv.IconView.ToWidget().DragDestSet(gtk.DEST_DEFAULT_MOTION|gtk.DEST_DEFAULT_DROP, ivTargets, gdk.ACTION_MOVE)

	fv.IconView.Connect("drag-begin", func(iv *gtk.IconView, ctx *gdk.DragContext) {
		globalDragPaths = fv.SelectedPaths()
	})

	fv.IconView.Connect("drag-drop", func(iv *gtk.IconView, ctx *gdk.DragContext, x, y int, tm uint32) bool {
		destDir := fv.Tab.State.Path()
		if tp := fv.IconView.GetPathAtPos(x, y); tp != nil {
			if e, ok := fv.entryAt(tp); ok && e.IsDir {
				destDir = e.Path
			}
		}
		return fv.dropDragged(destDir)
	})

	fv.IconView.Connect("drag-end", func(iv *gtk.IconView, ctx *gdk.DragContext) {
		globalDragPaths = nil
	})
}

// dropDragged moves the dragged paths into destDir, skipping paths already
// there, and clears the drag state. It returns true whenever drag data was
// present — even if nothing needed moving — so the drop never falls
// through to GTK's default model-row drop, which would reorder the store
// behind our back.
func (fv *FileView) dropDragged(destDir string) bool {
	if len(globalDragPaths) == 0 {
		return false
	}
	sources := make([]string, 0, len(globalDragPaths))
	for _, p := range globalDragPaths {
		if filepath.Dir(p) != destDir {
			sources = append(sources, p)
		}
	}
	globalDragPaths = nil
	if len(sources) > 0 {
		runFileOp(fv.Tab.App, itemCountMsg(len(sources), "moved"), func() error {
			return fileops.PasteFiles(sources, destDir, true)
		})
	}
	return true
}

// selectedItems returns the tree paths of the current selection in the
// active view.
func (fv *FileView) selectedItems() []*gtk.TreePath {
	var items []*gtk.TreePath
	collect := func(item interface{}) {
		if tp, ok := item.(*gtk.TreePath); ok {
			items = append(items, tp)
		}
	}
	if fv.Tab.State.ViewMode == core.ListMode {
		sel, err := fv.TreeView.GetSelection()
		if err != nil {
			return nil
		}
		sel.GetSelectedRows(fv.Store).Foreach(collect)
	} else {
		fv.IconView.GetSelectedItems().Foreach(collect)
	}
	return items
}

// SelectedPaths returns the filesystem paths of all selected entries.
func (fv *FileView) SelectedPaths() []string {
	items := fv.selectedItems()
	paths := make([]string, 0, len(items))
	for _, tp := range items {
		if e, ok := fv.entryAt(tp); ok {
			paths = append(paths, e.Path)
		}
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

// SelectAll selects every item in the active view.
func (fv *FileView) SelectAll() {
	if fv.Tab.State.ViewMode == core.ListMode {
		sel, err := fv.TreeView.GetSelection()
		if err != nil {
			return
		}
		sel.SelectAll()
	} else {
		fv.IconView.SelectAll()
	}
}

// OpenSelected activates the single selected item (opens folder or file)
// and reports whether it did; with no or multiple selection it declines so
// the key event can fall through to the focused widget's own handling.
func (fv *FileView) OpenSelected() bool {
	items := fv.selectedItems()
	if len(items) != 1 {
		return false
	}
	fv.activateRow(items[0])
	return true
}

// HasFocus reports whether one of the file view widgets has keyboard
// focus.
func (fv *FileView) HasFocus() bool {
	return fv.TreeView.ToWidget().HasFocus() || fv.IconView.ToWidget().HasFocus()
}
