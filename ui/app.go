package ui

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"

	"filex/i18n"
)

// App holds the global application state.
type App struct {
	Window    *gtk.Window
	MainBox   *gtk.Box
	Notebook  *gtk.Notebook
	Sidebar   *Sidebar
	Statusbar *Statusbar

	// Clipboard state for copy/cut/paste
	ClipboardPaths []string
	ClipboardCut   bool
}

func NewApp() *App {
	app := &App{}
	app.buildWindow()
	return app
}

func (a *App) buildWindow() {
	var err error
	a.Window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	a.Window.SetTitle(i18n.T("Files"))
	a.Window.SetDefaultSize(900, 600)
	a.Window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	a.Window.SetIconName("system-file-manager")

	BuildWindow(a)
}

// ActiveTab returns the currently visible tab.
func (a *App) ActiveTab() *Tab {
	page := a.Notebook.GetCurrentPage()
	if page < 0 {
		return nil
	}
	widget, err := a.Notebook.GetNthPage(page)
	if err != nil || widget == nil {
		return nil
	}
	w := widget.ToWidget()
	tab, ok := tabRegistry[w.Native()]
	if !ok {
		return nil
	}
	return tab
}

// GetStartPath returns the path to open on launch (argument or home).
func GetStartPath() string {
	if len(os.Args) > 1 {
		info, err := os.Stat(os.Args[1])
		if err == nil && info.IsDir() {
			return os.Args[1]
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return home
}
