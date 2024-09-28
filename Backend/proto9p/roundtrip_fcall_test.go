package proto9p

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"
	"time"
)

func RoundTripFCall(t *testing.T, fc FCall) {
	var tw = TypedWriter{&bytes.Buffer{}}
	var err error

	err = fc.writeTo(tw)
	if err != nil {
		t.Fatalf("Failed to write %T (%+v): %v", fc, fc, err)
	}

	bs := tw.Finish()
	var tr = TypedReader{bytes.NewBuffer(bs)}
	fc2, err := ParseFCall(tr)
	if err != nil {
		t.Fatalf("Failed to parse %T (%+v): %v", fc, fc, err)
	}

	if fc.Type() != fc2.Type() {
		t.Fatalf("Wrong type of fcall. Wanted %v, got %v", fc.Type(), fc2.Type())
	}

	if !reflect.DeepEqual(fc, fc2) {
		t.Fatalf("Value did not survive round trip. Wanted %v of type %T, got %v of type %T", fc, fc, fc2, fc2)
	}

}

func FuzzRoundTripRError(f *testing.F) {
	f.Add(uint16(1), "")
	f.Add(uint16(1), "broken")
	f.Add(uint16(0), strings.Repeat("uh oh ", 60))
	f.Fuzz(func(t *testing.T, tag uint16, ename string) {
		RoundTripFCall(t, &RError{
			Tag:   Tag(tag),
			Ename: ename,
		})

	})
}

func FuzzRoundTripTCreate(f *testing.F) {
	f.Add(uint16(1), uint32(0), "", uint32(1), uint8(0))
	f.Fuzz(func(t *testing.T, tag uint16, fid uint32, name string, perm uint32, mode uint8) {
		RoundTripFCall(t, &TCreate{
			Tag:  Tag(tag),
			Fid:  Fid(fid),
			Name: name,
			perm: Perm(perm),
			Mode: Mode(mode),
		})

	})
}

func FuzzRoundTripTClunk(f *testing.F) {
	f.Add(uint16(1), uint32(0))
	f.Fuzz(func(t *testing.T, tag uint16, fid uint32) {
		RoundTripFCall(t, &TClunk{
			Tag: Tag(tag),
			Fid: Fid(fid),
		})

	})
}

func FuzzRoundTripTRemove(f *testing.F) {
	f.Add(uint16(1), uint32(0))
	f.Fuzz(func(t *testing.T, tag uint16, fid uint32) {
		RoundTripFCall(t, &TRemove{
			Tag: Tag(tag),
			Fid: Fid(fid),
		})

	})
}

func FuzzRoundTripTStat(f *testing.F) {
	f.Add(uint16(1), uint32(1))
	f.Fuzz(func(t *testing.T, tag uint16, fid uint32) {
		RoundTripFCall(t, &TStat{
			Tag: Tag(tag),
			Fid: Fid(fid),
		})

	})
}

func FuzzRoundTripRStat(f *testing.F) {
	f.Add(uint16(1), uint16(2), uint32(3), uint8(4), uint32(5), uint64(6), uint8(7), uint32(8), uint32(9), uint64(10), "file.txt", "name", "group", "name")
	f.Fuzz(func(t *testing.T, tag uint16, type_ uint16, dev uint32, qidT uint8, qidV uint32, qidU uint64, mode uint8, atime uint32, mtime uint32, length uint64, name string, uid string, gid string, muid string) {
		RoundTripFCall(t, &RStat{
			Tag: Tag(tag),
			Stat: Stat{
				type_: type_,
				dev:   dev,
				qid: Qid{
					Qtype: QType(qidT),
					Vers:  qidV,
					Uid:   qidU,
				},
				mode:   Mode(mode),
				atime:  time.Unix(int64(atime), 0),
				mtime:  time.Unix(int64(mtime), 0),
				length: length,
				name:   name,
				uid:    uid,
				gid:    gid,
				muid:   muid,
			},
		})

	})
}

