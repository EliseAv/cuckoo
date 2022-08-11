package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

//go:embed windows.manifest *.png
var resources embed.FS

type guiView struct {
	window    *walk.MainWindow
	active    *walk.Action
	intervals []struct {
		action *walk.Action
		value  int
		label  string
	}
	tray                  *walk.NotifyIcon
	tiEnabled, tiDisabled walk.Image
}

func makeForm() *guiView {
	var gui guiView
	err := d.MainWindow{
		AssignTo: &gui.window,
		Title:    "Cuckoo",
		Visible:  false,
	}.Create()
	if err != nil {
		log.Panic(err)
	}

	if gui.tray, err = walk.NewNotifyIcon(gui.window); err != nil {
		log.Panic(err)
	}
	_ = gui.tray.SetVisible(true)
	gui.tray.MouseDown().Attach(gui.onTrayMouse)

	gui.tiEnabled = newIconFromEmbeddedFilename("imgyeah.png")
	gui.tiDisabled = newIconFromEmbeddedFilename("imgok.png")

	menu := gui.tray.ContextMenu()

	gui.active = walk.NewAction()
	_ = gui.active.SetText("&Active")
	_ = gui.active.SetCheckable(true)
	gui.active.Triggered().Attach(func() { gui.setActive(gui.active.Checked()) })
	_ = menu.Actions().Add(gui.active)
	gui.setActive(settings.Active)

	_ = menu.Actions().Add(walk.NewSeparatorAction())

	gui.intervals = []struct {
		action *walk.Action
		value  int
		label  string
	}{
		{walk.NewAction(), 1, "Every &minute"},
		{walk.NewAction(), 5, "&5 minutes"},
		{walk.NewAction(), 10, "&10 minutes"},
		{walk.NewAction(), 15, "15 m&inutes"},
		{walk.NewAction(), 30, "&30 minutes"},
		{walk.NewAction(), 60, "1 &hour"},
	}
	for _, item := range gui.intervals {
		_ = item.action.SetText(item.label)
		_ = item.action.SetCheckable(true)
		item.action.Triggered().Attach(gui.setInterval(item.value))
		_ = menu.Actions().Add(item.action)
	}
	gui.setInterval(settings.IntervalMinutes)()

	_ = menu.Actions().Add(walk.NewSeparatorAction())

	action := walk.NewAction()
	_ = action.SetText("E&xit")
	action.Triggered().Attach(func() { walk.App().Exit(0) })
	_ = menu.Actions().Add(action)

	return &gui
}

func (gui *guiView) Dispose() {
	if tray := gui.tray; tray != nil {
		_ = tray.Dispose()
	}
	if window := gui.window; window != nil {
		window.Dispose()
	}
}

func (gui *guiView) onTrayMouse(_, _ int, button walk.MouseButton) {
	if button == walk.LeftButton {
		// TODO: figure out how to open a damn menu
		gui.setActive(!settings.Active)

		message := "Time notifications are now disabled."
		if settings.Active {
			s := "s"
			if settings.IntervalMinutes == 1 {
				s = ""
			}
			message = fmt.Sprintf("Speaking the time every %d minute%s.", settings.IntervalMinutes, s)
		}

		_ = gui.tray.ShowInfo("Cuckoo", message+"\nRight click to change settings or exit")
	}
}

func (gui *guiView) setActive(value bool) {
	settings.Active = value
	_ = gui.active.SetChecked(value)
	if settings.Active {
		_ = gui.tray.SetIcon(gui.tiEnabled)
	} else {
		_ = gui.tray.SetIcon(gui.tiDisabled)
	}
	_ = settings.Save()
}

func (gui guiView) setInterval(value int) walk.EventHandler {
	return func() {
		settings.IntervalMinutes = value
		for _, info := range gui.intervals {
			_ = info.action.SetChecked(info.value == value)
		}
		_ = settings.Save()
	}
}

func newIconFromEmbeddedFilename(filename string) walk.Image {
	payload, err := resources.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}

	var img image.Image
	img, _, err = image.Decode(bytes.NewReader(payload))
	if err != nil {
		log.Panic(err)
	}

	var icon walk.Image
	icon, err = walk.NewBitmapFromImageForDPI(img, 96)
	if err != nil {
		log.Panic(err)
	}
	return icon
}
