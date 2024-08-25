package main

import (
	"log"
	"log/slog"
)

type BackendState int

const (
	Opening BackendState = iota
	LoadingConfig
	LoadingThemes
	LoadingExtensions
	LoadingCache
	BuildingCache
	Ready
	Error
)

type BackendUpdate interface {
	Describe() string
}

var _ BackendUpdate = ConfigLoaded{}

type ConfigLoaded struct {
	config Config
}

func (cl ConfigLoaded) Describe() string {
	return "Loaded Configuration"
}

func LoadWorkspace(flags Flags) chan BackendUpdate {
	updateChan := make(chan BackendUpdate, 5)

	go loadWorkspace(flags, updateChan)

	return updateChan
}
func loadWorkspace(flags Flags, updates chan BackendUpdate) {
	cfg, err := loadWorkspaceConfig(flags.VaultPath)
	if err != nil {
		slog.Error("Unable to load workspace config, using default...", "err", err)
		err = cfg.Save(flags.VaultPath)
		if err != nil {
			slog.Error("Couldn't write default config", "err", err)
		}
	}

	updates <- ConfigLoaded{
		config: cfg,
	}

	dc, err := loadWorkspaceCache(flags.VaultPath)
	if err != nil {
		slog.Warn("Failed to load cache, rebuilding...", "err", err)
	}
	log.Println("Cache", dc)

}
