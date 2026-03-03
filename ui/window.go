package ui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// BuildWindow assembles the main layout: toolbar, paned(sidebar|notebook), statusbar.
func BuildWindow(app *App) {
	mainBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	if err != nil {
		log.Fatal(err)
	}
	app.MainBox = mainBox

	// Toolbar
	app.Toolbar = NewToolbar(app)
	mainBox.PackStart(app.Toolbar.Box, false, false, 0)

	// Horizontal paned: sidebar | notebook
	paned, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		log.Fatal(err)
	}
	paned.SetPosition(180)

	// Sidebar
	app.Sidebar = NewSidebar(app)
	paned.Pack1(app.Sidebar.ScrollWin, false, false)

	// Notebook for tabs
	app.Notebook, err = gtk.NotebookNew()
	if err != nil {
		log.Fatal(err)
	}
	app.Notebook.SetScrollable(true)
	app.Notebook.SetShowBorder(false)
	app.Notebook.Connect("switch-page", func(nb *gtk.Notebook, page *gtk.Widget, pageNum uint) {
		tab, ok := tabRegistry[page.Native()]
		if ok {
			app.Toolbar.UpdateForTab(tab)
			app.Statusbar.Update(tab)
		}
	})
	paned.Pack2(app.Notebook, true, false)

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
