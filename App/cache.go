package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
)

var cacheFolderName VaultLocation = ".cache"
var dataCacheFilename string = "data.json"
var dataCachePath VaultLocation = cacheFolderName.Append(dataCacheFilename)

type VaultCache struct {
	Version Version     `json:"version"`
	Notes   []NoteCache `json:"notes"`
}

type NoteCache struct {
	tags     TagSet
	outlinks Link
	metadata map[string]MetaDataValue
}

type FullPath struct {
	vaultLocation OSPath
	notePath      VaultLocation
}

func (fp FullPath) ToPath() string {
	return ToOSPath(fp.vaultLocation, fp.notePath)
}

func GetTags(doc ast.Node) TagSet {
	var tags TagSet = NewTagSet()

	// Tags on the inside of the doc `#tag` syntax
	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if n, ok := node.(*hashtag.Node); ok && enter {
			tags.Add(Tag(string(n.Tag)))
		}
		return ast.WalkContinue, nil
	})

	// tags from front matter
	meta := doc.OwnerDocument().Meta()

	maybeTags, exists := meta["tags"]

	if !exists {
		return tags
	}

	list, isList := maybeTags.([]interface{})
	if !isList {
		return tags
	}

	for _, maybeTag := range list {
		tag, isTag := maybeTag.(string)
		if isTag {
			tags.Add(Tag(tag))
		}
	}

	return tags
}

func MakeNoteCache(bytes []byte) (NoteCache, ast.Node, error) {

	doc := VaultParser().Parse(text.NewReader(bytes))

	nc := NoteCache{
		tags:     GetTags(doc),
		outlinks: "",
		metadata: map[string]MetaDataValue{},
	}

	// List the tags.

	return nc, doc, nil
}

type MetaDataValue interface{}

type Link string
type Tag string

func loadWorkspaceCache(vault_location OSPath) (*VaultCache, error) {
	var cache_folder OSPath = OSPath(ToOSPath(vault_location, cacheFolderName))
	err := os.MkdirAll(string(cache_folder), 0777)
	if err != nil {
		return nil, err
	}

	cachepath := ToOSPath(vault_location, dataCachePath)
	slog.Info("Trying to load cache", "cachepath", cachepath)

	f, err := os.Open(cachepath)
	if err != nil {
		return nil, err
	}

	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	dc := VaultCache{
		Version: CURRENT_VERSION,
		Notes:   []NoteCache{},
	}

	err = json.Unmarshal(bs, &dc)
	if err != nil {
		return nil, err
	}

	return &VaultCache{}, nil
}

type CacheBuildStatus struct {
	currentFile VaultLocation
	fileNumber  int
	totalFiles  int
}
