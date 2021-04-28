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

//go:embed *.png
var icons embed.FS

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
	gui := &guiView{}

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
	gui.tray.SetVisible(true)
	gui.tray.MouseDown().Attach(gui.onTrayMouse)

	gui.tiEnabled = newIconFromEmbeddedFilename("imgyeah.png")
	gui.tiDisabled = newIconFromEmbeddedFilename("imgok.png")

	menu := gui.tray.ContextMenu()

	gui.active = walk.NewAction()
	gui.active.SetText("&Active")
	gui.active.SetCheckable(true)
	gui.active.Triggered().Attach(func() { gui.setActive(gui.active.Checked()) })
	menu.Actions().Add(gui.active)
	gui.setActive(settings.Active)

	menu.Actions().Add(walk.NewSeparatorAction())

	gui.intervals = []struct {
		action *walk.Action
		value  int
		label  string
	}{
		{walk.NewAction(), 1, "Every &minute"},
		{walk.NewAction(), 15, "1&5 minutes"},
		{walk.NewAction(), 30, "&30 minutes"},
		{walk.NewAction(), 60, "1 &hour"},
	}
	for _, item := range gui.intervals {
		item.action.SetText(item.label)
		item.action.SetCheckable(true)
		item.action.Triggered().Attach( gui.setInterval(item.value) )
		menu.Actions().Add(item.action)
	}
	gui.setInterval(settings.IntervalMinutes)()

	menu.Actions().Add(walk.NewSeparatorAction())

	action := walk.NewAction()
	action.SetText("&Quit")
	action.Triggered().Attach(func() { walk.App().Exit(0) })
	menu.Actions().Add(action)

	return gui
}

func (gui *guiView) Dispose() {
	if tray := gui.tray; tray != nil {
		tray.Dispose()
	}
	if window := gui.window; window != nil {
		window.Dispose()
	}
}

func (gui *guiView) onTrayMouse(x, y int, button walk.MouseButton) {
	if button == walk.LeftButton {
		gui.window.Menu() // TODO: figure out how to open a damn menu
	}
}

func (gui *guiView) setActive(value bool) {
	settings.Active = value
	gui.active.SetChecked(value)
	if settings.Active {
		gui.tray.SetIcon(gui.tiEnabled)
	} else {
		gui.tray.SetIcon(gui.tiDisabled)
	}
	settings.Save()
}

func (gui guiView) setInterval(value int) walk.EventHandler {
	return func() {
		fmt.Println("Setting interval to", value)
		settings.IntervalMinutes = value
		for _, info := range gui.intervals {
			info.action.SetChecked(info.value == value)
		}
		settings.Save()
	}
}

func newIconFromEmbeddedFilename(filename string) walk.Image {
	payload, err := icons.ReadFile(filename)
	if err != nil {
		log.Panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(payload))
	if err != nil {
		log.Panic(err)
	}
	icon, err := walk.NewBitmapFromImageForDPI(img, 96)
	if err != nil {
		log.Panic(err)
	}
	return icon
}
