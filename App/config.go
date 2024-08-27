package main

import (
	"encoding/json"
	"io"
	"os"

	"fyne.io/fyne/v2"
	"github.com/cowsed/Pumice/App/data"
)

var configFolderName data.VaultLocation = ".config"
var configFilename string = "config.json"
var configFilePath data.VaultLocation = configFolderName.Append(configFilename)

func loadWorkspaceConfig(vault_path data.OSPath) (Config, error) {
	path := data.ToOSPath(vault_path, configFilePath)
	default_cfg := NewConfig()

	f, err := os.Open(path)
	if err != nil {
		return default_cfg, err
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return default_cfg, err
	}

	cfg := Config{}

	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		return default_cfg, err
	}

	return cfg, nil
}

func NewConfig() Config {
	return Config{
		Themes:       []Theme{},
		CurrentTheme: "builtin",
		WindowSize:   fyne.NewSize(400, 300),
		// extensions: []Extension{}
	}
}

type Config struct {
	Themes       []Theme   `json:"theme"`
	CurrentTheme ThemeID   `json:"current_theme"`
	WindowSize   fyne.Size `json:"size"`
}

func (c Config) Save(vault_location data.OSPath) error {
	var config_folder data.OSPath = data.OSPath(data.ToOSPath(vault_location, configFolderName))
	err := os.MkdirAll(string(config_folder), 0777)
	if err != nil {
		return err
	}

	configpath := data.ToOSPath(vault_location, configFilePath)

	bs, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configpath, bs, 0644)
}
