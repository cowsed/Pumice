package proto9p

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"
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

	if !reflect.DeepEqual(fc, fc2) {
		t.Fatalf("Value did not survive round trip. Wanted %v of type %T, got %v of type %T", fc, fc, fc2, fc2)
	}

}

func FuzzRoundTripTVersion(f *testing.F) {
	f.Add(uint16(1), uint32(2), "")
	f.Add(uint16(1), uint32(65535), "9p2000")
	f.Add(uint16(0), uint32(65536), "")
	f.Fuzz(func(t *testing.T, a uint16, b uint32, c string) {
		RoundTripFCall(t, &TVersion{
			Tag:     Tag(a),
			msize:   b,
			version: c,
		})

	})
}
func FuzzRoundTripRVersion(f *testing.F) {
	f.Add(uint16(1), uint32(2), "")
	f.Add(uint16(0), uint32(65536), "")
	f.Fuzz(func(t *testing.T, a uint16, b uint32, c string) {
		RoundTripFCall(t, &RVersion{
			Tag:     Tag(a),
			msize:   b,
			version: c,
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
			newFid: Fid(c),
			wnames: strings.Split(extra, split),
		})

	})
}
func FuzzRoundTripRWalk(f *testing.F) {
	f.Add(uint16(1), bytes.Repeat([]byte{1, 2, 3, 4, 5}, 50))
	f.Fuzz(func(t *testing.T, a uint16, extra []byte) {
		qids := make([]Qid, len(extra)/3)
		for i := 0; i < len(extra)/(13); i++ {
			qids[i].Qtype = extra[i*13]
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
				Qtype: qidT,
				Vers:  qidVers,
				Uid:   qidUid,
			},
		})

	})
}
