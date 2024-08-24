package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var configFolderName VaultLocation = ".config"
var configFilename string = "config.json"
var configFilePath VaultLocation = configFolderName.Append(configFilename)

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.patch, v.comment)
}

type NoConfig struct {
	err error
}

func (nc NoConfig) Error() string {
	return "Unable to load config: " + nc.err.Error()
}

func loadWorkspaceConfig(vault_path OSPath) (Config, error) {
	path := ToOSPath(vault_path, configFilePath)
	default_cfg := NewConfig()

	f, err := os.Open(path)
	if err != nil {
		return default_cfg, NoConfig{err}
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return default_cfg, NoConfig{err}
	}

	cfg := Config{}

	err = json.Unmarshal(bs, &cfg)
	if err != nil {
		return default_cfg, NoConfig{err}
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
	err := os.MkdirAll(string(config_folder), 0644)

	if err != nil {
		return err
	}

	bs, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(string(vault_location), bs, 0644)
}

type ThemeID string
type Theme struct {
	name string
}
