package ui

import (
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
)

// ShowContextMenu displays the right-click context menu. Selection state
// is read once, up front; every item dispatches through the same shared
// actions as the keyboard shortcuts.
func ShowContextMenu(tab *Tab, event *gdk.EventButton) {
	menu, _ := gtk.MenuNew()

	app := tab.App
	selectedPaths := tab.FileView.SelectedPaths()
	hasSelection := len(selectedPaths) > 0

	// Single selection, if any, and whether it is a directory
	selectedPath := ""
	isDir := false
	if len(selectedPaths) == 1 {
		selectedPath = selectedPaths[0]
		for _, e := range tab.Entries {
			if e.Path == selectedPath {
				isDir = e.IsDir
				break
			}
		}
	}

	// Open
	if hasSelection {
		openItem := menuItem("Open", "document-open")
		openItem.Connect("activate", func() {
			if isDir {
				tab.NavigateTo(selectedPath)
				return
			}
			for _, p := range selectedPaths {
				fileops.OpenFile(p)
			}
		})
		menu.Append(openItem)
	}

	// Open in Terminal
	if isDir || !hasSelection {
		termItem := menuItem("Open in Terminal", "utilities-terminal")
		termItem.Connect("activate", func() {
			dir := tab.State.Path()
			if isDir {
				dir = selectedPath
			}
			fileops.OpenTerminal(dir)
		})
		menu.Append(termItem)
	}

	// Add Bookmark (selected folder, or the current directory)
	if isDir || !hasSelection {
		bookmarkItem := menuItem("Add Bookmark", "bookmark-new")
		bookmarkItem.Connect("activate", func() {
			dir := tab.State.Path()
			if isDir {
				dir = selectedPath
			}
			app.Sidebar.AddBookmark(dir)
		})
		menu.Append(bookmarkItem)
	}

	addSeparator(menu)

	// New Folder
	newFolderItem := menuItem("New Folder…", "folder-new")
	newFolderItem.Connect("activate", func() {
		ShowNewFolderDialog(tab)
	})
	menu.Append(newFolderItem)

	addSeparator(menu)

	// Cut / Copy / Paste
	cutItem := menuItem("Cut", "edit-cut")
	cutItem.SetSensitive(hasSelection)
	cutItem.Connect("activate", func() {
		setClipboard(app, selectedPaths, true)
	})
	menu.Append(cutItem)

	copyItem := menuItem("Copy", "edit-copy")
	copyItem.SetSensitive(hasSelection)
	copyItem.Connect("activate", func() {
		setClipboard(app, selectedPaths, false)
	})
	menu.Append(copyItem)

	pasteItem := menuItem("Paste", "edit-paste")
	pasteItem.SetSensitive(len(app.Clipboard.Paths) > 0)
	pasteItem.Connect("activate", func() {
		pasteClipboard(app, tab)
	})
	menu.Append(pasteItem)

	addSeparator(menu)

	// Copy Path
	if hasSelection {
		copyPathItem := menuItem("Copy Path", "edit-copy")
		copyPathItem.Connect("activate", func() {
			copyPathsToClipboard(app, selectedPaths)
		})
		menu.Append(copyPathItem)
	}

	// Rename
	if selectedPath != "" {
		renameItem := menuItem("Rename…", "document-edit")
		renameItem.Connect("activate", func() {
			ShowRenameDialog(tab, selectedPath)
		})
		menu.Append(renameItem)
	}

	// Unzip (for .zip files)
	if strings.HasSuffix(strings.ToLower(selectedPath), ".zip") {
		unzipItem := menuItem("Extract Here", "package-x-generic")
		unzipItem.Connect("activate", func() {
			runFileOp(app, "Archive extracted", func() error {
				return fileops.Unzip(selectedPath, filepath.Dir(selectedPath))
			})
		})
		menu.Append(unzipItem)
	}

	addSeparator(menu)

	// Delete
	if hasSelection {
		deleteItem := menuItem("Move to Trash", "user-trash")
		deleteItem.Connect("activate", func() {
			ShowDeleteConfirmDialog(tab, selectedPaths)
		})
		menu.Append(deleteItem)
	}

	// Properties
	if selectedPath != "" {
		addSeparator(menu)
		propsItem := menuItem("Properties", "document-properties")
		propsItem.Connect("activate", func() {
			ShowPropertiesDialog(tab, selectedPath)
		})
		menu.Append(propsItem)
	}

	menu.ShowAll()
	menu.PopupAtPointer((*gdk.Event)(nil))
}

func menuItem(label, iconName string) *gtk.MenuItem {
	item, _ := gtk.MenuItemNew()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	icon, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_MENU)
	lbl, _ := gtk.LabelNew(label)
	lbl.SetHAlign(gtk.ALIGN_START)
	box.PackStart(icon, false, false, 0)
	box.PackStart(lbl, true, true, 0)
	item.Add(box)
	return item
}

func addSeparator(menu *gtk.Menu) {
	sep, _ := gtk.SeparatorMenuItemNew()
	menu.Append(sep)
}