func FuzzRoundTripWTStat(f *testing.F) {
	f.Add(uint16(1), uint32(1), uint16(1), uint32(1), uint8(1), uint32(1), uint64(1), uint8(1), uint32(1), uint32(1), uint64(1), "file.txt", "name", "group", "name")
	f.Fuzz(func(t *testing.T, tag uint16, fid uint32, type_ uint16, dev uint32, qidT uint8, qidV uint32, qidU uint64, mode uint8, atime uint32, mtime uint32, length uint64, name string, uid string, gid string, muid string) {
		RoundTripFCall(t, &TWStat{
			Tag: Tag(tag),
			Fid: Fid(fid),
			Stat: Stat{
				type_: type_,
				dev:   dev,
				qid: Qid{
					Qtype: QType(qidT),
					Vers:  qidV,
					Uid:   qidU,
				},
				mode:   Mode(mode),
				atime:  time.Unix(int64(atime), 0),
				mtime:  time.Unix(int64(mtime), 0),
				length: length,
				name:   name,
				uid:    uid,
				gid:    gid,
				muid:   muid,
			},
		})

	})
}

func FuzzRoundTripRWStat(f *testing.F) {
	f.Add(uint16(1))
	f.Fuzz(func(t *testing.T, tag uint16) {
		RoundTripFCall(t, &RWStat{
			Tag: Tag(tag),
		})

	})
}

func FuzzRoundTripRClunk(f *testing.F) {
	f.Add(uint16(1))
	f.Fuzz(func(t *testing.T, tag uint16) {
		RoundTripFCall(t, &TClunk{
			Tag: Tag(tag),
		})

	})
}
func FuzzRoundTripRRemove(f *testing.F) {
	f.Add(uint16(1))
	f.Fuzz(func(t *testing.T, tag uint16) {
		RoundTripFCall(t, &RRemove{
			Tag: Tag(tag),
		})

	})
}

func FuzzRoundTripRCreate(f *testing.F) {
	f.Add(uint16(1), uint8(1), uint32(0), uint64(0), uint32(0))
	f.Fuzz(func(t *testing.T, tag uint16, qT uint8, qV uint32, qU uint64, iou uint32) {
		RoundTripFCall(t, &RCreate{
			Tag: Tag(tag),
			Qid: Qid{
				Qtype: QType(qT),
				Vers:  qV,
				Uid:   qU,
			},
			IOUnit: iou,
		})

	})
}

func FuzzRoundTripTopen(f *testing.F) {
	f.Add(uint16(1), uint32(0), uint8(0))
	f.Fuzz(func(t *testing.T, tag uint16, fid uint32, mode uint8) {
		RoundTripFCall(t, &TOpen{
			Tag:  Tag(tag),
			Fid:  Fid(fid),
			Mode: Mode(mode),
		})

	})
}
func FuzzRoundTripRopen(f *testing.F) {
	f.Add(uint16(1), uint8(1), uint32(0), uint64(0), uint32(0))
	f.Fuzz(func(t *testing.T, tag uint16, qT uint8, qV uint32, qU uint64, iou uint32) {
		RoundTripFCall(t, &ROpen{
			Tag: Tag(tag),
			Qid: Qid{
				Qtype: QType(qT),
				Vers:  qV,
				Uid:   qU,
			},
			IOUnit: iou,
		})

	})
}

func FuzzRoundTripTVersion(f *testing.F) {
	f.Add(uint16(1), uint32(2), "")
	f.Add(uint16(1), uint32(65535), "9p2000")
	f.Add(uint16(0), uint32(65536), "")
	f.Fuzz(func(t *testing.T, a uint16, b uint32, c string) {
		RoundTripFCall(t, &TVersion{
			Tag:     Tag(a),
			MSize:   b,
			Version: c,
		})

	})
}
func FuzzRoundTripRVersion(f *testing.F) {
	f.Add(uint16(1), uint32(2), "")
	f.Add(uint16(0), uint32(65536), "")
	f.Fuzz(func(t *testing.T, a uint16, b uint32, c string) {
		RoundTripFCall(t, &RVersion{
			Tag:     Tag(a),
			Msize:   b,
			Version: c,
		})
	})
}

