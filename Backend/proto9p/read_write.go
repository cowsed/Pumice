package proto9p

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type TypedReader struct {
	io.Reader
}

func (tr *TypedReader) ReadQid() (Qid, error) {
	var q Qid
	var err error
	q.Qtype, err = tr.Read8()
	if err != nil {
		return q, err
	}
	q.Vers, err = tr.Read32()
	if err != nil {
		return q, err
	}
	q.Uid, err = tr.Read64()
	if err != nil {
		return q, err
	}

	return q, nil
}
func (tr *TypedReader) ReadFid() (Fid, error) {
	u32, err := tr.Read32()
	return Fid(u32), err
}

func (tr *TypedReader) Read8() (uint8, error) {
	bs, err := tr.ReadN(1)
	if err != nil {
		return 0, err
	}

	return bs[0], nil
}

func (tr *TypedReader) ReadTag() (Tag, error) {
	u16, err := tr.Read16()
	if err != nil {
		return 0, err
	}
	return Tag(u16), nil

}

func (tr *TypedReader) ReadN(n int) ([]byte, error) {
	bs := make([]byte, n)
	read, err := tr.Read(bs)
	if read < n {
		return []byte{}, NewErrBufferTooShort(n, read)
	}
	if err != nil && !errors.Is(err, io.EOF) {
		return []byte{}, err
	}
	return bs, nil
}

func (tr *TypedReader) ReadString() (string, error) {
	length, err := tr.Read16()
	if err != nil {
		return "", err
	}

	buf, err := tr.ReadN(int(length))
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
func (tr *TypedReader) Read64() (uint64, error) {
	b64, err := tr.ReadN(8)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(b64), nil
}
func (tr *TypedReader) Read32() (uint32, error) {
	b32, err := tr.ReadN(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b32), nil
}
func (tr *TypedReader) Read16() (uint16, error) {
	b16, err := tr.ReadN(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b16), nil
}

func (tr *TypedReader) ReadStat() (Stat, error) {
	var st Stat
	var err error

	// size, but we don't actually need it at this point idt
	_, err = tr.Read16()
	if err != nil {
		return st, err
	}

	st.type_, err = tr.Read16()
	if err != nil {
		return st, err
	}
	st.dev, err = tr.Read32()
	if err != nil {
		return st, err
	}
	st.qid, err = tr.ReadQid()
	if err != nil {
		return st, err
	}
	m32, err := tr.Read32()
	st.mode = Mode(m32)
	if err != nil {
		return st, err
	}
	atime32, err := tr.Read32()
	if err != nil {
		return st, err
	}
	st.atime = time.Unix(int64(atime32), 0)
	mtime32, err := tr.Read32()
	if err != nil {
		return st, err
	}
	st.mtime = time.Unix(int64(mtime32), 0)

	st.length, err = tr.Read64()
	if err != nil {
		return st, err
	}
	st.name, err = tr.ReadString()
	if err != nil {
		return st, err
	}
	st.uid, err = tr.ReadString()
	if err != nil {
		return st, err
	}
	st.gid, err = tr.ReadString()
	if err != nil {
		return st, err
	}
	st.muid, err = tr.ReadString()
	if err != nil {
		return st, err
	}

	return st, nil
}

func (tw *TypedWriter) WriteType(t Type) error {
	return tw.Write8(uint8(t))
}

type TypedWriter struct {
	io.ReadWriter
}

func (tw *TypedWriter) Finish() []byte {
	bs, _ := io.ReadAll(tw)
	full := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(full, uint32(len(bs)+4))
	full = append(full, bs...)
	return full
}

func (tw *TypedWriter) WriteQid(q Qid) error {
	// type[1] version[4] path[8]
	err := tw.Write8(q.Qtype)
	if err != nil {
		return err
	}
	err = tw.Write32(q.Vers)
	if err != nil {
		return err
	}
	err = tw.Write64(q.Uid)
	if err != nil {
		return err
	}
	return nil
}

func (tw *TypedWriter) WriteFid(f Fid) error {
	return tw.Write32(uint32(f))
}

func (tw *TypedWriter) WriteN(bs []byte) error {
	n, err := tw.Write(bs)
	if n < len(bs) {
		return ErrCouldNotWriteAll
	}
	if err != nil {
		return err
	}
	return nil
}

func (tw *TypedWriter) Write8(v uint8) error {
	n, err := tw.Write([]byte{v})
	if n < 1 {
		return ErrCouldNotWriteAll
	}
	if err != nil {
		return err
	}
	return nil
}

func (tw *TypedWriter) Write16(v uint16) error {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	n, err := tw.Write(b)
	if n < 2 {
		return ErrCouldNotWriteAll
	}
	if err != nil {
		return err
	}
	return nil
}
func (tw *TypedWriter) Write32(v uint32) error {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	n, err := tw.Write(b)
	if n < 4 {
		return ErrCouldNotWriteAll
	}
	if err != nil {
		return err
	}
	return nil
}

func (tw *TypedWriter) Write64(v uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	n, err := tw.Write(b)
	if n < 8 {
		return ErrCouldNotWriteAll
	}
	if err != nil {
		return err
	}
	return nil
}

func (tw *TypedWriter) WriteString(s string) error {
	bs := []byte(s)
	l := len(bs)
	if l > MaxStrSize {
		return ErrStringTooBig
	}
	err := tw.Write16(uint16(l))
	if err != nil {
		return err
	}
	n, err := tw.Write(bs)
	if n < l {
		return ErrCouldNotWriteAll
	}
	if err != nil {
		return err
	}
	return nil

}

func (tw *TypedWriter) WriteTag(tag Tag) error {
	return tw.Write16(uint16(tag))
}
func (outer *TypedWriter) WriteStat(stat Stat) error {
	tw := TypedWriter{&bytes.Buffer{}}
	var err error

	err = tw.Write16(stat.type_)
	if err != nil {
		return err
	}

	err = tw.Write32(stat.dev)
	if err != nil {
		return err
	}
	err = tw.WriteQid(stat.qid)
	if err != nil {
		return err
	}
	err = tw.Write32(uint32(stat.mode))
	if err != nil {
		return err
	}
	err = tw.Write32(uint32(stat.atime.Unix()))
	if err != nil {
		return err
	}
	err = tw.Write32(uint32(stat.mtime.Unix()))
	if err != nil {
		return err
	}
	err = tw.Write64(stat.length)
	if err != nil {
		return err
	}
	err = tw.WriteString(stat.name)
	if err != nil {
		return err
	}
	err = tw.WriteString(stat.uid)
	if err != nil {
		return err
	}
	err = tw.WriteString(stat.gid)
	if err != nil {
		return err
	}
	err = tw.WriteString(stat.muid)
	if err != nil {
		return err
	}

	bs, err := io.ReadAll(tw)
	if err != nil {
		return err
	}
	err = outer.Write16(uint16(len(bs) + 4))
	if err != nil {
		return err
	}
	err = outer.WriteN(bs)
	if err != nil {
		return err
	}
	return nil
}
