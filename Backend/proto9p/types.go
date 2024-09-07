package proto9p

import (
	"errors"
	"fmt"
)

// http://web.archive.org/web/20190304021146/https://knusbaum.com/useful/rfc9p2000

const MaxStrSize = 65535

var ErrBufferTooShort = errors.New("buffer does not contain all the bytes needed to parse this element")

func NewErrBufferTooShort(wanted, got int) error {
	return fmt.Errorf("wanted: %d, got: %d: %w", wanted, got, ErrBufferTooShort)
}

var ErrCouldNotWriteAll = errors.New("could not write entire message")
var ErrStringTooBig = errors.New("string too long to write to wire. max length 65535")

type FCall interface {
	// With an empty FCall (constructed from size and tag read form the wire), fillFrom in the rest of the data
	fillFrom(r TypedReader) (FCall, error)
	writeTo(w TypedWriter) error
	Type() Type
}
type Type uint8
type Tag uint16
type Fid uint32
type Stat uint32

const NOTAG Tag = 0

var (
	_ FCall = &TVersion{}
	_ FCall = &RVersion{}
	_ FCall = &TAuth{}
	_ FCall = &RAuth{}
	_ FCall = &TAttach{}
	_ FCall = &RAttach{}
	_ FCall = &RError{}
	_ FCall = &TWalk{}
	_ FCall = &RWalk{}
	_ FCall = &TOpen{}
	_ FCall = &ROpen{}
	_ FCall = &TCreate{}
	_ FCall = &RCreate{}
	_ FCall = &TRead{}
	_ FCall = &RRead{}
	_ FCall = &TWrite{}
	_ FCall = &RWrite{}
	_ FCall = &TClunk{}
	_ FCall = &RClunk{}
	_ FCall = &TRemove{}
	_ FCall = &RRemove{}
	_ FCall = &TStat{}
	_ FCall = &RStat{}
	_ FCall = &TWStat{}
	_ FCall = &RWStat{}
)

const (
	Tversion Type = 100
	Rversion Type = 101
	Tauth    Type = 102
	Rauth    Type = 103
	Tattach  Type = 104
	Rattach  Type = 105
	Terror   Type = 106 /* illegal */
	Rerror   Type = 107
	Tflush   Type = 108
	Rflush   Type = 109
	Twalk    Type = 110
	Rwalk    Type = 111
	Topen    Type = 112
	Ropen    Type = 113
	Tcreate  Type = 114
	Rcreate  Type = 115
	Tread    Type = 116
	Rread    Type = 117
	Twrite   Type = 118
	Rwrite   Type = 119
	Tclunk   Type = 120
	Rclunk   Type = 121
	Tremove  Type = 122
	Rremove  Type = 123
	Tstat    Type = 124
	Rstat    Type = 125
	Twstat   Type = 126
	Rwstat   Type = 127
)

// http://9p.io/magic/man2html/5/stat
//
// wire format:
//
//	size[4] Tstat tag[2] fid[4]
type TStat struct {
	Tag
	Fid
}

func (ts *TStat) Type() Type {
	return Tstat
}

// http://9p.io/magic/man2html/5/stat
//
// wire format:
//
//	size[4] Rstat tag[2] stat[n]
type RStat struct {
	Tag
	Stat
}

func (ts *RStat) Type() Type {
	return Rstat
}

// http://9p.io/magic/man2html/5/stat
//
// wire format:
//
//	size[4] Twstat tag[2] fid[4] stat[n]
type TWStat struct {
	Tag
	Fid
	Stat
}

func (ts *TWStat) Type() Type {
	return Twstat
}

// http://9p.io/magic/man2html/5/stat
//
// wire format:
//
//	size[4] Rwstat tag[2]
type RWStat struct {
	Tag
}

func (ts *RWStat) Type() Type {
	return Rwstat
}

// http://9p.io/magic/man2html/5/clunk
//
// wire format:
//
//	size[4] Tclunk tag[2] fid[4]
type TClunk struct {
	Tag
	Fid
}

func (t *TClunk) Type() Type {
	return Tclunk
}

// http://9p.io/magic/man2html/5/clunk
//
// wire format:
//
//	size[4] Tclunk tag[2] fid[4]
type RClunk struct {
	Tag
}

func (t *RClunk) Type() Type {
	return Rclunk
}

// http://9p.io/magic/man2html/5/remove
//
// wire format:
//
//	size[4] Tclunk tag[2] fid[4]
type TRemove struct {
	Tag
	Fid
}

func (t *TRemove) Type() Type {
	return Tremove
}

// http://9p.io/magic/man2html/5/remove
//
// wire format:
//
//	size[4] Tclunk tag[2] fid[4]
type RRemove struct {
	Tag
}

func (t *RRemove) Type() Type {
	return Tremove
}

// http://9p.io/magic/man2html/5/open
//
// wire format:
//
//	size[4] Topen tag[2] fid[4] mode[1]
type TCreate struct {
	Tag
	Fid
	Name string
	perm uint32
	Mode
}

func (t *TCreate) Type() Type {
	return Tcreate
}

// http://9p.io/magic/man2html/5/open
//
// wire format:
//
//	size[4] Ropen tag[2] qid[13] iounit[4]
type RCreate struct {
	Tag
	Qid
	IOUnit uint32
}

