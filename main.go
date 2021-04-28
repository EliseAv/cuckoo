package main

func main() {
	form := makeForm()
	defer form.Dispose()

	channel := make(chan string)
	go emitEnglishSpeechEvents(channel)
	go absorbSpeechEvents(channel)

	form.window.Run()
}
