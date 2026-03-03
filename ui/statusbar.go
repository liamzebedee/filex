package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"golang.org/x/sys/unix"
)

// Statusbar shows item count, feedback messages, and free disk space.
type Statusbar struct {
	Box        *gtk.Box
	ItemLabel  *gtk.Label
	MsgLabel   *gtk.Label
	SpaceLabel *gtk.Label
	msgExpiry  time.Time
}

func NewStatusbar() *Statusbar {
	sb := &Statusbar{}
	var err error

	sb.Box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err != nil {
		log.Fatal(err)
	}
	sc, _ := sb.Box.GetStyleContext()
	sc.AddClass("statusbar")

	sb.ItemLabel, _ = gtk.LabelNew("")
	sb.ItemLabel.SetHAlign(gtk.ALIGN_START)
	sb.Box.PackStart(sb.ItemLabel, false, false, 0)

	sb.MsgLabel, _ = gtk.LabelNew("")
	sb.MsgLabel.SetHAlign(gtk.ALIGN_CENTER)
	sb.Box.SetCenterWidget(sb.MsgLabel)

	sb.SpaceLabel, _ = gtk.LabelNew("")
	sb.SpaceLabel.SetHAlign(gtk.ALIGN_END)
	sb.Box.PackEnd(sb.SpaceLabel, false, false, 0)

	return sb
}

// ShowMessage displays a temporary feedback message in the statusbar.
func (sb *Statusbar) ShowMessage(msg string) {
	sb.MsgLabel.SetText(msg)
	sb.msgExpiry = time.Now().Add(4 * time.Second)
	glib.TimeoutAdd(4000, func() bool {
		if time.Now().Before(sb.msgExpiry) {
			return false // a newer message replaced this one
		}
		sb.MsgLabel.SetText("")
		return false
	})
}

// Update refreshes the statusbar for the given tab.
func (sb *Statusbar) Update(tab *Tab) {
	if tab == nil {
		return
	}
	count := len(tab.FileView.Entries)
	if count == 1 {
		sb.ItemLabel.SetText("1 item")
	} else {
		sb.ItemLabel.SetText(fmt.Sprintf("%d items", count))
	}

	// Free space
	var stat unix.Statfs_t
	if err := unix.Statfs(tab.Path, &stat); err == nil {
		free := stat.Bavail * uint64(stat.Bsize)
		sb.SpaceLabel.SetText(fmt.Sprintf("Free space: %s", formatBytes(free)))
	} else {
		sb.SpaceLabel.SetText("")
	}
}

func formatBytes(b uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case b >= TB:
		return fmt.Sprintf("%.1f TB", float64(b)/float64(TB))
	case b >= GB:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", b)
	}
}
