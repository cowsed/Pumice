package main

import (
	"log"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func doUI(update chan BackendUpdate, win fyne.Window) {
	label := widget.NewLabel("Openning...")
	progress := widget.NewProgressBarInfinite()
	progress.Start()

	win.SetContent(container.NewVBox(label, progress))
	for update := range update {
		label.SetText(update.message)
	}
}

func main() {
	flags := parseFlags()
	slog.Info("Loaded flags", "flags", flags)

	updates := LoadWorkspace(flags)
	log.Println(updates)

	myApp := app.New()
	win := myApp.NewWindow("Pumice")

	go doUI(updates, win)

	win.ShowAndRun()

}
