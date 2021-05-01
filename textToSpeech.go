package main

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func speak(message string) {
	// Not sure if we'll need initialization
	err := ole.CoInitialize(0)
	if err == nil {
		defer ole.CoUninitialize()
	}

	object, err := oleutil.CreateObject("SAPI.SpVoice")
	if err != nil {
		log.Panic(err)
	}
	sapiVoice, err := object.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Panic(err)
	}

	_, err = oleutil.CallMethod(sapiVoice, "Speak", message)
	if err != nil {
		log.Fatal(err)
	}
}
