package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/cowsed/Pumice/App/config"
	"github.com/cowsed/Pumice/App/data"
	"github.com/knusbaum/go9p"
	fs9p "github.com/knusbaum/go9p/fs"
	"github.com/knusbaum/go9p/proto"
)

func (v *Vault) Server() go9p.Srv {
	v.Trace("Starting server")

	vaultCache := makeVaultCacheFS(v.caches, v)

	return vaultCache.Server()
}

func StringsFile(links []string) func() []byte {
	return func() []byte {
		b := strings.Builder{}
		for _, l := range links {
			b.Write([]byte(l))
			b.WriteByte('\n')
		}
		return []byte(b.String())
	}
}
func LinksFile(links []data.VaultLocation) func() []byte {
	ls := make([]string, len(links))
	for i, l := range links {
		ls[i] = string(l)
	}
	return StringsFile(ls)
}

func makeAboutDir(filesys *fs9p.FS) fs9p.Dir {
	version := fs9p.NewStaticFile(filesys.NewStat("version", User, Group, 0444), []byte(config.VERSION.String()+"\n"))
	dir := fs9p.NewStaticDir(filesys.NewStat("about", User, Group, 0755))

	dir.AddChild(version)
	return dir
}

func makeDirFromCache(cache data.NoteCache, filesys *fs9p.FS) *fs9p.StaticDir {
	dir := fs9p.NewStaticDir(filesys.NewStat(string(cache.Path.Name()), User, Group, 0755))
	tags := fs9p.NewDynamicFile(filesys.NewStat("tags", User, Group, 0444), StringsFile(cache.Tags.StringList()))
	dir.AddChild(tags)

	outlinks := fs9p.NewDynamicFile(filesys.NewStat("outlinks", User, Group, 0444), func() []byte {
		buf := bytes.Buffer{}
		for _, link := range cache.Outlinks {
			buf.WriteString(string(link))
			buf.WriteByte('\n')
		}
		return buf.Bytes()
	})
	dir.AddChild(outlinks)

	metadata := fs9p.NewDynamicFile(filesys.NewStat("metadata", User, Group, 0444), func() []byte {
		bs, err := json.MarshalIndent(cache.Metadata, "", "  ")
		if err != nil {
			log.Println("Error marshalling", err)
		}
		return bs
	})
	dir.AddChild(metadata)

	sum := fs9p.NewDynamicFile(filesys.NewStat("sum", User, Group, 0444), func() []byte {
		return []byte(hexdump(cache.Md5sum[:]))
	})
	dir.AddChild(sum)

	return dir
}

func hexdump(sum []byte) string {
	s := ""
	for _, b := range sum {
		s += fmt.Sprintf("%02x", b)
	}
	s += "\n"
	return s
}

var User = "glenda"
var Group = User

type FSSTate struct {
	fs        *fs9p.FS
	cachedirs map[data.VaultLocation]*fs9p.StaticDir
	dataRoot  *fs9p.StaticDir
}

func (ft *FSSTate) GetOrMakeDir(path data.VaultLocation) *fs9p.StaticDir {
	if path == "." {
		return ft.dataRoot
	}
	if dir, exists := ft.cachedirs[path]; exists {
		return dir
	}
	name := path.Name()
	parentDir := path.Dir()
	parent := ft.GetOrMakeDir(parentDir)
	me := fs9p.NewStaticDir(ft.fs.NewStat(string(name), User, Group, 0755))
	parent.AddChild(me)
	ft.cachedirs[path] = me
	return me
}

func makeDataDir(caches []data.NoteCache, filesys *fs9p.FS) fs9p.Dir {
	dir := fs9p.NewStaticDir(filesys.NewStat("data", User, Group, 0755))
	vfst := FSSTate{
		fs:        filesys,
		cachedirs: map[data.VaultLocation]*fs9p.StaticDir{},
		dataRoot:  dir,
	}

	for _, cache := range caches {
		parentPath := cache.Path.Dir()

		noteDir := makeDirFromCache(cache, filesys)

		parentDir := vfst.GetOrMakeDir(parentPath)
		parentDir.AddChild(noteDir)

	}
	return dir
}

func makeVaultCacheFS(caches []data.NoteCache, vault *Vault) *fs9p.FS {
	vfs, root := fs9p.NewFS(User, Group, 0755)

	AboutDir := makeAboutDir(vfs)
	DataDir := makeDataDir(caches, vfs)

	ActionDir := fs9p.NewStaticDir(vfs.NewStat("actions", User, Group, 0755))
	searchContentsFile := fs9p.NewDynamicFile(vfs.NewStat("search_contents", User, Group, 0444), func() []byte { return []byte("coming soon\n") })
	// searchFile := fs9p.NewDynamicFile(vfs.NewStat("search_files", User, Group, 0444), func() []byte { return []byte("coming soon\n") })
	baseFile := fs9p.NewBaseFile(vfs.NewStat("search_files", User, Group, 0666))
	var sf fs9p.File = &fs9p.WrappedFile{
		File: baseFile,
		OpenF: func(fid uint64, omode proto.Mode) error {
			return baseFile.Open(fid, omode)
		},
		ReadF: func(fid uint64, offset uint64, count uint64) ([]byte, error) {
			s := fmt.Sprintf("%s\n%d\n", vault.fileSearchTerm, len(vault.fileSearchContents))
			for _, path := range vault.fileSearchContents {
				s += string(path) + "\n"
			}

			r := bytes.NewReader([]byte(s))
			if offset >= uint64(r.Len()) {
				return []byte{}, nil
			}

			buf := make([]byte, count)
			_, err := r.ReadAt(buf, int64(offset))
			if err == nil || errors.Is(err, io.EOF) {
				return buf, nil
			}
			return buf, err
		},
		WriteF: func(fid uint64, offset uint64, data []byte) (uint32, error) {
			vault.fileSearchTerm = string(data)
			vault.SearchFiles(vault.fileSearchTerm)
			return uint32(len(data)), nil
		},
		CloseF: func(fid uint64) error {
			return baseFile.Close(fid)
		},
	}
	ActionDir.AddChild(searchContentsFile)
	ActionDir.AddChild(sf)

	lf := fs9p.NewDynamicFile(vfs.NewStat("all_files", User, Group, 0444), LinksFile(vault.AllFiles()))
	ActionDir.AddChild(lf)

	root.AddChild(AboutDir)
	root.AddChild(DataDir)
	root.AddChild(ActionDir)

	return vfs
}
