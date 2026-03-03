package ui

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"filex/fileops"
)

// glib_idle_add schedules a function to run on the GTK main thread.
func glib_idle_add(fn func()) {
	glib.IdleAdd(func() bool {
		fn()
		return false // don't repeat
	})
}

// setupKeyboardShortcuts registers global keyboard shortcuts on the window.
func setupKeyboardShortcuts(app *App) {
	app.Window.Connect("key-press-event", func(win *gtk.Window, event *gdk.Event) bool {
		keyEvent := gdk.EventKeyNewFromEvent(event)
		key := keyEvent.KeyVal()
		state := gdk.ModifierType(keyEvent.State()) & gtk.AcceleratorGetDefaultModMask()

		ctrl := state&gdk.CONTROL_MASK != 0
		shift := state&gdk.SHIFT_MASK != 0

		tab := app.ActiveTab()
		if tab == nil {
			return false
		}

		switch {
		// Ctrl+T: New Tab
		case ctrl && !shift && key == gdk.KEY_t:
			NewTab(app, tab.Path)
			return true

		// Ctrl+W: Close Tab
		case ctrl && !shift && key == gdk.KEY_w:
			tab.Close()
			return true

		// Ctrl+Tab / Ctrl+Page_Down: Next Tab
		case ctrl && !shift && (key == gdk.KEY_Tab || key == gdk.KEY_Page_Down):
			page := app.Notebook.GetCurrentPage()
			if page < app.Notebook.GetNPages()-1 {
				app.Notebook.SetCurrentPage(page + 1)
			} else {
				app.Notebook.SetCurrentPage(0)
			}
			return true

		// Ctrl+Shift+Tab / Ctrl+Page_Up: Previous Tab
		case ctrl && shift && key == gdk.KEY_ISO_Left_Tab,
			ctrl && !shift && key == gdk.KEY_Page_Up:
			page := app.Notebook.GetCurrentPage()
			if page > 0 {
				app.Notebook.SetCurrentPage(page - 1)
			} else {
				app.Notebook.SetCurrentPage(app.Notebook.GetNPages() - 1)
			}
			return true

		// Ctrl+L: Focus path entry
		case ctrl && !shift && key == gdk.KEY_l:
			tab.Toolbar.ShowPathEntry()
			return true

		// Ctrl+H: Toggle hidden files
		case ctrl && !shift && key == gdk.KEY_h:
			tab.ToggleHidden()
			return true

		// F2: Rename
		case key == gdk.KEY_F2:
			if path := tab.FileView.SelectedPath(); path != "" {
				ShowRenameDialog(tab, path)
			}
			return true

		// Delete: Move to trash
		case key == gdk.KEY_Delete:
			paths := tab.FileView.SelectedPaths()
			if len(paths) > 0 {
				ShowDeleteConfirmDialog(tab, paths)
			}
			return true

		// Ctrl+C: Copy
		case ctrl && !shift && key == gdk.KEY_c:
			paths := tab.FileView.SelectedPaths()
			if len(paths) > 0 {
				app.ClipboardPaths = paths
				app.ClipboardCut = false
			}
			return true

		// Ctrl+X: Cut
		case ctrl && !shift && key == gdk.KEY_x:
			paths := tab.FileView.SelectedPaths()
			if len(paths) > 0 {
				app.ClipboardPaths = paths
				app.ClipboardCut = true
			}
			return true

		// Ctrl+V: Paste
		case ctrl && !shift && key == gdk.KEY_v:
			if len(app.ClipboardPaths) > 0 {
				go func() {
					PasteAndRefresh(app, tab)
				}()
			}
			return true

		// Alt+Left: Back
		case state&gdk.MOD1_MASK != 0 && key == gdk.KEY_Left:
			tab.GoBack()
			return true

		// Alt+Right: Forward
		case state&gdk.MOD1_MASK != 0 && key == gdk.KEY_Right:
			tab.GoForward()
			return true

		// Alt+Up: Parent directory
		case state&gdk.MOD1_MASK != 0 && key == gdk.KEY_Up:
			tab.GoUp()
			return true

		// Backspace: Go back
		case key == gdk.KEY_BackSpace:
			tab.GoBack()
			return true

		// Ctrl+Shift+N: New folder
		case ctrl && shift && key == gdk.KEY_N:
			ShowNewFolderDialog(tab)
			return true

		// F5: Refresh
		case key == gdk.KEY_F5:
			tab.FileView.Refresh()
			if app.Statusbar != nil {
				app.Statusbar.Update(tab)
			}
			return true
		}

		return false
	})
}

// PasteAndRefresh handles paste operation with UI refresh.
func PasteAndRefresh(app *App, tab *Tab) {
	fileops.PasteFiles(app.ClipboardPaths, tab.Path, app.ClipboardCut)
	if app.ClipboardCut {
		app.ClipboardPaths = nil
		app.ClipboardCut = false
	}
	gtkIdleRefresh(tab)
}
