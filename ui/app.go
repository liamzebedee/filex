package ui

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
)

// Clipboard holds paths staged by Copy/Cut for a later Paste.
type Clipboard struct {
	Paths []string
	Cut   bool
}

// App holds the global application state: the window chrome and the
// clipboard. Per-tab state lives in each Tab.
type App struct {
	Window    *gtk.Window
	MainBox   *gtk.Box
	Notebook  *gtk.Notebook
	Sidebar   *Sidebar
	Statusbar *Statusbar

	Clipboard Clipboard
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
	a.Window.SetTitle("Files")
	a.Window.SetDefaultSize(900, 600)
	a.Window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Try loading icon from file (works without system install)
	iconPaths := []string{
		"/usr/share/icons/hicolor/scalable/apps/filex.svg",
	}
	if exe, err2 := os.Executable(); err2 == nil {
		iconPaths = append([]string{filepath.Join(filepath.Dir(exe), "assets", "filex.svg")}, iconPaths...)
	}
	if cwd, err2 := os.Getwd(); err2 == nil {
		iconPaths = append([]string{filepath.Join(cwd, "assets", "filex.svg")}, iconPaths...)
	}
	for _, p := range iconPaths {
		if err2 := gtk.WindowSetDefaultIconFromFile(p); err2 == nil {
			break
		}
	}

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
	return tabRegistry[widget.ToWidget().Native()]
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
