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
		return (&TVersion{}).fillFrom(packet_reader)
	case Rversion:
		return (&RVersion{}).fillFrom(packet_reader)
	case Tflush:
		return (&TFlush{}).fillFrom(packet_reader)
	case Rflush:
		return (&RFlush{}).fillFrom(packet_reader)
	case Twalk:
		return (&TWalk{}).fillFrom(packet_reader)
	case Rwalk:
		return (&RWalk{}).fillFrom(packet_reader)
	case Tread:
		return (&TRead{}).fillFrom(packet_reader)
	case Rread:
		return (&RRead{}).fillFrom(packet_reader)
	case Twrite:
		return (&TWrite{}).fillFrom(packet_reader)
	case Rwrite:
		return (&RWrite{}).fillFrom(packet_reader)
	default:
		return nil, fmt.Errorf("got %d: %w", pt8, ErrUnknownType)
	}

}

// size[4] Tread tag[2] fid[4] offset[8] count[4]
func (tv *TRead) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	tv.Offset, err = r.Read64()
	if err != nil {
		return nil, err
	}
	tv.Count, err = r.Read32()
	if err != nil {
		return nil, err
	}
	return tv, nil

}

// size[4] Rread tag[2] count[4] data[count]
func (tv *RRead) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	count, err := r.Read32()
	if err != nil {
		return nil, err
	}

	tv.Data, err = r.ReadN(int(count))
	if err != nil {
		return nil, err
	}
	return tv, nil
}

// size[4] Twrite tag[2] fid[4] offset[8] count[4] data[count]
func (tv *TWrite) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	tv.Offset, err = r.Read64()
	if err != nil {
		return nil, err
	}
	count, err := r.Read32()
	if err != nil {
		return nil, err
	}
	tv.Data, err = r.ReadN(int(count))
	if err != nil {
		return nil, err
	}
	return tv, nil
}

// size[4] Rwrite tag[2] count[4]
func (tv *RWrite) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Count, err = r.Read32()
	if err != nil {
		return nil, err
	}
	return tv, nil
}

func (tv *TWalk) fillFrom(r TypedReader) (FCall, error) {
	// size[4] Twalk tag[2] fid[4] newfid[4] nwname[2] nwname*(wname[s])
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
	tv.wnames = make([]string, numWname)
	for i := 0; i < int(numWname); i++ {
		tv.wnames[i], err = r.ReadString()
		if err != nil {
			return nil, err
		}
	}
	return tv, nil

}

// wire format:
//
//	size[4] Rwalk tag[2] nwqid[2] nwqid*(wqid[13])
func (tv *RWalk) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	numQids, err := r.Read16()
	tv.Wqids = make([]Qid, numQids)
	for i := 0; i < int(numQids); i++ {
		tv.Wqids[i], err = r.ReadQid()
		if err != nil {
			return nil, err
		}
	}

	return tv, err
}

// wire format:
//
// size[4] Tflush tag[2] oldtag[2]
func (tv *TFlush) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Oldtag, err = r.ReadTag()
	return tv, err

}

func (tv *RFlush) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	return tv, err
}

func (tv *TVersion) fillFrom(r TypedReader) (FCall, error) {
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

// size[4] Rversion tag[2] msize[4] version[s]
func (tv *RVersion) fillFrom(r TypedReader) (FCall, error) {
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
