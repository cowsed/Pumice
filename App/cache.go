package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
)

var cacheFolderName VaultLocation = ".cache"
var dataCacheFilename string = "data.json"
var dataCachePath VaultLocation = cacheFolderName.Append(dataCacheFilename)

type VaultCache struct {
	Version Version     `json:"version"`
	Notes   []NoteCache `json:"notes"`
}

type NoteCache struct {
	notepath VaultLocation
}

type NoCache struct {
	err error
}

func (nc NoCache) Error() string {
	return "No Cache Found: " + nc.err.Error()
}

func loadWorkspaceCache(vault_location OSPath) (*VaultCache, error) {
	var cache_folder OSPath = OSPath(ToOSPath(vault_location, cacheFolderName))
	err := os.MkdirAll(string(cache_folder), 0644)
	if err != nil {
		return nil, err
	}

	cachepath := ToOSPath(vault_location, dataCachePath)
	slog.Info("Trying to load cache", "cachepath", cachepath)

	f, err := os.Open(cachepath)
	if err != nil {
		return nil, NoCache{err}
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, NoCache{err}
	}

	dc := VaultCache{
		Version: CURRENT_VERSION,
		Notes:   []NoteCache{},
	}

	err = json.Unmarshal(bs, &dc)
	if err != nil {
		return nil, NoCache{err}
	}

	return &VaultCache{}, nil
}

type CacheBuildStatus struct {
	currentFile VaultLocation
	fileNumber  int
	totalFiles  int
}
