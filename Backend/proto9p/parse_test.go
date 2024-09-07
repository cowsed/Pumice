package proto9p

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	fuzz "github.com/AdaLogics/go-fuzz-headers"
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
		fz := fuzz.NewConsumer(extra)
		var targetSlice []Qid
		err := fz.CreateSlice(&targetSlice)
		if err != nil {
			t.Fatalf("Test Setup Error creating fuzz slice: %v", err)
		}
		RoundTripFCall(t, &RWalk{
			Tag:   Tag(a),
			Wqids: targetSlice,
		})

	})
}
