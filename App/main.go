package main

import (
	"fmt"
	"log/slog"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func doUI(vaultPath OSPath, cfg Config, update chan BackendUpdate, win fyne.Window) {
	vaultName := vaultPath.Base()
	label := widget.NewLabel(fmt.Sprintf("Openning %s...", vaultName))
	progress := widget.NewProgressBarInfinite()
	progress.Start()

	win.SetContent(container.NewVBox(label, progress))
	for update := range update {
		label.SetText(update.Describe())
	}
}

func handleConfigLoad(bu BackendUpdate) Config {
	cfgUpdate, ok := bu.(ConfigLoaded)
	if !ok {
		panic("Application error")
	}
	return cfgUpdate.config
}
func main() {
	flags := parseFlags()
	slog.Info("Loaded flags", "flags", flags)

	updates := LoadWorkspace(flags)

	myApp := app.New()

	var update BackendUpdate = <-updates

	initialCfg := handleConfigLoad(update)
	win := myApp.NewWindow(fmt.Sprintf("Pumice - %s", path.Base(string(flags.VaultPath))))

	go doUI(flags.VaultPath, initialCfg, updates, win)

	win.ShowAndRun()

}
