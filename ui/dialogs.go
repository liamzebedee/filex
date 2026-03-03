package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
	"filex/util"
)

// ShowRenameDialog displays a dialog to rename a file.
func ShowRenameDialog(tab *Tab, filePath string) {
	dialog, _ := gtk.DialogNewWithButtons(
		"Rename",
		tab.App.Window,
		gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
		[]interface{}{"Cancel", gtk.RESPONSE_CANCEL},
		[]interface{}{"Rename", gtk.RESPONSE_OK},
	)
	dialog.SetDefaultSize(350, -1)

	// Mark OK button as suggested
	okBtnW, err := dialog.GetWidgetForResponse(gtk.RESPONSE_OK)
	if err == nil {
		sc, _ := okBtnW.ToWidget().GetStyleContext()
		sc.AddClass("suggested-action")
	}

	contentArea, _ := dialog.GetContentArea()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)

	label, _ := gtk.LabelNew("Enter new name:")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackStart(label, false, false, 0)

	entry, _ := gtk.EntryNew()
	oldName := filepath.Base(filePath)
	entry.SetText(oldName)

	// Select name without extension
	ext := filepath.Ext(oldName)
	nameWithoutExt := strings.TrimSuffix(oldName, ext)
	entry.SetActivatesDefault(true)
	box.PackStart(entry, false, false, 0)

	contentArea.PackStart(box, true, true, 0)
	dialog.SetDefaultResponse(gtk.RESPONSE_OK)
	dialog.ShowAll()

	// Select the name part (without extension) after showing
	entry.GrabFocus()
	entry.SelectRegion(0, len(nameWithoutExt))

	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		newName, _ := entry.GetText()
		newName = strings.TrimSpace(newName)
		if newName != "" && newName != oldName {
			dir := filepath.Dir(filePath)
			newPath := filepath.Join(dir, newName)
			if err := fileops.Rename(filePath, newPath); err != nil {
				showErrorDialog(tab.App.Window, "Rename Failed", err.Error())
			} else {
				tab.FileView.Refresh()
			}
		}
	}
	dialog.Destroy()
}

// ShowNewFolderDialog displays a dialog to create a new folder.
func ShowNewFolderDialog(tab *Tab) {
	dialog, _ := gtk.DialogNewWithButtons(
		"New Folder",
		tab.App.Window,
		gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
		[]interface{}{"Cancel", gtk.RESPONSE_CANCEL},
		[]interface{}{"Create", gtk.RESPONSE_OK},
	)
	dialog.SetDefaultSize(350, -1)

	okBtnW, err := dialog.GetWidgetForResponse(gtk.RESPONSE_OK)
	if err == nil {
		sc, _ := okBtnW.ToWidget().GetStyleContext()
		sc.AddClass("suggested-action")
	}

	contentArea, _ := dialog.GetContentArea()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)

	label, _ := gtk.LabelNew("Folder name:")
	label.SetHAlign(gtk.ALIGN_START)
	box.PackStart(label, false, false, 0)

	entry, _ := gtk.EntryNew()
	entry.SetText("New Folder")
	entry.SetActivatesDefault(true)
	box.PackStart(entry, false, false, 0)

	contentArea.PackStart(box, true, true, 0)
	dialog.SetDefaultResponse(gtk.RESPONSE_OK)
	dialog.ShowAll()

	entry.GrabFocus()
	entry.SelectRegion(0, -1)

	response := dialog.Run()
	if response == gtk.RESPONSE_OK {
		name, _ := entry.GetText()
		name = strings.TrimSpace(name)
		if name != "" {
			newPath := filepath.Join(tab.Path, name)
			if err := fileops.NewFolder(newPath); err != nil {
				showErrorDialog(tab.App.Window, "Create Folder Failed", err.Error())
			} else {
				tab.FileView.Refresh()
			}
		}
	}
	dialog.Destroy()
}

// ShowDeleteConfirmDialog asks for confirmation before deleting files.
func ShowDeleteConfirmDialog(tab *Tab, paths []string) {
	var msg string
	if len(paths) == 1 {
		msg = fmt.Sprintf("Move \"%s\" to the Trash?", filepath.Base(paths[0]))
	} else {
		msg = fmt.Sprintf("Move %d items to the Trash?", len(paths))
	}

	dialog := gtk.MessageDialogNew(
		tab.App.Window,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_QUESTION,
		gtk.BUTTONS_NONE,
		"%s",
		msg,
	)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	trashBtnW, _ := dialog.AddButton("Move to Trash", gtk.RESPONSE_OK)
	sc, _ := trashBtnW.ToWidget().GetStyleContext()
	sc.AddClass("suggested-action")

	response := dialog.Run()
	dialog.Destroy()

	if response == gtk.RESPONSE_OK {
		n := len(paths)
		go func() {
			for _, p := range paths {
				fileops.TrashFile(p)
			}
			glib_idle_add(func() {
				tab.FileView.Refresh()
				if tab.App.Statusbar != nil {
					tab.App.Statusbar.Update(tab)
					tab.App.Statusbar.ShowMessage(itemCountMsg(n, "moved to trash"))
				}
			})
		}()
	}
}

// ShowPropertiesDialog shows basic file properties.
func ShowPropertiesDialog(tab *Tab, filePath string) {
	info, err := os.Stat(filePath)
	if err != nil {
		showErrorDialog(tab.App.Window, "Error", err.Error())
		return
	}

	dialog, _ := gtk.DialogNewWithButtons(
		"Properties",
		tab.App.Window,
		gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
		[]interface{}{"Close", gtk.RESPONSE_CLOSE},
	)
	dialog.SetDefaultSize(350, -1)

	contentArea, _ := dialog.GetContentArea()
	grid, _ := gtk.GridNew()
	grid.SetColumnSpacing(12)
	grid.SetRowSpacing(8)
	grid.SetMarginTop(12)
	grid.SetMarginBottom(12)
	grid.SetMarginStart(12)
	grid.SetMarginEnd(12)

	row := 0
	addPropRow := func(key, value string) {
		kLabel, _ := gtk.LabelNew(key)
		kLabel.SetHAlign(gtk.ALIGN_END)
		sc, _ := kLabel.GetStyleContext()
		_ = sc
		vLabel, _ := gtk.LabelNew(value)
		vLabel.SetHAlign(gtk.ALIGN_START)
		vLabel.SetSelectable(true)
		grid.Attach(kLabel, 0, row, 1, 1)
		grid.Attach(vLabel, 1, row, 1, 1)
		row++
	}

	addPropRow("Name:", info.Name())
	if info.IsDir() {
		addPropRow("Type:", "Folder")
	} else {
		mime := util.DetectMimeType(info)
		addPropRow("Type:", mime)
		addPropRow("Size:", util.FormatSize(info.Size()))
	}
	addPropRow("Location:", filepath.Dir(filePath))
	addPropRow("Modified:", util.FormatDate(info.ModTime()))
	addPropRow("Permissions:", fmt.Sprintf("%o", info.Mode().Perm()))

	contentArea.PackStart(grid, true, true, 0)
	dialog.ShowAll()
	dialog.Run()
	dialog.Destroy()
}

func showErrorDialog(parent *gtk.Window, title, message string) {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		"%s",
		message,
	)
	dialog.SetTitle(title)
	dialog.Run()
	dialog.Destroy()
}
