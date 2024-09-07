package fs

import (
	"PumiceBackend/proto9p"
	"PumiceBackend/server"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"sync"
)

var Version = "9P2000"

func NewServer(root string) (FS, error) {

	if _, err := os.Stat(root); errors.Is(err, os.ErrNotExist) {
		return FS{}, fmt.Errorf("couldnt open directory `%s`: %w", root, os.ErrNotExist)
	}

	return FS{
		maxSize:  65535,
		root:     root,
		qids:     map[ServedPath]proto9p.Qid{},
		qidCount: 0,
		clients:  map[string]ClientEntry{},
		Mutex:    sync.Mutex{},
	}, nil
}

type FidEntry struct {
	path string
}
type ClientEntry struct {
	attachPoint ServedPath
}
type conn struct {
	uname string
}

func (conn *conn) SetUsername(s string) {
	conn.uname = s
}
func (conn *conn) Username() string {
	return conn.uname
}

var _ server.Conn = &conn{}

type FS struct {
	maxSize  uint32
	root     string // location in server OS filesystem that this server is serving
	qids     map[ServedPath]proto9p.Qid
	qidCount uint64

	clients map[string]ClientEntry
	sync.Mutex
}

func (fs *FS) QidFor(rpath ServedPath) (proto9p.Qid, error) {
	fs.Lock()
	defer fs.Unlock()
	q, exists := fs.qids[rpath]
	if exists {
		return q, nil
	}

	q.Qtype = 0
	q.Vers = 0

	fs.qidCount++
	q.Uid = fs.qidCount
	fs.qids[rpath] = q

	return q, nil

}

type ServedPath string

func (f *FS) FileExists(rpath ServedPath) bool {
	if _, err := os.Stat(path.Join(f.root, string(rpath))); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// 9p.io/magic/man2html/5/attach
func (f *FS) Attach(conn server.Conn, p *proto9p.TAttach) (proto9p.FCall, error) {
	if p.Afid != proto9p.NOFID {
		return &proto9p.RError{
			Tag:   p.Tag,
			Ename: "this server does not support authentication",
		}, nil
	}
	if !f.FileExists(ServedPath(p.Aname)) {
		return &proto9p.RError{
			Tag:   p.Tag,
			Ename: fmt.Sprintf("requested path `%s` not found", p.Aname),
		}, nil
	}
	if _, exists := f.clients[p.Uname]; exists {
		return &proto9p.RError{Tag: p.Tag, Ename: "server does not support multiple users of same uname right now"}, nil
	}
	f.clients[p.Uname] = ClientEntry{
		attachPoint: ServedPath(p.Aname),
	}
	conn.SetUsername(p.Uname)

	q, err := f.QidFor(ServedPath(p.Aname))
	if err != nil {
		return nil, err
	}
	return &proto9p.RAttach{
		Tag: p.Tag,
		Qid: q,
	}, nil
}

// Auth implements server.Server.
func (f *FS) Auth(server.Conn, *proto9p.TAuth) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Clunk implements server.Server.
func (f *FS) Clunk(server.Conn, *proto9p.TClunk) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Create implements server.Server.
func (f *FS) Create(server.Conn, *proto9p.TCreate) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Flush implements server.Server.
func (f *FS) Flush(server.Conn, *proto9p.TFlush) (proto9p.FCall, error) {
	panic("unimplemented")
}

// NewConnection implements server.Server.
func (f *FS) NewConnection() server.Conn {
	return &conn{
		uname: "",
	}
}

// Open implements server.Server.
func (f *FS) Open(server.Conn, *proto9p.TOpen) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Read implements server.Server.
func (f *FS) Read(server.Conn, *proto9p.TRead) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Remove implements server.Server.
func (f *FS) Remove(server.Conn, *proto9p.TRemove) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Stat implements server.Server.
func (f *FS) Stat(server.Conn, *proto9p.TStat) (proto9p.FCall, error) {
	panic("unimplemented")
}

// http://9p.io/magic/man2html/5/version
func (f *FS) Version(c server.Conn, vers *proto9p.TVersion) (proto9p.FCall, error) {
	if vers.Tag != proto9p.NOTAG {
		slog.Warn("version tag should be ~0 (65535)", "tag", vers.Tag)
	}
	f.maxSize = vers.MSize
	if vers.Version != Version {
		return &proto9p.RVersion{
			Tag:     vers.Tag,
			Msize:   f.maxSize,
			Version: "unknown",
		}, nil
	}
	// TODO version resets the connection. if we have received a connection already, we should reset
	return &proto9p.RVersion{
		Tag:     vers.Tag,
		Msize:   f.maxSize,
		Version: Version,
	}, nil
}

// Walk implements server.Server.
func (f *FS) Walk(server.Conn, *proto9p.TWalk) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Write implements server.Server.
func (f *FS) Write(server.Conn, *proto9p.TWrite) (proto9p.FCall, error) {
	panic("unimplemented")
}

// Wstat implements server.Server.
func (f *FS) Wstat(server.Conn, *proto9p.TWStat) (proto9p.FCall, error) {
	panic("unimplemented")
}

var _ server.Server = &FS{}
