package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cowsed/Pumice/App/data"
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

type CacheEntry struct {
	path  string
	err   error
	cache data.NoteCache
}

func NewCacheEntryErr(path string, err error) CacheEntry {
	return CacheEntry{
		path:  path,
		err:   err,
		cache: data.NoteCache{},
	}
}

func readFiles(filesys fs.FS, in chan string, out chan CacheEntry) {
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
		cache, _, err := data.MakeNoteCache(bs)
		if err != nil {
			out <- NewCacheEntryErr(path, err)
			continue
		}
		// log.Println(i, "Im looking at", path)
		out <- CacheEntry{
			path:  path,
			err:   nil,
			cache: cache,
		}

	}
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

	num_threads := 1

	in := make(chan string, num_threads)
	out := make(chan CacheEntry)

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

	caches := []CacheEntry{}
	count := 0
	for ent := range out {
		count++
		if count == len(mds)-1 {
			close(out)
		}
		if ent.err != nil {
			log.Printf("Err on %s, %v\n", ent.path, ent.err)
			continue
		}
		caches = append(caches, ent)
	}

	log.Printf("Read %v of %v files", len(caches), len(mds))

	// myApp := app.New()

	// var update BackendUpdate = <-updates

	// initialCfg := handleConfigLoad(update)
	// win := myApp.NewWindow(fmt.Sprintf("Pumice - %s", path.Base(string(flags.VaultPath))))

	// go doUI(flags.VaultPath, initialCfg, updates, win)

	// win.ShowAndRun()

}
