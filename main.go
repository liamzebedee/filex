package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"filex/theme"
	"filex/ui"
)

func main() {
	gtk.Init(&os.Args)

	// Load Ambiance CSS
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Failed to create CSS provider:", err)
	}
	if err := cssProvider.LoadFromData(theme.AmbianceCSS); err != nil {
		log.Fatal("Failed to load CSS:", err)
	}
	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Fatal("Failed to get default screen:", err)
	}
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	app := ui.NewApp()
	app.Window.ShowAll()

	gtk.Main()
}
