package main

import "path"

var vaultCacheDir = "../Mount"

func MetadataDirectory(filepath string) string {
	path := path.Join(vaultCacheDir, "data", filepath)
	return path
}
