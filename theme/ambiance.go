package theme

// AmbianceCSS provides the Ubuntu Ambiance ~2012 theme styling.
const AmbianceCSS = `
/* === Global === */
* {
    -gtk-icon-style: regular;
}

window {
    background-color: #f2f1f0;
}

/* === Header / Toolbar === */
.toolbar-box {
    background: linear-gradient(to bottom, #4a4944, #3c3b37);
    border-bottom: 1px solid #2b2a27;
    padding: 4px 6px;
    min-height: 38px;
}

.toolbar-box button {
    background: linear-gradient(to bottom, #565550, #484743);
    border: 1px solid #2b2a27;
    border-radius: 3px;
    color: #dfdbd2;
    padding: 2px 8px;
    min-height: 24px;
    min-width: 24px;
    box-shadow: inset 0 1px 0 rgba(255,255,255,0.05);
}

.toolbar-box button:hover {
    background: linear-gradient(to bottom, #636260, #555450);
    border-color: #222;
}

.toolbar-box button:active,
.toolbar-box button:checked {
    background: linear-gradient(to bottom, #3a3935, #444340);
    box-shadow: inset 0 2px 3px rgba(0,0,0,0.2);
}

.toolbar-box button:disabled {
    opacity: 0.4;
}

/* === Breadcrumb / Path Bar === */
.breadcrumb-bar {
    background: #4e4d49;
    border: 1px solid #2b2a27;
    border-radius: 3px;
    padding: 0px 2px;
}

.breadcrumb-btn {
    background: none;
    border: none;
    border-radius: 2px;
    color: #dfdbd2;
    padding: 2px 8px;
    min-height: 22px;
    box-shadow: none;
}

.breadcrumb-btn:hover {
    background: rgba(255,255,255,0.1);
}

.breadcrumb-btn:active {
    background: rgba(0,0,0,0.15);
}

.breadcrumb-sep {
    color: #7a7870;
    padding: 0 0;
    min-width: 8px;
}

.path-entry {
    background: #4e4d49;
    border: 1px solid #f07746;
    border-radius: 3px;
    color: #dfdbd2;
    padding: 2px 6px;
    min-height: 26px;
    caret-color: #dfdbd2;
}

.path-entry:focus {
    border-color: #f07746;
    box-shadow: 0 0 0 1px rgba(240,119,70,0.3);
}

/* === Sidebar === */
.sidebar {
    background-color: #e8e6e3;
    border-right: 1px solid #c5c3c0;
}

.sidebar row {
    padding: 4px 8px;
    border-radius: 0;
    min-height: 28px;
}

.sidebar row:selected {
    background-color: #f07746;
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
    background-color: #f07746;
    color: white;
}

treeview header button {
    background: linear-gradient(to bottom, #f7f6f5, #e8e6e3);
    border-bottom: 1px solid #c5c3c0;
    border-right: 1px solid #c5c3c0;
    color: #3c3b37;
    padding: 4px 8px;
    font-weight: normal;
    font-size: 12px;
}

treeview header button:hover {
    background: linear-gradient(to bottom, #ffffff, #eeede9);
}

treeview {
    -GtkTreeView-grid-line-width: 0;
    -GtkTreeView-horizontal-separator: 4;
}

treeview row:nth-child(even) {
    background-color: #fafaf9;
}

treeview row:selected {
    background-color: #f07746;
    color: white;
}

iconview {
    background-color: #ffffff;
}

iconview:selected {
    background-color: #f07746;
    color: white;
}

/* === Notebook Tabs === */
notebook header {
    background: linear-gradient(to bottom, #e8e6e3, #dddbd8);
    border-bottom: 1px solid #a1a09e;
}

notebook header tabs {
    padding: 0;
}

notebook tab {
    background: linear-gradient(to bottom, #dddbd8, #d0cecc);
    border: 1px solid #a1a09e;
    border-bottom: none;
    border-radius: 4px 4px 0 0;
    padding: 4px 8px;
    margin: 2px 1px 0 1px;
    color: #3c3b37;
}

notebook tab:checked {
    background: linear-gradient(to bottom, #f7f6f5, #f2f1f0);
    border-bottom: 1px solid #f2f1f0;
    color: #3c3b37;
}

notebook tab:hover:not(:checked) {
    background: linear-gradient(to bottom, #e8e6e3, #dddbd8);
}

notebook tab button {
    background: none;
    border: none;
    border-radius: 2px;
    padding: 0px;
    min-height: 16px;
    min-width: 16px;
    color: #888580;
}

notebook tab button:hover {
    background: rgba(0,0,0,0.1);
    color: #3c3b37;
}

/* === Statusbar === */
.statusbar {
    background: linear-gradient(to bottom, #e8e6e3, #dddbd8);
    border-top: 1px solid #c5c3c0;
    padding: 2px 12px;
    font-size: 12px;
    color: #555550;
    min-height: 24px;
}

/* === Context Menus === */
menu {
    background-color: #f7f6f5;
    border: 1px solid #a1a09e;
    border-radius: 0;
    padding: 4px 0;
    box-shadow: 2px 2px 6px rgba(0,0,0,0.25);
}

menu menuitem {
    padding: 4px 20px;
    color: #3c3b37;
    min-height: 22px;
}

menu menuitem:hover {
    background-color: #f07746;
    color: white;
}

menu separator {
    background-color: #c5c3c0;
    min-height: 1px;
    margin: 4px 0;
}

/* === Scrollbars === */
scrollbar {
    background-color: #e8e6e3;
}

scrollbar slider {
    background-color: #b5b3b0;
    border-radius: 4px;
    min-width: 8px;
    min-height: 8px;
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
    background: linear-gradient(to bottom, #f7f6f5, #e8e6e3);
    border: 1px solid #a1a09e;
    border-radius: 4px;
    padding: 6px 16px;
    color: #3c3b37;
}

dialog .dialog-action-area button:hover {
    background: linear-gradient(to bottom, #ffffff, #eeede9);
}

dialog .dialog-action-area button.suggested-action {
    background: linear-gradient(to bottom, #f4945e, #f07746);
    border-color: #c55a2b;
    color: white;
}

dialog .dialog-action-area button.suggested-action:hover {
    background: linear-gradient(to bottom, #f5a576, #f28c5e);
}

/* === Search Entry === */
.search-entry {
    background: #4e4d49;
    border: 1px solid #2b2a27;
    border-radius: 3px;
    color: #dfdbd2;
    padding: 2px 6px;
    min-height: 26px;
    min-width: 150px;
}

.search-entry:focus {
    border-color: #f07746;
}

/* === Paned separator === */
paned separator {
    background-color: #c5c3c0;
    min-width: 1px;
    min-height: 1px;
}

/* === View toggle === */
.view-toggle {
    background: linear-gradient(to bottom, #565550, #484743);
    border: 1px solid #2b2a27;
    border-radius: 3px;
    padding: 0;
}

.view-toggle button {
    background: none;
    border: none;
    border-radius: 2px;
    color: #dfdbd2;
    padding: 2px 6px;
    min-height: 24px;
    min-width: 28px;
    box-shadow: none;
}

.view-toggle button:checked {
    background: rgba(0,0,0,0.2);
    box-shadow: inset 0 1px 3px rgba(0,0,0,0.3);
}

.view-toggle button:hover:not(:checked) {
    background: rgba(255,255,255,0.05);
}

/* === Drop highlight === */
.drop-target {
    border: 2px dashed #f07746;
    border-radius: 4px;
}
`
