package proto9p

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var ErrUnknownType = errors.New("9p packet type not recognized")

func ParseFCall(r io.Reader) (FCall, error) {
	wire_reader := TypedReader{r}
	size, err := wire_reader.Read32()
	if err != nil {
		return nil, err
	}
	// size - 4 since size takes up 4 bytes
	wanted := size - 4
	bs, err := wire_reader.ReadN(int(wanted))
	if err != nil {
		return nil, err
	}
	packet_reader := TypedReader{bytes.NewReader(bs)}
	var pt8 uint8
	pt8, err = packet_reader.Read8()
	if err != nil {
		return nil, err
	}

	switch Type(pt8) {
	case Tversion:
		return (&TVersion{}).fill(packet_reader)
	case Rversion:
		return (&RVersion{}).fill(packet_reader)
	case Tflush:
		return (&TFlush{}).fill(packet_reader)
	case Rflush:
		return (&RFlush{}).fill(packet_reader)
	case Twalk:
		return (&TWalk{}).fill(packet_reader)
	case Rwalk:
		return (&RWalk{}).fill(packet_reader)
	default:
		return nil, fmt.Errorf("got %d: %w", pt8, ErrUnknownType)
	}

}

func (tv *TWalk) fill(r TypedReader) (FCall, error) {
	// size[4] Twalk tag[2] fid[4] newfid[4] nwname[2] nwname*(wname[s])	var err error
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	tv.newFid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	numWname, err := r.Read16()
	if err != nil {
		return nil, err
	}
	wnames := make([]string, numWname)
	for i := 0; i < int(numWname); i++ {
		wnames[i], err = r.ReadString()
		if err != nil {
			return nil, err
		}
	}
	return tv, nil

}
func (tv *RWalk) fill(r TypedReader) (FCall, error) {
	// size[4] Rwalk tag[2] nwqid[2] nwqid*(wqid[13])
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	numQids, err := r.Read16()
	qids := make([]Qid, numQids)
	for i := 0; i < int(numQids); i++ {
		qids[i], err = r.ReadQid()
		if err != nil {
			return nil, err
		}
	}

	return tv, err
}

func (tv *TFlush) fill(r TypedReader) (FCall, error) {
	// size[4] Tflush tag[2] oldtag[2]
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Oldtag, err = r.ReadTag()
	return tv, err

}
func (tv *RFlush) fill(r TypedReader) (FCall, error) {
	// size[4] Rflush tag[2]
	var err error
	tv.Tag, err = r.ReadTag()
	return tv, err
}

func (tv *TVersion) fill(r TypedReader) (FCall, error) {
	// size[4] Tversion tag[2] msize[4] version[s]
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	tv.msize, err = r.Read32()
	if err != nil {
		return nil, err
	}
	tv.version, err = r.ReadString()
	return tv, err
}

func (tv *RVersion) fill(r TypedReader) (FCall, error) {
	// size[4] Rversion tag[2] msize[4] version[s]
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	tv.msize, err = r.Read32()
	if err != nil {
		return nil, err
	}
	tv.version, err = r.ReadString()
	return tv, err
}
