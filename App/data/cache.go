package data

import (
	"crypto/md5"
	"fmt"

	"github.com/cowsed/Pumice/App/config"
	"github.com/cowsed/Pumice/App/parser"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
)

type VaultCache struct {
	Version config.Version `json:"version"`
	Notes   []NoteCache    `json:"notes"`
}

type NoteCache struct {
	Path     VaultLocation
	Tags     TagSet
	Outlinks []VaultLocation
	Metadata map[string]MetaDataValue
	Md5sum   [16]byte
}

func (nc NoteCache) Title() string {
	return string(nc.Path.Name())
}

type FullPath struct {
	vaultLocation OSPath
	notePath      VaultLocation
}

func (fp FullPath) ToPath() string {
	return ToOSPath(fp.vaultLocation, fp.notePath)
}

func GetOutlinks(doc ast.Node, source []byte) []string {
	var outs = mapset.NewSet[string]()
	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		var linkSource = ""
		if n, ok := node.(*ast.Link); ok && enter {
			linkSource = string(n.Destination)
		} else if n, ok := node.(*wikilink.Node); ok && enter {
			linkSource = string(n.Target)
		}
		if linkSource != "" {
			outs.Add(linkSource)
		}

		return ast.WalkContinue, nil
	})
	return outs.ToSlice()
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

type PannicedErr struct {
	err any
}

func (pe PannicedErr) Error() string {
	return fmt.Sprintf("error recovered: %+v", pe.err)
}

func MakeNoteCache(path VaultLocation, bytes []byte) (cache NoteCache, doc ast.Node, err error) {
	// defer func() {
	// if r := recover(); r != nil {
	// err = PannicedErr{r}
	// }
	// }()

	doc = parser.VaultParser().Parse(text.NewReader(bytes))

	meta := map[string]MetaDataValue{}
	for k, v := range doc.OwnerDocument().Meta() {
		meta[k] = v
	}

	unresolvedOutlinks := GetOutlinks(doc, bytes)
	resolvedOutlinks := make([]VaultLocation, len(unresolvedOutlinks))
	for i, unresolved := range unresolvedOutlinks {
		resolvedOutlinks[i] = VaultLocation(unresolved)
	}

	// log.Println("outlinks", outlinks)

	cache = NoteCache{
		Path:     path,
		Tags:     GetTags(doc),
		Outlinks: resolvedOutlinks,
		Metadata: meta,
		Md5sum:   md5.Sum(bytes),
	}

	// List the tags.

	return cache, doc, nil
}

type MetaDataValue interface{}

type Link string
type Tag string

func (t Tag) String() string {
	return string(t)
}
