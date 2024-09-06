package proto9p

import (
	"encoding/binary"
	"io"
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
	bs := []byte{0}
	read, err := tr.Read(bs)
	if read < 1 {
		return 0, ErrBufferTooShort
	}
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
	writ, err := tr.Read(bs)
	if writ < n {
		return []byte{}, ErrBufferTooShort
	}
	if err != nil {
		return []byte{}, err
	}
	return bs, nil
}

func (tr *TypedReader) ReadString() (string, error) {
	length, err := tr.Read16()
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	n, err := tr.Read(buf)
	if n < int(length) {
		return "", ErrBufferTooShort
	} else if err != nil {
		return "", err
	}

	return string(buf), nil
}
func (tr *TypedReader) Read64() (uint64, error) {
	b64 := make([]byte, 8)
	n, err := tr.Read(b64)
	if n < 8 {
		return 0, ErrBufferTooShort
	} else if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b64), nil
}
func (tr *TypedReader) Read32() (uint32, error) {
	b32 := make([]byte, 4)
	n, err := tr.Read(b32)
	if n < 4 {
		return 0, ErrBufferTooShort
	} else if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b32), nil
}
func (tr *TypedReader) Read16() (uint16, error) {
	b16 := make([]byte, 2)
	n, err := tr.Read(b16)
	if n < 2 {
		return 0, ErrBufferTooShort
	}
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b16), nil
}

type TypedWriter struct {
	io.Writer
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
