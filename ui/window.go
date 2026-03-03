package ui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"

	"filex/i18n"
)

// BuildWindow assembles the main layout:
// Sidebar full left, right side has notebook (each tab has its own toolbar + content).
func BuildWindow(app *App) {
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}
	app.MainBox = mainBox

	// Dark header bar for window decoration
	headerBar, err := gtk.HeaderBarNew()
	if err != nil {
		log.Fatal(err)
	}
	headerBar.SetShowCloseButton(true)
	headerBar.SetTitle(i18n.T("Files"))
	hsc, _ := headerBar.GetStyleContext()
	hsc.AddClass("dark-headerbar")
	app.Window.SetTitlebar(headerBar)

	// Horizontal paned: sidebar (full height) | right panel
	paned, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		log.Fatal(err)
	}
	paned.SetPosition(180)

	// Sidebar — full left column
	app.Sidebar = NewSidebar(app)
	paned.Pack1(app.Sidebar.ScrollWin, false, false)

	// Right panel: just the notebook (toolbar is per-tab, inside each tab's page)
	rightBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}

	app.Notebook, err = gtk.NotebookNew()
	if err != nil {
		log.Fatal(err)
	}
	app.Notebook.SetScrollable(true)
	app.Notebook.SetShowBorder(false)
	app.Notebook.Connect("switch-page", func(nb *gtk.Notebook, page *gtk.Widget, pageNum uint) {
		tab, ok := tabRegistry[page.Native()]
		if ok {
			app.Statusbar.Update(tab)
		}
	})
	rightBox.PackStart(app.Notebook, true, true, 0)

	paned.Pack2(rightBox, true, false)

	mainBox.PackStart(paned, true, true, 0)

	// Statusbar
	app.Statusbar = NewStatusbar()
	mainBox.PackStart(app.Statusbar.Box, false, false, 0)

	app.Window.Add(mainBox)

	// Open initial tab
	startPath := GetStartPath()
	NewTab(app, startPath)

	// Set up keyboard shortcuts
	setupKeyboardShortcuts(app)
}