func FuzzRoundTripTFlush(f *testing.F) {
	f.Add(uint16(1), uint16(1))
	f.Add(uint16(99), uint16(4))
	f.Fuzz(func(t *testing.T, a uint16, b uint16) {
		RoundTripFCall(t, &TFlush{
			Tag:    Tag(a),
			Oldtag: Tag(b),
		})
	})
}
func FuzzRoundTripRFlush(f *testing.F) {
	f.Add(uint16(1))
	f.Fuzz(func(t *testing.T, a uint16) {
		RoundTripFCall(t, &RFlush{
			Tag: Tag(a),
		})

	})
}

func FuzzRoundTripTRead(f *testing.F) {
	f.Add(uint16(1), uint32(1), uint64(1), uint32(1))
	f.Fuzz(func(t *testing.T, a uint16, b uint32, c uint64, d uint32) {
		RoundTripFCall(t, &TRead{
			Tag:    Tag(a),
			Fid:    Fid(b),
			Offset: c,
			Count:  d,
		})

	})
}
func FuzzRoundTripRRead(f *testing.F) {
	f.Add(uint16(1), []byte{2, 4})
	f.Fuzz(func(t *testing.T, a uint16, b []byte) {
		RoundTripFCall(t, &RRead{
			Tag:  Tag(a),
			Data: []byte{},
		})

	})
}
func FuzzRoundTripTWrite(f *testing.F) {
	f.Add(uint16(1), uint32(1), uint64(1), []byte{2, 4})
	f.Fuzz(func(t *testing.T, a uint16, b uint32, c uint64, d []byte) {
		RoundTripFCall(t, &TWrite{
			Tag:    Tag(a),
			Fid:    Fid(b),
			Offset: c,
			Data:   d,
		})

	})
}
func FuzzRoundTripRWrite(f *testing.F) {
	f.Add(uint16(1), uint32(1))
	f.Fuzz(func(t *testing.T, a uint16, b uint32) {
		RoundTripFCall(t, &RWrite{
			Tag:   Tag(a),
			Count: b,
		})

	})
}

func FuzzRoundTripTWalk(f *testing.F) {
	f.Add(uint16(1), uint32(1), uint32(1), "abababababababababa", "b")

	f.Fuzz(func(t *testing.T, a uint16, b uint32, c uint32, extra string, split string) {
		RoundTripFCall(t, &TWalk{
			Tag:    Tag(a),
			Fid:    Fid(b),
			NewFid: Fid(c),
			WNames: strings.Split(extra, split),
		})

	})
}
func FuzzRoundTripRWalk(f *testing.F) {
	f.Add(uint16(1), bytes.Repeat([]byte{1, 2, 3, 4, 5}, 50))
	f.Fuzz(func(t *testing.T, a uint16, extra []byte) {
		qids := make([]Qid, len(extra)/3)
		for i := 0; i < len(extra)/(13); i++ {
			qids[i].Qtype = QType(extra[i*13])
			qids[i].Vers = binary.LittleEndian.Uint32(extra[i*13+1:])
			qids[i].Uid = binary.LittleEndian.Uint64(extra[i*13+5:])

		}

		RoundTripFCall(t, &RWalk{
			Tag:   Tag(a),
			Wqids: qids,
		})

	})
}

func FuzzRoundTripTAuth(f *testing.F) {
	f.Add(uint16(1), uint32(1), "", "")
	f.Fuzz(func(t *testing.T, tag uint16, afid uint32, uname, aname string) {
		RoundTripFCall(t, &TAuth{
			Tag:   Tag(tag),
			afid:  0,
			uname: "",
			aname: "",
		})

	})
}

func FuzzRoundTripRAuth(f *testing.F) {
	f.Add(uint16(1), uint8(1), uint32(1), uint64(1))
	f.Fuzz(func(t *testing.T, tag uint16, qidT uint8, qidVers uint32, qidUid uint64) {
		RoundTripFCall(t, &RAuth{
			Tag: Tag(tag),
			aqid: Qid{
				Qtype: QType(qidT),
				Vers:  qidVers,
				Uid:   qidUid,
			},
		})

	})
}
