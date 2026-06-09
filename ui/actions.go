package ui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
)

// This file holds every side-effecting action shared by shortcuts, menus
// and drag-and-drop. The pattern is uniform: read state on the main
// thread, run file IO in the background, then re-render on the main
// thread via glibIdle.

// itemCount returns e.g. "1 item" or "3 items".
func itemCount(n int) string {
	if n == 1 {
		return "1 item"
	}
	return fmt.Sprintf("%d items", n)
}

// itemCountMsg returns e.g. "3 items moved" or "1 item copied".
func itemCountMsg(n int, action string) string {
	return itemCount(n) + " " + action
}

// glibIdle schedules fn to run once on the GTK main thread.
func glibIdle(fn func()) {
	glib.IdleAdd(func() bool {
		fn()
		return false
	})
}

// refreshAllTabs re-reads every open tab. Used after operations that
// change the filesystem, since any tab may be showing an affected
// directory.
func refreshAllTabs(app *App) {
	for _, tab := range tabRegistry {
		tab.Refresh()
	}
}

// runFileOp runs a blocking file operation off the main thread, then
// refreshes all tabs and reports the outcome in the statusbar.
func runFileOp(app *App, successMsg string, op func() error) {
	go func() {
		err := op()
		glibIdle(func() {
			refreshAllTabs(app)
			if err != nil {
				app.Statusbar.ShowMessage("Error: " + err.Error())
			} else {
				app.Statusbar.ShowMessage(successMsg)
			}
		})
	}()
}

// setClipboard stages paths for a later paste and reports it.
func setClipboard(app *App, paths []string, cut bool) {
	if len(paths) == 0 {
		return
	}
	app.Clipboard = Clipboard{Paths: paths, Cut: cut}
	verb := "copied to clipboard"
	if cut {
		verb = "cut to clipboard"
	}
	app.Statusbar.ShowMessage(itemCountMsg(len(paths), verb))
}

// pasteClipboard pastes the staged paths into the tab's directory. The
// clipboard is read and (for cut) cleared here on the main thread; only
// the file IO runs in the background.
func pasteClipboard(app *App, tab *Tab) {
	cb := app.Clipboard
	if len(cb.Paths) == 0 {
		return
	}
	verb := "pasted"
	if cb.Cut {
		app.Clipboard = Clipboard{}
		verb = "moved"
	}
	dest := tab.State.Path()
	runFileOp(app, itemCountMsg(len(cb.Paths), verb), func() error {
		return fileops.PasteFiles(cb.Paths, dest, cb.Cut)
	})
}

// trashPaths moves paths to the freedesktop trash.
func trashPaths(app *App, paths []string) {
	runFileOp(app, itemCountMsg(len(paths), "moved to trash"), func() error {
		var firstErr error
		for _, p := range paths {
			if err := fileops.TrashFile(p); err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	})
}

// copyPathsToClipboard puts the paths, newline-separated, on the system
// text clipboard.
func copyPathsToClipboard(app *App, paths []string) {
	clip, err := gtk.ClipboardGet(gdk.GdkAtomIntern("CLIPBOARD", false))
	if err != nil {
		return
	}
	clip.SetText(strings.Join(paths, "\n"))
	app.Statusbar.ShowMessage("Path copied to clipboard")
}
