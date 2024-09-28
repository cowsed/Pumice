package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cowsed/Pumice/App/config"
	"github.com/cowsed/Pumice/App/data"
	"github.com/knusbaum/go9p"
	fs9p "github.com/knusbaum/go9p/fs"
)

func doUI(vaultPath data.OSPath, cfg Config, update chan BackendUpdate, win fyne.Window) {
	vaultName := vaultPath.Base()
	label := widget.NewLabel(fmt.Sprintf("Openning %s...", vaultName))
	progress := widget.NewProgressBarInfinite()
	progress.Start()

	win.SetContent(container.NewVBox(label, progress))
	for update := range update {
		label.SetText(update.Describe())
	}
}

func handleConfigLoad(bu BackendUpdate) Config {
	cfgUpdate, ok := bu.(ConfigLoaded)
	if !ok {
		panic("Application error")
	}
	return cfgUpdate.config
}

func vaultFS(path data.OSPath) fs.FS {
	return os.DirFS(path.String())
}

func allFilesOfType(filesys fs.FS, ext string) ([]string, error) {
	mds := []string{}
	err := fs.WalkDir(filesys, ".",
		func(fpath string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if path.Ext(fpath) != ext {
				return nil
			}

			mds = append(mds, fpath)
			return nil
		})
	if err != nil {
		return mds, err
	}
	return mds, nil

}

type CacheResponse struct {
	path  string
	err   error
	cache data.NoteCache
}

func NewCacheEntryErr(path string, err error) CacheResponse {
	return CacheResponse{
		path:  path,
		err:   err,
		cache: data.NoteCache{},
	}
}

func readFiles(filesys fs.FS, in chan string, out chan CacheResponse) {
	i := 0
	for path := range in {
		i++
		// Open File
		fil, err := filesys.Open(path)
		if err != nil {
			out <- NewCacheEntryErr(path, err)
			continue
		}

		// Read file
		bs, err := io.ReadAll(fil)
		if err != nil {
			out <- NewCacheEntryErr(path, err)
			continue
		}

		// Parse File
		cache, _, err := data.MakeNoteCache(data.VaultLocation(path), bs)
		if err != nil {
			out <- NewCacheEntryErr(path, err)
			continue
		}
		log.Println(i, "Im looking at", path)
		out <- CacheResponse{
			path:  path,
			err:   nil,
			cache: cache,
		}

	}
}

func CacheAll(mds []string, filesys fs.FS) []data.NoteCache {
	num_threads := 1

	in := make(chan string, num_threads)
	out := make(chan CacheResponse)

	// Start workers
	for i := 0; i < num_threads; i++ {
		go readFiles(filesys, in, out)
	}

	//Dump in
	go func() {
		defer close(in)
		for _, path := range mds {
			in <- path
		}
	}()

	caches := []CacheResponse{}
	count := 0
	for ent := range out {
		count++
		if count == len(mds) {
			close(out)
		}
		if ent.err != nil {
			log.Printf("Err on %s, %v\n", ent.path, ent.err)
			continue
		}
		caches = append(caches, ent)
	}

	values := []data.NoteCache{}
	for _, v := range caches {
		values = append(values, v.cache)
	}

	return values
}

func main() {
	flags := parseFlags()
	slog.Info("Loaded flags", "flags", flags)

	// updates := LoadWorkspace(flags)
	// fmt.Println(updates)

	filesys := vaultFS(flags.VaultPath)
	mds, err := allFilesOfType(filesys, ".md")
	if err != nil {
		panic(err)
	}

	log.Println("There are ", len(mds), "markdown files here")

	caches := CacheAll(mds, filesys)

	log.Printf("Read %v of %v files", len(caches), len(mds))

	vaultCache := makeVaultCacheFS(caches)
	log.Println("serving")
	go9p.PostSrv("vaultfs", vaultCache.Server())

	fmt.Println(caches[0])

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

var User = "glenda"
var Group = User

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

	return dir
}

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

func makeVaultCacheFS(caches []data.NoteCache) *fs9p.FS {
	vfs, root := fs9p.NewFS(User, Group, 0755)

	AboutDir := makeAboutDir(vfs)
	DataDir := makeDataDir(caches, vfs)

	ActionDir := fs9p.NewStaticDir(vfs.NewStat("actions", User, Group, 0755))
	searchFile := fs9p.NewDynamicFile(vfs.NewStat("search", User, Group, 0444), func() []byte { return []byte("coming soon\n") })
	ActionDir.AddChild(searchFile)

	root.AddChild(AboutDir)
	root.AddChild(DataDir)
	root.AddChild(ActionDir)

	return vfs
}
