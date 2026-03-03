package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
	"filex/i18n"
)

// ShowContextMenu displays the right-click context menu.
func ShowContextMenu(tab *Tab, event *gdk.EventButton) {
	menu, _ := gtk.MenuNew()

	selectedPaths := tab.FileView.SelectedPaths()
	selectedPath := ""
	if len(selectedPaths) == 1 {
		selectedPath = selectedPaths[0]
	}
	hasSelection := len(selectedPaths) > 0

	// If single selection is a directory
	isDir := false
	if selectedPath != "" {
		info, err := os.Stat(selectedPath)
		if err == nil {
			isDir = info.IsDir()
		}
	}

	// Open
	if hasSelection {
		openItem := menuItem(i18n.T("Open"), "document-open")
		openItem.Connect("activate", func() {
			if isDir && selectedPath != "" {
				tab.NavigateAndPush(selectedPath)
			} else {
				for _, p := range selectedPaths {
					fileops.OpenFile(p)
				}
			}
		})
		menu.Append(openItem)
	}

	// Open in Terminal
	if isDir || !hasSelection {
		termItem := menuItem(i18n.T("Open in Terminal"), "utilities-terminal")
		termItem.Connect("activate", func() {
			dir := tab.Path
			if isDir && selectedPath != "" {
				dir = selectedPath
			}
			fileops.OpenTerminal(dir)
		})
		menu.Append(termItem)
	}

	addSeparator(menu)

	// New Folder
	newFolderItem := menuItem(i18n.T("New Folder…"), "folder-new")
	newFolderItem.Connect("activate", func() {
		ShowNewFolderDialog(tab)
	})
	menu.Append(newFolderItem)

	addSeparator(menu)

	// Cut
	cutItem := menuItem(i18n.T("Cut"), "edit-cut")
	cutItem.SetSensitive(hasSelection)
	cutItem.Connect("activate", func() {
		tab.App.ClipboardPaths = selectedPaths
		tab.App.ClipboardCut = true
	})
	menu.Append(cutItem)

	// Copy
	copyItem := menuItem(i18n.T("Copy"), "edit-copy")
	copyItem.SetSensitive(hasSelection)
	copyItem.Connect("activate", func() {
		tab.App.ClipboardPaths = selectedPaths
		tab.App.ClipboardCut = false
	})
	menu.Append(copyItem)

	// Paste
	pasteItem := menuItem(i18n.T("Paste"), "edit-paste")
	pasteItem.SetSensitive(len(tab.App.ClipboardPaths) > 0)
	pasteItem.Connect("activate", func() {
		dest := tab.Path
		go func() {
			fileops.PasteFiles(tab.App.ClipboardPaths, dest, tab.App.ClipboardCut)
			if tab.App.ClipboardCut {
				tab.App.ClipboardPaths = nil
				tab.App.ClipboardCut = false
			}
			gtkIdleRefresh(tab)
		}()
	})
	menu.Append(pasteItem)

	addSeparator(menu)

	// Copy Path
	if hasSelection {
		copyPathItem := menuItem(i18n.T("Copy Path"), "edit-copy")
		copyPathItem.Connect("activate", func() {
			fileops.CopyPathToClipboard(selectedPaths)
		})
		menu.Append(copyPathItem)
	}

	// Rename
	if len(selectedPaths) == 1 {
		renameItem := menuItem(i18n.T("Rename…"), "document-edit")
		renameItem.Connect("activate", func() {
			ShowRenameDialog(tab, selectedPath)
		})
		menu.Append(renameItem)
	}

	// Unzip (for .zip files)
	if selectedPath != "" && strings.HasSuffix(strings.ToLower(selectedPath), ".zip") {
		unzipItem := menuItem(i18n.T("Extract Here"), "package-x-generic")
		unzipItem.Connect("activate", func() {
			go func() {
				fileops.Unzip(selectedPath, filepath.Dir(selectedPath))
				gtkIdleRefresh(tab)
			}()
		})
		menu.Append(unzipItem)
	}

	addSeparator(menu)

	// Delete
	if hasSelection {
		deleteItem := menuItem(i18n.T("Move to Trash"), "user-trash")
		deleteItem.Connect("activate", func() {
			ShowDeleteConfirmDialog(tab, selectedPaths)
		})
		menu.Append(deleteItem)
	}

	// Properties (placeholder)
	if len(selectedPaths) == 1 {
		addSeparator(menu)
		propsItem := menuItem(i18n.T("Properties"), "document-properties")
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

func gtkIdleRefresh(tab *Tab) {
	glib_idle_add(func() {
		tab.FileView.Refresh()
		if tab.App.Statusbar != nil {
			tab.App.Statusbar.Update(tab)
		}
	})
}
