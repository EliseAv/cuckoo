package main

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/lxn/walk"
	d "github.com/lxn/walk/declarative"
)

//go:embed *.ico *.png
var icons embed.FS

type formWidgets struct {
	window                *walk.MainWindow
	checkBox              *walk.CheckBox
	radioButtons          [3]*walk.RadioButton
	tray                  *walk.NotifyIcon
	tiEnabled, tiDisabled walk.Image
}

func makeForm() *formWidgets {
	form := &formWidgets{}
	err := d.MainWindow{
		AssignTo: &form.window,
		Title:    "Cuckoo",
		Size:     d.Size{Width: 400, Height: 300},
		Layout:   d.VBox{},
		Children: []d.Widget{
			d.CheckBox{
				AssignTo:            &form.checkBox,
				Checked:             true,
				Text:                "&Active",
				OnCheckStateChanged: form.checkStateChanged,
			},
			d.RadioButtonGroupBox{
				Title:  "Interval",
				Layout: d.VBox{},
				Buttons: []d.RadioButton{
					{
						Text:      "15 minutes",
						OnClicked: form.setInterval(15),
						AssignTo:  &form.radioButtons[0],
						Value:     15,
					},
					{
						Text:      "30 minutes",
						OnClicked: form.setInterval(30),
						AssignTo:  &form.radioButtons[1],
						Value:     30,
					},
					{
						Text:      "1 hour",
						OnClicked: form.setInterval(60),
						AssignTo:  &form.radioButtons[2],
						Value:     60,
					},
				},
			},
		},
	}.Create()
	if err != nil {
		panic(err)
	}
	form.window.SetVisible(false)

	form.radioButtons[0].SetChecked(true)

	form.tray, err = walk.NewNotifyIcon(form.window)
	if err != nil {
		panic(err)
	}

	form.tiEnabled = newIconFromEmbeddedFilename("imgyeah.png")
	form.tiDisabled = newIconFromEmbeddedFilename("imgok.png")
	form.tray.SetIcon(form.tiEnabled)

	if err := form.tray.SetVisible(true); err != nil {
		panic(err)
	}

	form.tray.MouseDown().Attach(form.mouseDown)
	return form
}

func (form *formWidgets) Dispose() {
	if form.tray != nil {
		form.tray.Dispose()
	}
	if form.window != nil {
		form.window.Dispose()
	}
}

func (form *formWidgets) checkStateChanged() {
	settings.Active = form.checkBox.Checked()
	if form.tray != nil {
		if settings.Active {
			form.tray.SetIcon(form.tiEnabled)
		} else {
			form.tray.SetIcon(form.tiDisabled)
		}
	}
	settings.Save()
}

func (form *formWidgets) mouseDown(x, y int, button walk.MouseButton) {
	if button == walk.LeftButton {
		form.checkBox.SetChecked(!form.checkBox.Checked())
	} else if button == walk.RightButton {
		form.window.SetVisible(!form.window.Visible())
	}
}

func (form formWidgets) setInterval(value int) walk.EventHandler {
	return func() {
		settings.IntervalMinutes = value
		settings.Save()
	}
}

func newIconFromEmbeddedFilename(filename string) walk.Image {
	payload, err := icons.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	icon, err := walk.NewBitmapFromImageForDPI(img, 96)
	if err != nil {
		panic(err)
	}
	return icon
}
