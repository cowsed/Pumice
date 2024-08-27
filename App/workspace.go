package main

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/cowsed/Pumice/App/config"
	"github.com/cowsed/Pumice/App/data"
)

var cacheFolderName data.VaultLocation = ".cache"
var dataCacheFilename string = "data.json"
var dataCachePath data.VaultLocation = cacheFolderName.Append(dataCacheFilename)

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

func loadWorkspaceCache(vault_location data.OSPath) (*data.VaultCache, error) {
	var cache_folder data.OSPath = data.OSPath(data.ToOSPath(vault_location, cacheFolderName))
	err := os.MkdirAll(string(cache_folder), 0777)
	if err != nil {
		return nil, err
	}

	cachepath := data.ToOSPath(vault_location, dataCachePath)
	canon, err := filepath.Abs(cachepath)
	if err != nil {
		return nil, err
	}
	slog.Info("Trying to load cache", "cachepath", canon)

	f, err := os.Open(cachepath)
	if err != nil {
		return nil, err
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	dc := data.VaultCache{
		Version: config.VERSION,
		Notes:   []data.NoteCache{},
	}

	err = json.Unmarshal(bs, &dc)
	if err != nil {
		return nil, err
	}

	return &data.VaultCache{}, nil
}
