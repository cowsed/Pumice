package main

import (
	"log"
	"log/slog"
	"time"
)

type BackendState int

const (
	Opening BackendState = iota
	LoadingConfig
	LoadingCache
	BuildingCache
	Ready
	Error
)

type BackendUpdate struct {
	state   BackendState
	message string
	err     error
}

func LoadWorkspace(flags Flags) chan BackendUpdate {
	updateChan := make(chan BackendUpdate, 5)
	updateChan <- BackendUpdate{
		state:   Opening,
		message: "Opening",
		err:     nil,
	}

	go loadWorkspace(flags, updateChan)

	return updateChan
}
func loadWorkspace(flags Flags, updates chan BackendUpdate) {
	time.Sleep(1 * time.Second)

	updates <- BackendUpdate{
		state:   LoadingConfig,
		message: "Loading Config",
		err:     nil,
	}
	time.Sleep(1 * time.Second)

	cfg, err := loadWorkspaceConfig(flags.VaultPath)
	if err != nil {
		slog.Error("Unable to load workspace config, using default...", "err", err)
		err = cfg.Save(flags.VaultPath)
		if err != nil {
			slog.Error("Couldn't write default config", "err", err)
		}
	}
	updates <- BackendUpdate{
		state:   LoadingCache,
		message: "Loading Cache",
		err:     nil,
	}
	time.Sleep(1 * time.Second)

	dc, err := loadWorkspaceCache(flags.VaultPath)
	if err != nil {
		slog.Warn("Failed to load cache, rebuilding...", "err", err)
	}
	log.Println("Cache", dc)

	updates <- BackendUpdate{
		state:   Ready,
		message: "Done",
		err:     nil,
	}

}
