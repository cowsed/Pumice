package main

import (
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/cowsed/Pumice/App/data"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Vault struct {
	filesys fs.FS
	caches  []data.NoteCache

	fileSearchContents []data.VaultLocation
	fileSearchTerm     string
}

func (v *Vault) AllFiles() []data.VaultLocation {
	ls := make([]data.VaultLocation, len(v.caches))
	for i, l := range v.caches {
		ls[i] = l.Path
	}
	return ls
}
func NewVaultFromFS(filesys fs.FS) Vault {
	v := Vault{
		filesys: filesys,
	}
	mds, err := allFilesOfType(v.filesys, ".md")
	if err != nil {
		panic(err)
	}

	log.Println("There are ", len(mds), "markdown files here")
	v.caches = CacheAll(mds, v.filesys)
	log.Printf("Read %v of %v files", len(v.caches), len(mds))
	v.fileSearchTerm = ""
	v.SearchFiles(v.fileSearchTerm)

	return v
}
func NewVaultFromCacheAndFS(filesys fs.FS) Vault {
	panic("Unimplemented")
}

func (v *Vault) Trace(rest ...interface{}) {
	log.Println(rest...)
}
func (v *Vault) TraceF(format string, rest ...interface{}) {
	log.Printf(format, rest...)
}

func (v *Vault) ResolveFile(link string) (data.VaultLocation, error) {
	return "", fmt.Errorf("could not resolve %s (unimplemented)", link)
}
func (v *Vault) SearchFiles(search string) []data.VaultLocation {
	search = strings.TrimSpace(search)
	v.fileSearchTerm = search

	if search == "" {
		v.fileSearchContents = v.AllFiles()
		return v.fileSearchContents
	}

	strs := make([]string, len(v.caches))
	for i, l := range v.caches {
		strs[i] = string(l.Path)
	}

	ranks := fuzzy.RankFind(search, strs)
	sort.Sort(ranks)

	vls := []data.VaultLocation{}
	for _, r := range ranks {
		vls = append(vls, v.caches[r.OriginalIndex].Path)
	}
	v.fileSearchContents = vls
	return vls
}
