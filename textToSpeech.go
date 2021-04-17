package main

import (
	"log"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func speak(eventSource chan string) {
	// Not sure if we'll need initialization
	err := ole.CoInitialize(0)
	if err == nil {
		defer ole.CoUninitialize()
	}
	object, err := oleutil.CreateObject("SAPI.SpVoice")
	if err != nil {
		panic(err)
	}
	sapiVoice, err := object.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		panic(err)
	}
	for message := range eventSource {
		_, err = oleutil.CallMethod(sapiVoice, "Speak", message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
