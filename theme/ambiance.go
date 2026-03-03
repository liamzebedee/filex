package theme

// AmbianceCSS provides the Ubuntu Ambiance ~2012 theme styling.
const AmbianceCSS = `
/* === Global === */
* {
    -gtk-icon-style: regular;
    font-size: 13px;
}

window {
    background-color: #f2f1f0;
}

/* === Dark HeaderBar (window title bar) === */
headerbar, .dark-headerbar {
    background: linear-gradient(to bottom, #4a4944, #3c3b37);
    border-bottom: 1px solid #2b2a27;
    color: #dfdbd2;
    min-height: 0;
    padding: 0px 4px;
    margin: 0;
}

headerbar .title, .dark-headerbar .title {
    color: #dfdbd2;
    font-weight: bold;
    font-size: 12px;
}

headerbar button, .dark-headerbar button {
    background: transparent;
    border: none;
    border-radius: 3px;
    color: #dfdbd2;
    box-shadow: none;
    min-height: 12px;
    min-width: 12px;
    padding: 2px 4px;
    margin: 0;
}

headerbar button:hover, .dark-headerbar button:hover {
    background: rgba(255,255,255,0.1);
}

headerbar button:active, .dark-headerbar button:active {
    background: rgba(0,0,0,0.15);
}

/* === Compact Navigation Toolbar === */
.toolbar-box {
    background: #edeceb;
    border-bottom: 1px solid #d0cfcd;
    padding: 3px 4px;
    min-height: 0;
}

/* Nav buttons: small flat arrows */
.nav-btn {
    background: transparent;
    border: 1px solid #c5c3c0;
    border-radius: 3px;
    color: #5c5b59;
    padding: 2px 6px;
    min-height: 14px;
    min-width: 14px;
    box-shadow: none;
}

.nav-btn:hover {
    background: #dddbd8;
    border-color: #b5b3b0;
}

.nav-btn:active {
    background: #d0cecc;
}

.nav-btn:disabled {
    opacity: 0.35;
}

/* === Breadcrumb / Path Bar === */
.breadcrumb-bar {
    background: #ffffff;
    border: 1px solid #c5c3c0;
    border-radius: 4px;
    padding: 0px 0px;
    min-height: 24px;
}

.breadcrumb-btn {
    background: transparent;
    border: none;
    border-radius: 0;
    color: #5c5b59;
    padding: 2px 6px;
    min-height: 18px;
    box-shadow: none;
    font-size: 13px;
    font-weight: normal;
}

.breadcrumb-btn:last-child {
    color: #2e2d2b;
    font-weight: bold;
}

.breadcrumb-btn:hover {
    background: rgba(0,0,0,0.06);
}

.breadcrumb-btn:active {
    background: rgba(0,0,0,0.1);
}

.breadcrumb-sep {
    color: #a1a09e;
    padding: 0 1px;
    min-width: 0;
    font-size: 13px;
    font-weight: normal;
}

.path-entry {
    background: #ffffff;
    border: 1px solid #c8622f;
    border-radius: 4px;
    color: #2e2d2b;
    padding: 2px 6px;
    min-height: 24px;
    caret-color: #2e2d2b;
}

.path-entry:focus {
    border-color: #c8622f;
    box-shadow: 0 0 0 1px rgba(200,98,47,0.25);
}

/* Small icon buttons in toolbar (search, view toggle, menu) */
.tool-icon-btn {
    background: transparent;
    border: 1px solid #c5c3c0;
    border-radius: 3px;
    color: #5c5b59;
    padding: 2px 4px;
    min-height: 14px;
    min-width: 14px;
    box-shadow: none;
}

.tool-icon-btn:hover {
    background: #dddbd8;
    border-color: #b5b3b0;
}

.tool-icon-btn:active,
.tool-icon-btn:checked {
    background: #d0cecc;
}

/* === Sidebar === */
.sidebar {
    background-color: #e8e6e3;
    border-right: 1px solid #c5c3c0;
}

.sidebar row {
    padding: 3px 8px;
    border-radius: 0;
    min-height: 24px;
}

.sidebar row:selected {
    background-color: #c8622f;
    color: white;
}

.sidebar row:hover:not(:selected) {
    background-color: #dddbd8;
}

.sidebar-header {
    font-weight: bold;
    font-size: 11px;
    color: #888580;
    padding: 8px 12px 4px 12px;
}

.sidebar label {
    font-size: 13px;
    color: #3c3b37;
}

.sidebar row:selected label {
    color: white;
}

/* === File View (TreeView / IconView) === */
.file-view {
    background-color: #ffffff;
    color: #3c3b37;
}

.file-view:selected {
    background-color: #c8622f;
    color: white;
}

treeview header button {
    background: linear-gradient(to bottom, #f7f6f5, #edeceb);
    border-bottom: 1px solid #c5c3c0;
    border-right: 1px solid #d8d7d5;
    color: #5c5b59;
    padding: 3px 8px;
    font-weight: normal;
    font-size: 12px;
}

treeview header button:hover {
    background: linear-gradient(to bottom, #ffffff, #f2f1f0);
}

treeview {
    -GtkTreeView-grid-line-width: 0;
    -GtkTreeView-horizontal-separator: 4;
}

treeview row:nth-child(even) {
    background-color: #fafaf9;
}

treeview row:selected {
    background-color: #c8622f;
    color: white;
}

iconview {
    background-color: #ffffff;
}

iconview:selected {
    background-color: #c8622f;
    color: white;
}

/* === Notebook Tabs === */
notebook header {
    background: #edeceb;
    border-bottom: 1px solid #c5c3c0;
    min-height: 0;
}

notebook header tabs {
    padding: 0;
    margin: 0;
}

notebook tab {
    background: #e2e1df;
    border: 1px solid #c5c3c0;
    border-bottom: none;
    border-radius: 4px 4px 0 0;
    padding: 2px 6px;
    margin: 2px 1px 0 1px;
    color: #5c5b59;
}

notebook tab:checked {
    background: #ffffff;
    border-bottom: 1px solid #ffffff;
    color: #2e2d2b;
}

notebook tab:hover:not(:checked) {
    background: #eaeae8;
}

notebook tab button {
    background: none;
    border: none;
    border-radius: 2px;
    padding: 0px;
    min-height: 12px;
    min-width: 12px;
    color: #888580;
}

notebook tab button:hover {
    background: rgba(0,0,0,0.1);
    color: #3c3b37;
}

/* === Statusbar === */
.statusbar {
    background: #edeceb;
    border-top: 1px solid #d0cfcd;
    padding: 1px 12px;
    font-size: 12px;
    color: #5c5b59;
    min-height: 18px;
}

/* === Context Menus === */
menu {
    background-color: #f7f6f5;
    border: 1px solid #a1a09e;
    border-radius: 0;
    padding: 4px 0;
    box-shadow: 2px 2px 6px rgba(0,0,0,0.2);
}

menu menuitem {
    padding: 3px 20px;
    color: #3c3b37;
    min-height: 20px;
}

menu menuitem:hover {
    background-color: #c8622f;
    color: white;
}

menu separator {
    background-color: #d0cfcd;
    min-height: 1px;
    margin: 3px 0;
}

/* === Scrollbars === */
scrollbar {
    background-color: #edeceb;
}

scrollbar slider {
    background-color: #c5c3c0;
    border-radius: 4px;
    min-width: 7px;
    min-height: 7px;
}

scrollbar slider:hover {
    background-color: #a1a09e;
}

scrollbar slider:active {
    background-color: #888580;
}

/* === Dialogs === */
dialog {
    background-color: #f2f1f0;
}

dialog .dialog-action-area button {
    background: linear-gradient(to bottom, #f7f6f5, #edeceb);
    border: 1px solid #c5c3c0;
    border-radius: 4px;
    padding: 4px 14px;
    color: #3c3b37;
}

dialog .dialog-action-area button:hover {
    background: linear-gradient(to bottom, #ffffff, #f2f1f0);
}

dialog .dialog-action-area button.suggested-action {
    background: linear-gradient(to bottom, #d4773e, #c8622f);
    border-color: #a04e22;
    color: white;
}

dialog .dialog-action-area button.suggested-action:hover {
    background: linear-gradient(to bottom, #db8a55, #d4773e);
}

/* === Search Entry === */
.search-entry {
    background: #ffffff;
    border: 1px solid #c5c3c0;
    border-radius: 4px;
    color: #2e2d2b;
    padding: 2px 6px;
    min-height: 24px;
    min-width: 140px;
}

.search-entry:focus {
    border-color: #c8622f;
}

/* === Paned separator === */
paned separator {
    background-color: #c5c3c0;
    min-width: 1px;
    min-height: 1px;
}

/* === View toggle === */
.view-toggle {
    border: 1px solid #c5c3c0;
    border-radius: 3px;
    padding: 0;
    background: transparent;
}

.view-toggle button {
    background: transparent;
    border: none;
    border-radius: 2px;
    color: #5c5b59;
    padding: 2px 4px;
    min-height: 14px;
    min-width: 18px;
    box-shadow: none;
}

.view-toggle button:checked {
    background: #d0cecc;
}

.view-toggle button:hover:not(:checked) {
    background: rgba(0,0,0,0.06);
}

/* === Drop highlight === */
.drop-target {
    border: 2px dashed #c8622f;
    border-radius: 4px;
}
`
