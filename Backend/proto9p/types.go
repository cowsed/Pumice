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
}
type Type uint8
type Tag uint16
type Fid uint32

const NOTAG Tag = 0

var (
	_ FCall = &TVersion{}
	_       = &RVersion{}
	// _       = &TAuth{}
	// _       = &RAuth{}
	// _       = &TAttach{}
	// _       = &RAttach{}
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

// http://9p.io/magic/man2html/5/flush
//
// wire format:
//
//	size[4] Tflush tag[2] oldtag[2]
type TFlush struct {
	Tag
	Oldtag Tag
}

// http://9p.io/magic/man2html/5/flush
//
// wire format:
//
//	size[4] Rflush tag[2]
type RFlush struct {
	Tag
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

// http://9p.io/magic/man2html/5/walk
//
// wire format:
//
//	size[4] Rwalk tag[2] nwqid[2] nwqid*(wqid[13])
type RWalk struct {
	Tag
	Wqids []Qid
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

// http://9p.io/magic/man2html/5/read
//
// wire format:
//
//	size[4] Rread tag[2] count[4] data[count]
type RRead struct {
	Tag
	Data []byte
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

// http://9p.io/magic/man2html/5/read
//
// wire_format:
//
//	size[4] Rwrite tag[2] count[4]
type RWrite struct {
	Tag
	Count uint32
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
