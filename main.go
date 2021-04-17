package main

func main() {
	form := makeForm()
	defer form.Dispose()

	channel := make(chan string)
	go speak(channel)
	go emitEnglishSpeakEvents(channel)

	form.window.Run()
}
