package main

import (
	"encoding/json"
	"io"
	"os"
)

var configFolderName VaultLocation = ".config"
var configFilename string = "config.json"
var configFilePath VaultLocation = configFolderName.Append(configFilename)

func loadWorkspaceConfig(vault_path OSPath) (Config, error) {
	path := ToOSPath(vault_path, configFilePath)
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
		// extensions: []Extension{}
	}
}

type Config struct {
	Themes       []Theme `json:"theme`
	CurrentTheme ThemeID `json:"current_theme`
}

func (c Config) Save(vault_location OSPath) error {
	var config_folder OSPath = OSPath(ToOSPath(vault_location, configFolderName))
	err := os.MkdirAll(string(config_folder), 0777)
	if err != nil {
		return err
	}

	configpath := ToOSPath(vault_location, configFilePath)

	bs, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configpath, bs, 0)
}
