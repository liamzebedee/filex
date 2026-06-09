package ui

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/core"
	"filex/util"
)

const (
	gridIconSize = 48
	// thumbnails are only generated for images smaller than this; larger
	// files fall back to the themed icon.
	thumbMaxBytes = 20 << 20
	// at most this many thumbnails are generated per directory render so
	// huge photo directories stay responsive.
	thumbsPerRender = 200
)

var (
	iconCache  = map[string]*gdk.Pixbuf{} // themed icons by name
	thumbCache = map[string]*gdk.Pixbuf{} // thumbnails by path|mtime
)

// gridPixbufFor returns the pixbuf shown for an entry in the icon view:
// a scaled-down thumbnail for reasonably-sized images (budget permitting),
// else the themed icon for its mime type. Results are cached.
func gridPixbufFor(e core.FileEntry, thumbBudget *int) *gdk.Pixbuf {
	mime := util.MimeFor(e.Name, e.IsDir)

	if strings.HasPrefix(mime, "image/") && e.Size > 0 && e.Size < thumbMaxBytes {
		key := fmt.Sprintf("%s|%d", e.Path, e.ModTime)
		if pb, ok := thumbCache[key]; ok {
			return pb
		}
		if *thumbBudget > 0 {
			*thumbBudget--
			if pb, err := gdk.PixbufNewFromFileAtScale(e.Path, gridIconSize, gridIconSize, true); err == nil {
				if len(thumbCache) > 4096 {
					thumbCache = map[string]*gdk.Pixbuf{}
				}
				thumbCache[key] = pb
				return pb
			}
		}
	}

	return themedIcon(util.IconForMime(mime))
}

// themedIcon loads (and caches) a themed icon by name at grid size.
func themedIcon(name string) *gdk.Pixbuf {
	if pb, ok := iconCache[name]; ok {
		return pb
	}
	theme, err := gtk.IconThemeGetDefault()
	if err != nil {
		return nil
	}
	pb, err := theme.LoadIcon(name, gridIconSize, gtk.ICON_LOOKUP_USE_BUILTIN)
	if err != nil {
		// Cache the miss as the generic icon so we don't retry every render.
		pb, _ = theme.LoadIcon("text-x-generic", gridIconSize, gtk.ICON_LOOKUP_USE_BUILTIN)
	}
	iconCache[name] = pb
	return pb
}
