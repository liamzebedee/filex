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
	oldName := filepath.Base(filePath)
	newName, ok := runEntryDialog(tab.App.Window, entryDialogSpec{
		title:      "Rename",
		label:      "Enter new name:",
		initial:    oldName,
		confirm:    "Rename",
		selectStem: true,
	})
	if !ok || newName == "" || newName == oldName {
		return
	}
	newPath := filepath.Join(filepath.Dir(filePath), newName)
	if err := fileops.Rename(filePath, newPath); err != nil {
		showErrorDialog(tab.App.Window, "Rename Failed", err.Error())
		return
	}
	refreshAllTabs(tab.App)
}

// ShowNewFolderDialog displays a dialog to create a new folder.
func ShowNewFolderDialog(tab *Tab) {
	name, ok := runEntryDialog(tab.App.Window, entryDialogSpec{
		title:   "New Folder",
		label:   "Folder name:",
		initial: "New Folder",
		confirm: "Create",
	})
	if !ok || name == "" {
		return
	}
	if err := fileops.NewFolder(filepath.Join(tab.State.Path(), name)); err != nil {
		showErrorDialog(tab.App.Window, "Create Folder Failed", err.Error())
		return
	}
	refreshAllTabs(tab.App)
}

// entryDialogSpec describes a single-text-entry modal dialog.
type entryDialogSpec struct {
	title      string
	label      string
	initial    string
	confirm    string
	selectStem bool // pre-select the name without its extension
}

// runEntryDialog shows the dialog and returns the trimmed text, with
// ok=false when the user cancelled.
func runEntryDialog(parent *gtk.Window, spec entryDialogSpec) (string, bool) {
	dialog, _ := gtk.DialogNewWithButtons(
		spec.title,
		parent,
		gtk.DIALOG_MODAL|gtk.DIALOG_DESTROY_WITH_PARENT,
		[]interface{}{"Cancel", gtk.RESPONSE_CANCEL},
		[]interface{}{spec.confirm, gtk.RESPONSE_OK},
	)
	defer dialog.Destroy()
	dialog.SetDefaultSize(350, -1)

	// Mark the confirm button as suggested
	if okBtn, err := dialog.GetWidgetForResponse(gtk.RESPONSE_OK); err == nil {
		sc, _ := okBtn.ToWidget().GetStyleContext()
		sc.AddClass("suggested-action")
	}

	contentArea, _ := dialog.GetContentArea()
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)

	label, _ := gtk.LabelNew(spec.label)
	label.SetHAlign(gtk.ALIGN_START)
	box.PackStart(label, false, false, 0)

	entry, _ := gtk.EntryNew()
	entry.SetText(spec.initial)
	entry.SetActivatesDefault(true)
	box.PackStart(entry, false, false, 0)

	contentArea.PackStart(box, true, true, 0)
	dialog.SetDefaultResponse(gtk.RESPONSE_OK)
	dialog.ShowAll()

	entry.GrabFocus()
	selEnd := -1
	if spec.selectStem {
		selEnd = len(strings.TrimSuffix(spec.initial, filepath.Ext(spec.initial)))
	}
	entry.SelectRegion(0, selEnd)

	if dialog.Run() != gtk.RESPONSE_OK {
		return "", false
	}
	text, _ := entry.GetText()
	return strings.TrimSpace(text), true
}

// ShowDeleteConfirmDialog asks for confirmation before trashing files.
func ShowDeleteConfirmDialog(tab *Tab, paths []string) {
	var msg string
	if len(paths) == 1 {
		msg = fmt.Sprintf("Move %q to the Trash?", filepath.Base(paths[0]))
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
	trashBtn, _ := dialog.AddButton("Move to Trash", gtk.RESPONSE_OK)
	sc, _ := trashBtn.ToWidget().GetStyleContext()
	sc.AddClass("suggested-action")

	response := dialog.Run()
	dialog.Destroy()

	if response == gtk.RESPONSE_OK {
		trashPaths(tab.App, paths)
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
	defer dialog.Destroy()
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
		addPropRow("Type:", util.MimeFor(info.Name(), false))
		addPropRow("Size:", util.FormatSize(info.Size()))
	}
	addPropRow("Location:", filepath.Dir(filePath))
	addPropRow("Modified:", util.FormatDate(info.ModTime()))
	addPropRow("Permissions:", fmt.Sprintf("%o", info.Mode().Perm()))

	contentArea.PackStart(grid, true, true, 0)
	dialog.ShowAll()
	dialog.Run()
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
