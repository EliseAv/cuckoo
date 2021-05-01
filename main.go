package main

func main() {
	form := makeForm()
	defer form.Dispose()

	go emitEnglishSpeechEvents()

	form.window.Run()
}
