// This is a sample filesystem that serves a couple "utilities"
// There's /time, which when read, will return a human-readable
// string of the current time.
// There's also /random, which is a file of infinite-length
// containing random bytes.
// Finally, there's /events, which records all of the high-level
// callbacks invoked on the Server struct.
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/knusbaum/go9p"
	"github.com/knusbaum/go9p/fs"
	"github.com/knusbaum/go9p/proto"
)

func addEvent(s string) {
	// f.Lock()
	// defer f.Unlock()
	// f.Data = append(f.Data, []byte(s+"\n")...)
}

func WrapEvents(f fs.File) fs.File {
	fname := f.Stat().Name
	return &fs.WrappedFile{
		File: f,
		OpenF: func(fid uint64, omode proto.Mode) error {
			addEvent(fmt.Sprintf("Open %s: mode: %d", fname, omode))
			return f.Open(fid, omode)
		},
		ReadF: func(fid uint64, offset uint64, count uint64) ([]byte, error) {
			addEvent(fmt.Sprintf("Read %s: offset %d, count %d", fname, offset, count))
			return f.Read(fid, offset, count)
		},
		WriteF: func(fid uint64, offset uint64, data []byte) (uint32, error) {
			addEvent(fmt.Sprintf("Write %s: offset %d, data %d bytes", fname, offset, len(data)))
			return f.Write(fid, offset, data)
		},
		CloseF: func(fid uint64) error {
			addEvent(fmt.Sprintf("Close %s", fname))
			return f.Close(fid)
		},
	}
}

type VaultLink string
type MetaData struct {
	Name     string
	Outlinks []VaultLink
	Inlinks  []VaultLink
	Tags     []string
	Aliases  []string
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
func LinksFile(links []VaultLink) func() []byte {
	ls := make([]string, len(links))
	for i, l := range links {
		ls[i] = string(l)
	}
	return StringsFile(ls)
}

var User = "glenda"

func metaDataDir(md MetaData, node *fs.FS) fs.Dir {
	d := fs.NewStaticDir(node.NewStat(md.Name, User, User, 0555))
	outlinks := fs.NewDynamicFile(node.NewStat("outlinks", User, User, 0444), LinksFile(md.Outlinks))
	inlinks := fs.NewDynamicFile(node.NewStat("inlinks", User, User, 0444), LinksFile(md.Inlinks))
	tags := fs.NewDynamicFile(node.NewStat("tags", User, User, 0444), StringsFile(md.Tags))
	aliases := fs.NewDynamicFile(node.NewStat("aliases", User, User, 0444), StringsFile(md.Aliases))
	d.AddChild(outlinks)
	d.AddChild(inlinks)
	d.AddChild(tags)
	d.AddChild(aliases)

	return d
}

func main() {
	utilFS, root := fs.NewFS("glenda", "glenda", 0777)
	events := fs.NewStaticFile(utilFS.NewStat("events", User, User, 0444), []byte{})
	root.AddChild(events)
	root.AddChild(metaDataDir(MetaData{
		Name:     "Note.md",
		Outlinks: []VaultLink{"out.md"},
		Inlinks:  []VaultLink{"in.md"},
		Tags:     []string{},
		Aliases:  []string{"Paper", "Sticky"},
	}, utilFS))

	root.AddChild(
		WrapEvents(fs.NewDynamicFile(utilFS.NewStat("time", "glenda", "glenda", 0444),
			func() []byte {
				return []byte(time.Now().String() + "\n")
			},
		)),
	)
	// root.AddChild(
	// 	WrapEvents(events, &fs.WrappedFile{
	// 		File: fs.NewBaseFile(utilFS.NewStat("random", "glenda", "glenda", 0444)),
	// 		ReadF: func(fid uint64, offset uint64, count uint64) ([]byte, error) {
	// 			bs := make([]byte, count)
	// 			rand.Reader.Read(bs)
	// 			return bs, nil
	// 		},
	// 	}),
	// )
	// Post a local service.
	go9p.PostSrv("utilfs", utilFS.Server())
}
