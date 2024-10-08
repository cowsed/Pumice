package proto9p

import (
	"errors"
	"fmt"
	"time"
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
	String() string
	// With an empty FCall (constructed from size and tag read form the wire), fillFrom in the rest of the data
	fillFrom(r TypedReader) (FCall, error)
	writeTo(w TypedWriter) error
	Type() Type
}
type Type uint8
type Tag uint16

func (f Tag) String() string {
	return fmt.Sprintf("%d", f)
}

type Fid uint32

func (f Fid) String() string {
	return fmt.Sprintf("%d", f)
}

// http://9p.io/magic/man2html/5/stat
type Stat struct {
	type_  uint16
	dev    uint32
	qid    Qid
	mode   Mode
	atime  time.Time
	mtime  time.Time
	length uint64
	name   string
	uid    string
	gid    string
	muid   string
}

func (s Stat) String() string {
	return fmt.Sprintf("Stat{type:%v dev:%v qid:%v mode:%v atime:%v mtime:%v length:%d name:%s uid:%s gid:%s muid:%s}", s.type_, s.dev, s.qid, s.mode, s.atime, s.mode, s.length, s.name, s.uid, s.gid, s.muid)
}

const NOTAG Tag = 0
const NOFID Fid = Fid(^uint32(0))

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

func (t Type) String() string {
	switch t {
	case Rattach:
	case Rauth:
		return "Rauth"
	case Rclunk:
		return "Rclunk"
	case Rcreate:
		return "Rcreate"
	case Rerror:
		return "Rerror"
	case Rflush:
		return "Rflush"
	case Ropen:
		return "Ropen"
	case Rread:
		return "Rread"
	case Rremove:
		return "Rremove"
	case Rstat:
		return "Rstat"
	case Rversion:
		return "Rversion"
	case Rwalk:
		return "Rwalk"
	case Rwrite:
		return "Rwrite"
	case Rwstat:
		return "Rwstat"
	case Tattach:
		return "Tattach"
	case Tauth:
		return "Tauth"
	case Tclunk:
		return "Tclunk"
	case Tcreate:
		return "Tcreate"
	case Terror:
		return "Terror(invalid)"
	case Tflush:
		return "Tflush"
	case Topen:
		return "Topen"
	case Tread:
		return "Tread"
	case Tremove:
		return "Tremove"
	case Tstat:
		return "Tstat"
	case Tversion:
		return "Tversion"
	case Twalk:
		return "Twalk"
	case Twrite:
		return "Twrite"
	case Twstat:
		return "Twstat"
	}
	return "UNKNOWN TYPE"
}

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
func (ts *TStat) String() string {
	return fmt.Sprintf("%s %s %s", ts.Type().String(), ts.Tag.String(), ts.Fid.String())
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
func (ts *RStat) String() string {
	return fmt.Sprintf("%s %s %s", ts.Type().String(), ts.Tag.String(), ts.Stat.String())
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
func (ts *TWStat) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v stat:%v", ts.Type(), ts.Tag, ts.Fid, ts.Stat)
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
func (ts *TClunk) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v", ts.Type(), ts.Tag, ts.Fid)
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
func (ts *TRemove) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v", ts.Type(), ts.Tag, ts.Fid)
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
	perm Perm
	Mode
}

func (ts *TCreate) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v name:%s perm:%v mode:%v", ts.Type(), ts.Tag, ts.Fid, ts.Name, ts.perm, ts.Mode)
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
func (ts *TOpen) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v mode:%v", ts.Type(), ts.Tag, ts.Fid, ts.Mode)
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
	Ename string
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
	Afid  Fid
	Uname string
	Aname string
}

func (r *TAttach) Type() Type {
	return Tattach
}
func (r *TAttach) String() string {
	return fmt.Sprintf("%v tag%v fid:%v afid:%v uname:%s aname:%s", r.Type(), r.Tag, r.Fid, r.Afid, r.Uname, r.Aname)
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
	MSize   uint32
	Version string
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
	Msize   uint32
	Version string
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
	NewFid Fid
	WNames []string
}

func (r *TWalk) Type() Type {
	return Twalk
}
func (ts *TWalk) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v newFid:%v wnames:%v", ts.Type(), ts.Tag, ts.Fid, ts.NewFid, ts.WNames)
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
func (ts *TRead) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v offset:%v count:%v", ts.Type(), ts.Tag, ts.Fid, ts.Offset, ts.Count)
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
func (ts *TWrite) String() string {
	return fmt.Sprintf("%v tag:%v fid:%v offset:%v len:%d", ts.Type(), ts.Tag, ts.Fid, ts.Offset, len(ts.Data))
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

type Perm uint32

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

type QType uint8

// https://9fans.github.io/plan9port/man/man9/intro.html
type Qid struct {
	Qtype QType
	Vers  uint32
	Uid   uint64
}