func (t *RCreate) Type() Type {
	return Rcreate
}

// http://9p.io/magic/man2html/5/open
//
// wire format:
//
//	size[4] Topen tag[2] fid[4] mode[1]
type TOpen struct {
	Tag
	Fid
	Mode
}

func (t *TOpen) Type() Type {
	return Topen
}

// http://9p.io/magic/man2html/5/open
//
// wire format:
//
//	size[4] Ropen tag[2] qid[13] iounit[4]
type ROpen struct {
	Tag
	Qid
	IOUnit uint32
}

func (t *ROpen) Type() Type {
	return Ropen
}

// http://9p.io/magic/man2html/5/0intro
//
// wire format:
//
//	size[4] Rerror tag[2] ename[s]
type RError struct {
	Tag
	ename string
}

func (r *RError) Type() Type {
	return Rerror
}

// http://9p.io/magic/man2html/5/attach
//
// wire format:
//
//	size[4] Tauth tag[2] afid[4] uname[s] aname[s]
type TAuth struct {
	Tag
	afid  Fid
	uname string
	aname string
}

func (r *TAuth) Type() Type {
	return Tauth
}

// http://9p.io/magic/man2html/5/attach
//
// wire format:
//
//	size[4] Rauth tag[2] aqid[13]
type RAuth struct {
	Tag
	aqid Qid
}

func (r *RAuth) Type() Type {
	return Rauth
}

// http://9p.io/magic/man2html/5/attach
//
// wire format:
//
//	size[4] Tattach tag[2] fid[4] afid[4] uname[s] aname[s]
type TAttach struct {
	Tag
	Fid
	afid  Fid
	uname string
	aname string
}

func (r *TAttach) Type() Type {
	return Tattach
}

// http://9p.io/magic/man2html/5/attach
//
// wire format:
//
//	size[4] Rattach tag[2] qid[13]
type RAttach struct {
	Tag
	Qid
}

func (r *RAttach) Type() Type {
	return Rattach
}

// http://9p.io/magic/man2html/5/version
//
// wire format:
//
//	size[4] Tversion tag[2] msize[4] version[s]
type TVersion struct {
	Tag
	msize   uint32
	version string
}

func (r *TVersion) Type() Type {
	return Tversion
}

// http://9p.io/magic/man2html/5/version
//
// wire format:
//
//	size[4] Rversion tag[2] msize[4] version[s]
type RVersion struct {
	Tag
	msize   uint32
	version string
}

func (r *RVersion) Type() Type {
	return Rversion
}

// http://9p.io/magic/man2html/5/flush
//
// wire format:
//
//	size[4] Tflush tag[2] oldtag[2]
type TFlush struct {
	Tag
	Oldtag Tag
}

func (r *TFlush) Type() Type {
	return Tflush
}

// http://9p.io/magic/man2html/5/flush
//
// wire format:
//
//	size[4] Rflush tag[2]
type RFlush struct {
	Tag
}

func (r *RFlush) Type() Type {
	return Rflush
}

// http://9p.io/magic/man2html/5/walk
//
// wire format:
//
//	size[4] Twalk tag[2] fid[4] newfid[4] nwname[2] nwname*(wname[s])
type TWalk struct {
	Tag
	Fid
	newFid Fid
	wnames []string
}

func (r *TWalk) Type() Type {
	return Twalk
}

// http://9p.io/magic/man2html/5/walk
//
// wire format:
//
//	size[4] Rwalk tag[2] nwqid[2] nwqid*(wqid[13])
type RWalk struct {
	Tag
	Wqids []Qid
}

func (r *RWalk) Type() Type {
	return Rwalk
}

// http://9p.io/magic/man2html/5/read
//
// wire format:
//
//	size[4] Tread tag[2] fid[4] offset[8] count[4]
type TRead struct {
	Tag
	Fid
	Offset uint64
	Count  uint32
}

func (r *TRead) Type() Type {
	return Tread
}

// http://9p.io/magic/man2html/5/read
//
// wire format:
//
//	size[4] Rread tag[2] count[4] data[count]
type RRead struct {
	Tag
	Data []byte
}

func (r *RRead) Type() Type {
	return Rread
}

// http://9p.io/magic/man2html/5/read
//
// wire format:
//
//	size[4] Twrite tag[2] fid[4] offset[8] count[4] data[count]
type TWrite struct {
	Tag
	Fid
	Offset uint64
	Data   []byte
}

func (r *TWrite) Type() Type {
	return Twrite
}

// http://9p.io/magic/man2html/5/read
//
// wire_format:
//
//	size[4] Rwrite tag[2] count[4]
type RWrite struct {
	Tag
	Count uint32
}

func (r *RWrite) Type() Type {
	return Rwrite
}

type Mode uint8

const (
	Oread   Mode = 0
	Owrite  Mode = 1
	Ordwr   Mode = 2
	Oexec   Mode = 3
	None    Mode = 4
	Otrunc  Mode = 0x10
	Orclose Mode = 0x40
)

type Qid struct {
	Qtype uint8
	Vers  uint32
	Uid   uint64
}
