package main

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	Vaults     []string `json:"vaults"`
	OpenVaults []string `json:"open_vaults"`
}

var DefaultAppConfig AppConfig = AppConfig{
	Vaults:     []string{},
	OpenVaults: []string{},
}

func (cfg AppConfig) Save(path string) error {
	bs, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bs, 0644)
}
