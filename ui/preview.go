package ui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/util"
)

// quickLook is the single spacebar-toggled preview window.
type quickLook struct {
	win  *gtk.Window
	path string
}

var preview quickLook

const previewTextLimit = 256 << 10 // bytes of a text file shown

// ToggleQuickLook previews path in a transient window. Pressing space on
// the same selection closes it; a different selection replaces it.
func ToggleQuickLook(app *App, path string) {
	samePath := preview.path == path
	if preview.win != nil {
		preview.win.Destroy() // the destroy handler resets the state
	}
	if samePath {
		return
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return
	}
	win.SetTitle(filepath.Base(path))
	win.SetTransientFor(app.Window)
	win.SetDestroyWithParent(true)
	win.SetPosition(gtk.WIN_POS_CENTER_ON_PARENT)

	content, wantLarge := previewContent(path)
	if wantLarge {
		win.SetDefaultSize(720, 520)
	}
	win.Add(content)

	// Space or Escape closes, mirroring how it was opened.
	win.Connect("key-press-event", func(w *gtk.Window, event *gdk.Event) bool {
		key := gdk.EventKeyNewFromEvent(event).KeyVal()
		if key == gdk.KEY_space || key == gdk.KEY_Escape {
			w.Destroy()
			return true
		}
		return false
	})
	win.Connect("destroy", func() {
		preview = quickLook{}
	})

	preview = quickLook{win: win, path: path}
	win.ShowAll()
}

// previewContent builds the preview widget for path: a scaled image, a
// read-only text view, or an info panel as fallback. wantLarge reports
// whether the window should open at full preview size.
func previewContent(path string) (gtk.IWidget, bool) {
	info, err := os.Stat(path)
	if err != nil {
		label, _ := gtk.LabelNew("Cannot preview: " + err.Error())
		label.SetMarginTop(24)
		label.SetMarginBottom(24)
		label.SetMarginStart(24)
		label.SetMarginEnd(24)
		return label, false
	}

	mime := util.MimeFor(info.Name(), info.IsDir())

	if strings.HasPrefix(mime, "image/") {
		if pb, err := gdk.PixbufNewFromFileAtScale(path, 700, 500, true); err == nil {
			img, err := gtk.ImageNewFromPixbuf(pb)
			if err == nil {
				return img, true
			}
		}
	}

	if util.IsTextMime(mime) && !info.IsDir() {
		if widget, ok := textPreview(path); ok {
			return widget, true
		}
	}

	return infoPreview(path, info, mime), false
}

// textPreview shows the head of a text file in a monospace view.
func textPreview(path string) (gtk.IWidget, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, previewTextLimit))
	if err != nil {
		return nil, false
	}

	tv, err := gtk.TextViewNew()
	if err != nil {
		return nil, false
	}
	tv.SetEditable(false)
	tv.SetCursorVisible(false)
	tv.SetMonospace(true)
	tv.SetLeftMargin(8)
	tv.SetRightMargin(8)
	buf, err := tv.GetBuffer()
	if err != nil {
		return nil, false
	}
	text := strings.ToValidUTF8(string(data), "�")
	if int64(len(data)) == previewTextLimit {
		text += "\n…"
	}
	buf.SetText(text)

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, false
	}
	scroll.Add(tv)
	return scroll, true
}

// infoPreview shows a small icon-and-facts panel for anything that has no
// richer preview.
func infoPreview(path string, info os.FileInfo, mime string) gtk.IWidget {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
	box.SetMarginTop(24)
	box.SetMarginBottom(24)
	box.SetMarginStart(36)
	box.SetMarginEnd(36)

	icon, _ := gtk.ImageNewFromIconName(util.IconForMime(mime), gtk.ICON_SIZE_DIALOG)
	box.PackStart(icon, false, false, 0)

	name, _ := gtk.LabelNew(info.Name())
	box.PackStart(name, false, false, 0)

	detail := mime
	if info.IsDir() {
		if entries, err := os.ReadDir(path); err == nil {
			detail = itemCount(len(entries))
		} else {
			detail = "Folder"
		}
	} else {
		detail = fmt.Sprintf("%s — %s", mime, util.FormatSize(info.Size()))
	}
	detailLabel, _ := gtk.LabelNew(detail)
	box.PackStart(detailLabel, false, false, 0)

	modified, _ := gtk.LabelNew("Modified " + util.FormatDate(info.ModTime().Local().Truncate(time.Second)))
	box.PackStart(modified, false, false, 0)

	return box
}
