package main

import (
	"image/color"
	"log"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	flags := parseFlags()
	slog.Info("Loaded flags", "flags", flags)

	cfg, err := loadWorkspaceConfig(flags.VaultPath)
	if err != nil {
		slog.Error("Unable to load workspace config, using default...", "err", err)
		cfg.Save(flags.VaultPath)
	}
	log.Println(cfg)

	dc, err := loadWorkspaceCache(flags.VaultPath)
	if err != nil {
		slog.Warn("Failed to load cache, rebuilding", "err", err)
	}

	log.Println("dc", dc)

	return
	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")

	title := canvas.NewText("Pumice", color.Black)
	title.TextSize = 18
	title.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewButton("Open Vault", func() {}),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()

}
