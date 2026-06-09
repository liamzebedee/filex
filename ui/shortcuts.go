package ui

import (
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// shortcut maps one key chord to an action on the active tab — a
// declarative table instead of a switch, so every binding follows the
// same shape.
type shortcut struct {
	key    uint
	mods   gdk.ModifierType
	action func(app *App, tab *Tab)
}

const (
	modNone  = gdk.ModifierType(0)
	modCtrl  = gdk.CONTROL_MASK
	modShift = gdk.SHIFT_MASK
	modAlt   = gdk.MOD1_MASK
)

var windowShortcuts = []shortcut{
	// Tabs
	{gdk.KEY_t, modCtrl, func(app *App, tab *Tab) { NewTab(app, tab.State.Path()) }},
	{gdk.KEY_w, modCtrl, func(app *App, tab *Tab) { tab.Close() }},
	{gdk.KEY_Tab, modCtrl, func(app *App, tab *Tab) { cycleTab(app, +1) }},
	{gdk.KEY_Page_Down, modCtrl, func(app *App, tab *Tab) { cycleTab(app, +1) }},
	{gdk.KEY_ISO_Left_Tab, modCtrl | modShift, func(app *App, tab *Tab) { cycleTab(app, -1) }},
	{gdk.KEY_Page_Up, modCtrl, func(app *App, tab *Tab) { cycleTab(app, -1) }},

	// Navigation
	{gdk.KEY_Left, modAlt, func(app *App, tab *Tab) { tab.GoBack() }},
	{gdk.KEY_Right, modAlt, func(app *App, tab *Tab) { tab.GoForward() }},
	{gdk.KEY_Up, modAlt, func(app *App, tab *Tab) { tab.GoUp() }},
	{gdk.KEY_BackSpace, modNone, func(app *App, tab *Tab) { tab.GoBack() }},
	{gdk.KEY_bracketleft, modCtrl, func(app *App, tab *Tab) { tab.GoBack() }},
	{gdk.KEY_bracketright, modCtrl, func(app *App, tab *Tab) { tab.GoForward() }},
	{gdk.KEY_Up, modCtrl, func(app *App, tab *Tab) { tab.GoUp() }},
	{gdk.KEY_Down, modCtrl, func(app *App, tab *Tab) { tab.FileView.OpenSelected() }},
	{gdk.KEY_Return, modNone, func(app *App, tab *Tab) { tab.FileView.OpenSelected() }},
	{gdk.KEY_KP_Enter, modNone, func(app *App, tab *Tab) { tab.FileView.OpenSelected() }},
	{gdk.KEY_l, modCtrl, func(app *App, tab *Tab) { tab.Toolbar.ShowPathEntry() }},

	// View
	{gdk.KEY_h, modCtrl, func(app *App, tab *Tab) { tab.ToggleHidden() }},
	{gdk.KEY_F5, modNone, func(app *App, tab *Tab) { tab.Refresh() }},
	{gdk.KEY_a, modCtrl, func(app *App, tab *Tab) { tab.FileView.SelectAll() }},

	// File operations
	{gdk.KEY_c, modCtrl, func(app *App, tab *Tab) { setClipboard(app, tab.FileView.SelectedPaths(), false) }},
	{gdk.KEY_x, modCtrl, func(app *App, tab *Tab) { setClipboard(app, tab.FileView.SelectedPaths(), true) }},
	{gdk.KEY_v, modCtrl, func(app *App, tab *Tab) { pasteClipboard(app, tab) }},
	{gdk.KEY_N, modCtrl | modShift, func(app *App, tab *Tab) { ShowNewFolderDialog(tab) }},
	{gdk.KEY_F2, modNone, func(app *App, tab *Tab) {
		if path := tab.FileView.SelectedPath(); path != "" {
			ShowRenameDialog(tab, path)
		}
	}},
	{gdk.KEY_Delete, modNone, func(app *App, tab *Tab) {
		if paths := tab.FileView.SelectedPaths(); len(paths) > 0 {
			ShowDeleteConfirmDialog(tab, paths)
		}
	}},
}

// setupKeyboardShortcuts registers the shortcut table on the window.
func setupKeyboardShortcuts(app *App) {
	app.Window.Connect("key-press-event", func(win *gtk.Window, event *gdk.Event) bool {
		tab := app.ActiveTab()
		if tab == nil || focusInEditable(win) {
			return false
		}

		keyEvent := gdk.EventKeyNewFromEvent(event)
		key := keyEvent.KeyVal()
		mods := gdk.ModifierType(keyEvent.State()) & gtk.AcceleratorGetDefaultModMask()

		for _, s := range windowShortcuts {
			if s.key == key && s.mods == mods {
				s.action(app, tab)
				return true
			}
		}
		return false
	})
}

// focusInEditable reports whether keyboard focus is in a text-editing
// widget. Global shortcuts must not steal keys (Backspace, Enter, Ctrl+A,
// Ctrl+C…) from the path or search entries.
func focusInEditable(win *gtk.Window) bool {
	focused, err := win.GetFocus()
	if err != nil || focused == nil {
		return false
	}
	w := focused.ToWidget()
	if w == nil {
		return false
	}
	// With no explicit widget name set, GetName returns the type name,
	// e.g. "GtkEntry" or "GtkSearchEntry".
	name, err := w.GetName()
	return err == nil && strings.Contains(name, "Entry")
}

// cycleTab moves the active notebook page by delta, wrapping around.
func cycleTab(app *App, delta int) {
	n := app.Notebook.GetNPages()
	if n == 0 {
		return
	}
	page := (app.Notebook.GetCurrentPage() + delta + n) % n
	app.Notebook.SetCurrentPage(page)
}
