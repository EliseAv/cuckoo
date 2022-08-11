package main

func main() {
	form := makeForm()
	defer form.Dispose()

	go speakTimeEvents()

	form.window.Run()
}
