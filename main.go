package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

var AppConfigDir = "config"

func LoadConfig() AppConfig {
	appcfg_path := path.Join(AppConfigDir, "app.json")
	file, err := os.Open(appcfg_path)
	var cfg AppConfig = DefaultAppConfig

	if os.IsNotExist(err) {
		slog.Warn("Can not find app config. Using defaults...")
		cfg.Save(appcfg_path)
		return cfg
	} else if err != nil {
		slog.Error("Error openning config file. Using defaults...", "err", err)
		return cfg
	}

	bs, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Error reading app config file. Using defaults...", "err", err)
		return cfg
	}

	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		slog.Error("Failed to parse app config file. Using defaults...", "err", err)
		return cfg
	}

	return cfg
}

type Notification struct {
	Message string
}

func notificationHandler(queue chan Notification) {
	for elem := range queue {
		slog.Info("Notification", "noti", elem)
	}
}

func main() {
	cfg := LoadConfig()
	fmt.Println(cfg)

	if len(cfg.OpenVaults) == 0 {
		ShowVaultOpener()
	}

}

func ShowVaultOpener() {
	noti_chan := make(chan Notification)
	go notificationHandler(noti_chan)

	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")

	title := canvas.NewText("Pumice", color.Black)
	title.TextSize = 18
	title.Alignment = fyne.TextAlignCenter

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewButton("Open Vault", func() { LoadVaultToImport(noti_chan, myWindow) }),
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()

}

func LoadVaultToImport(notis chan Notification, window fyne.Window) {
	filename, err := PickVaultToImport(notis)

	if errors.Is(err, dialog.ErrCancelled) {
		slog.Debug("Vault import cancelled. File picker closed")
		return
	} else if err != nil {
		slog.Error("Vault import failed. File picker error. ", "err", err)
		notis <- Notification{
			Message: "Failed to open vault: " + err.Error(),
		}
		return
	}
	log.Println(filename)
	window.Close()

}

func PickVaultToImport(notis chan Notification) (string, error) {
	filename, err := dialog.Directory().Title("Open a Vault").Browse()
	if err != nil {
		return "", err
	}
	return filename, nil
}
