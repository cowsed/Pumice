package proto9p

import (
	"bytes"
	"strings"
	"testing"
)

func RoundTrip[T comparable](t *testing.T, thing T, enc func(TypedWriter, T) error, dec func(TypedReader) (T, error)) {
	var b1 bytes.Buffer
	var tw = TypedWriter{&b1}
	var tr = TypedReader{&b1}
	var err error
	err = enc(tw, thing)
	if err != nil {
		t.Fatalf("Error encoding %T (%+v): %v", thing, thing, err)
	}
	thing2, err := dec(tr)
	if err != nil {
		t.Fatalf("Error encoding %T (%+v): %v", thing, thing, err)
	}
	if thing != thing2 {
		t.Fatalf("Mismatch round tripping %T: Wanted %+v, got %+v", thing, thing, thing2)
	}
}
func FuzzStringEncoding(f *testing.F) {
	f.Add("hello")
	f.Add("你好世界")
	f.Add("")
	f.Add(strings.Repeat("a", MaxStrSize-1))
	f.Add(strings.Repeat("a", MaxStrSize))
	f.Add(strings.Repeat("a", MaxStrSize+1))

	f.Fuzz(func(t *testing.T, inS string) {
		if len([]byte(inS)) > MaxStrSize {
			// Don't test these here.
			return
		}
		RoundTrip(t, inS,
			func(tw TypedWriter, s string) error { return tw.WriteString(s) },
			func(tr TypedReader) (string, error) { return tr.ReadString() },
		)

	})
}

func FuzzQidEncoding(f *testing.F) {
	f.Add(uint8(1), uint32(2), uint64(5))
	f.Add(uint8(0), uint32(0), uint64(0))

	f.Fuzz(func(t *testing.T, qtype uint8, vers uint32, uid uint64) {
		var qin = Qid{
			Qtype: QType(qtype),
			Vers:  vers,
			Uid:   uid,
		}
		RoundTrip(t, qin,
			func(tw TypedWriter, q Qid) error { return tw.WriteQid(q) },
			func(tr TypedReader) (Qid, error) { return tr.ReadQid() },
		)
	})
}
