package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/knusbaum/go9p"
	"github.com/knusbaum/go9p/proto"
)

var _ go9p.Conn = &conn{}
var _ go9p.Srv = &DirFS{}

type OSPath struct {
	path  string
	atime time.Time
	qid   proto.Qid
}

func (op OSPath) ATime() uint32 {
	return uint32(op.atime.Unix())
}

func (op OSPath) Stat() (proto.Stat, error) {
	s, err := os.Stat(op.path)
	if err != nil {
		return proto.Stat{}, err
	}

	return proto.Stat{
		Type:   0,
		Dev:    0,
		Qid:    op.qid,
		Mode:   uint32(s.Mode()),
		Atime:  op.ATime(),
		Mtime:  uint32(s.ModTime().Unix()),
		Length: uint64(s.Size()),
		Name:   path.Base(op.path),
		Uid:    "",
		Gid:    "",
		Muid:   "",
	}, nil
}

type DirFS struct {
	root OSPath

	conCount uint32
	authFunc interface{}
}

// Clunk implements go9p.Srv.
func (s *DirFS) Clunk(go9p.Conn, *proto.TClunk) (proto.FCall, error) {
	panic("unimplemented")
}

// Create implements go9p.Srv.
func (s *DirFS) Create(go9p.Conn, *proto.TCreate) (proto.FCall, error) {
	panic("unimplemented")
}

// Open implements go9p.Srv.
func (s *DirFS) Open(go9p.Conn, *proto.TOpen) (proto.FCall, error) {
	panic("unimplemented")
}

// Read implements go9p.Srv.
func (s *DirFS) Read(go9p.Conn, *proto.TRead) (proto.FCall, error) {
	panic("unimplemented")

}

// Remove implements go9p.Srv.
func (s *DirFS) Remove(go9p.Conn, *proto.TRemove) (proto.FCall, error) {
	panic("unimplemented")
}

// Stat implements go9p.Srv.
func (s *DirFS) Stat(go9p.Conn, *proto.TStat) (proto.FCall, error) {
	panic("unimplemented")
}

// Walk implements go9p.Srv.
func (s *DirFS) Walk(go9p.Conn, *proto.TWalk) (proto.FCall, error) {
	panic("unimplemented")
}

// Write implements go9p.Srv.
func (s *DirFS) Write(go9p.Conn, *proto.TWrite) (proto.FCall, error) {
	panic("unimplemented")
}

// Wstat implements go9p.Srv.
func (s *DirFS) Wstat(go9p.Conn, *proto.TWstat) (proto.FCall, error) {
	panic("unimplemented")
}

func (s *DirFS) NewConn() go9p.Conn {
	s.conCount += 1
	return &conn{connID: s.conCount}
}
func (_ *DirFS) Version(gc go9p.Conn, t *proto.TRVersion) (proto.FCall, error) {
	var reply proto.TRVersion
	if t.Type == proto.Tversion && t.Version == "9P2000" {
		if t.Msize > proto.MaxMsgLen {
			t.Msize = proto.MaxMsgLen
		}
		gc.(*conn).msize = t.Msize
		reply = *t
		reply.Type = proto.Rversion
		return &reply, nil
	} else {
		return nil, fmt.Errorf("Cannot reply to type %d\n", t.Type)
	}
}
func (s *DirFS) Auth(gc go9p.Conn, t *proto.TAuth) (proto.FCall, error) {
	panic("Unimplemented")
}
func (s *DirFS) Attach(gc go9p.Conn, t *proto.TAttach) (proto.FCall, error) {
	c := gc.(*conn)

	if s.authFunc == nil {
		log.Printf("%s attached", t.Uname)
		c.fids.Store(t.Fid, newFidInfo(t.Uname, s.root))
		stat, err := s.root.Stat()
		return &proto.RAttach{Header: proto.Header{
			Type: proto.Rattach,
			Tag:  t.Tag,
		}, Qid: stat.Qid}, err
	}

	log.Printf("Loading info from C: %p, t.Afid: %d\n", c, t.Afid)
	i, ok := c.fids.Load(t.Afid)
	if !ok {
		return &proto.RError{
			Header: proto.Header{Type: proto.Rerror, Tag: t.Tag},
			Ename:  "Not Authenticated.",
		}, nil
	}
	info := i.(*fidInfo)
	if err, ok := info.extra.(error); ok {
		return &proto.RError{
			Header: proto.Header{
				Type: proto.Rerror,
				Tag:  t.Tag,
			},
			Ename: err.Error(),
		}, nil
	}

	authName, ok := info.extra.(string)
	if !ok {
		return &proto.RError{
			Header: proto.Header{
				Type: proto.Rerror,
				Tag:  t.Tag,
			},
			Ename: "Not Authenticated."}, nil
	}
	// TODO: For some reason, these don't seem to need to match.
	// User is authenticated as ai.Cuid, *not* necessarily as t.Uname.
	//	if t.Uname != ai.Cuid {
	//		return &proto.RError{proto.Header{t.Type, t.Tag}, "Bad attach uname"}, nil
	//	}
	stat, err := s.root.Stat()
	if err != nil {
		return nil, err
	}
	c.fids.Store(t.Fid, newFidInfo(authName, s.root))
	return &proto.RAttach{
		Header: proto.Header{
			Type: proto.Rattach,
			Tag:  t.Tag,
		}, Qid: stat.Qid,
	}, nil
}

type fidInfo struct {
	path       OSPath
	openMode   proto.Mode
	openOffset uint64
	uname      string // uname inherited during walk.
	extra      interface{}
}

func newFidInfo(uname string, path OSPath) *fidInfo {
	return &fidInfo{
		path:     path,
		openMode: proto.None,
		uname:    uname,
	}
}

type conn struct {
	connID uint32
	fids   sync.Map
	tags   sync.Map
	msize  uint32
}

type ctxCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *conn) toConnFid(fid uint32) uint64 {
	bcid := uint64(c.connID)
	bcid = bcid << 32
	bcid = bcid | uint64(fid)
	return bcid
}

func (c *conn) TagContext(tag uint16) context.Context {
	v, ok := c.tags.Load(tag)
	if !ok {
		ctx, cancel := context.WithCancel(context.Background())
		ctxc := &ctxCancel{ctx, cancel}
		c.tags.Store(tag, ctxc)
		return ctx
	}
	ctxc := v.(*ctxCancel)
	return ctxc.ctx
}

func (c *conn) DropContext(tag uint16) {
	v, ok := c.tags.Load(tag)
	if !ok {
		return
	}
	ctxc := v.(*ctxCancel)
	c.tags.Delete(tag)
	ctxc.cancel()
}
