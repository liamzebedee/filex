package ui

import (
	"log"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"golang.org/x/sys/unix"

	"filex/util"
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
	if sb == nil {
		return
	}
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

// Render shows the visible item count and the free space for path.
func (sb *Statusbar) Render(path string, count int) {
	if sb == nil {
		return
	}
	sb.ItemLabel.SetText(itemCount(count))

	var stat unix.Statfs_t
	if err := unix.Statfs(path, &stat); err == nil {
		free := int64(stat.Bavail) * stat.Bsize
		sb.SpaceLabel.SetText("Free space: " + util.FormatSize(free))
	} else {
		sb.SpaceLabel.SetText("")
	}
}
