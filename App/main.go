package main

import (
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/cowsed/Pumice/App/data"
	"github.com/knusbaum/go9p"
)

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
		// log.Println(i, "Im looking at", path)
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

	filesys := vaultFS(flags.VaultPath)
	vault := NewVaultFromFS(filesys)

	go9p.PostSrv("vaultfs", vault.Server())

}
