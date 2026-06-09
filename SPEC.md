# Filex — Feature Spec

Parity checklist with a standard desktop file explorer (circa 2020).
One line per feature; all items are implemented and verified.

## Navigation
- [x] Back / Forward with per-tab, browser-style history
- [x] Up to parent directory
- [x] Breadcrumb path bar; click any segment to jump there
- [x] Click breadcrumb whitespace or Ctrl+L to type a path; `~` expands; bad paths rejected with a message
- [x] Re-clicking the current location (breadcrumb/bookmark) refreshes it
- [x] macOS/Finder-style keys: Ctrl+[ back, Ctrl+] forward, Ctrl+Up parent, Ctrl+Down open
- [x] Alt+Left / Alt+Right / Alt+Up and Backspace navigation
- [x] Enter opens the selected item; folders navigate, files open with their app
- [x] Symlinks to folders navigate like folders

## Tabs
- [x] Multiple tabs, each with independent path, history, view mode, sort, and filter
- [x] New tab (Ctrl+T), close tab (Ctrl+W or ✕), reorder by drag
- [x] Cycle tabs (Ctrl+Tab / Ctrl+PgUp / Ctrl+PgDn), last tab cannot be closed
- [x] Tab title follows the current folder name

## Viewing
- [x] List view: Name / Size / Modified columns with mime-type icons
- [x] Icon grid view: 48 px themed icons and real image thumbnails (cached, budgeted)
- [x] Per-tab list/grid toggle in the toolbar
- [x] Sort by name, size, or date by clicking column headers; click again to reverse (indicator shown)
- [x] Directories always group before files, in either direction
- [x] Hidden-file toggle (Ctrl+H), off by default
- [x] Refresh (F5)

## Search & filter
- [x] Live, case-insensitive name filter from the toolbar search box (per tab)
- [x] Escape clears the filter; navigating away clears it automatically
- [x] Statusbar item count reflects the filtered view

## Preview
- [x] Spacebar Quick Look on the selection; Space or Escape closes it
- [x] Images preview scaled-to-fit; text/code previews in a monospace view
- [x] Folders and other files get an icon + facts panel (type, size, modified)
- [x] Properties dialog (type, size, location, modified, permissions)

## File operations
- [x] Open files with the default application (double-click / Enter)
- [x] Copy / Cut / Paste (Ctrl+C/X/V and context menu); collisions become "name (copy)"
- [x] Pasting a folder into itself is safely refused
- [x] Drag-and-drop move: between rows, onto folders, across tabs, in both views
- [x] Rename (F2 or context menu) with the name stem pre-selected
- [x] Move to Trash (Delete) with confirmation (Enter confirms), freedesktop-spec trash
- [x] New folder (Ctrl+Shift+N or context menu)
- [x] Multi-select: Ctrl+A, rubber-band, Ctrl/Shift+click
- [x] Extract .zip archives in place (context menu)
- [x] Copy full path(s) to the clipboard
- [x] Open in Terminal

## Bookmarks
- [x] Default places (Home, Desktop, Documents, …, Trash, File System); missing dirs hidden
- [x] Add: Ctrl+D, context-menu "Add Bookmark", or drag a folder onto the sidebar
- [x] Remove: right-click a user bookmark → Remove
- [x] Persisted across sessions (~/.config/filex/bookmarks.txt); duplicates ignored

## Statusbar
- [x] Visible item count for the active tab
- [x] Free disk space for the current location
- [x] Transient feedback for copy/cut/paste/move/trash/bookmark/extract (auto-clears)

## Keyboard correctness
- [x] Shortcuts never steal keys from the path/search entries (typing Backspace, Ctrl+A, Enter… works)
- [x] Enter falls through to focused widgets (buttons, multi-select lists) when it cannot open
