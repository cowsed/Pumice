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
	case Tauth:
		return (&TAuth{}).fillFrom(packet_reader)
	case Rauth:
		return (&RAuth{}).fillFrom(packet_reader)
	case Rerror:
		return (&RError{}).fillFrom(packet_reader)
	case Tflush:
		return (&TFlush{}).fillFrom(packet_reader)
	case Rflush:
		return (&RFlush{}).fillFrom(packet_reader)
	case Tattach:
		return (&TAttach{}).fillFrom(packet_reader)
	case Rattach:
		return (&RAttach{}).fillFrom(packet_reader)
	case Twalk:
		return (&TWalk{}).fillFrom(packet_reader)
	case Rwalk:
		return (&RWalk{}).fillFrom(packet_reader)

	case Topen:
		return (&TOpen{}).fillFrom(packet_reader)
	case Ropen:
		return (&ROpen{}).fillFrom(packet_reader)
	case Tcreate:
		return (&TCreate{}).fillFrom(packet_reader)
	case Rcreate:
		return (&RCreate{}).fillFrom(packet_reader)
	case Tread:
		return (&TRead{}).fillFrom(packet_reader)
	case Rread:
		return (&RRead{}).fillFrom(packet_reader)
	case Twrite:
		return (&TWrite{}).fillFrom(packet_reader)
	case Rwrite:
		return (&RWrite{}).fillFrom(packet_reader)
	case Tclunk:
		return (&TClunk{}).fillFrom(packet_reader)
	case Rclunk:
		return (&RClunk{}).fillFrom(packet_reader)
	case Tremove:
		return (&TRemove{}).fillFrom(packet_reader)
	case Rremove:
		return (&RRemove{}).fillFrom(packet_reader)

	case Tstat:
		return (&TStat{}).fillFrom(packet_reader)
	case Rstat:
		return (&RStat{}).fillFrom(packet_reader)
	case Twstat:
		return (&TWStat{}).fillFrom(packet_reader)
	case Rwstat:
		return (&RWStat{}).fillFrom(packet_reader)
	default:
		return nil, fmt.Errorf("got %d: %w", pt8, ErrUnknownType)
	}

}

func (t *TStat) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (t *RStat) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Stat, err = r.ReadStat()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TWStat) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	t.Stat, err = r.ReadStat()
	if err != nil {
		return nil, err
	}

	return t, nil
}
func (t *RWStat) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TClunk) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (t *RClunk) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TRemove) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (t *RRemove) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TCreate) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	t.Name, err = r.ReadString()
	if err != nil {
		return nil, err
	}
	perm32, err := r.Read32()
	if err != nil {
		return nil, err
	}
	t.perm = Perm(perm32)

	m8, err := r.Read8()
	if err != nil {
		return nil, err
	}
	t.Mode = Mode(m8)

	return t, nil
}

func (tv *RCreate) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Qid, err = r.ReadQid()
	if err != nil {
		return nil, err
	}

	tv.IOUnit, err = r.Read32()
	if err != nil {
		return nil, err
	}

	return tv, nil
}

func (t *TOpen) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	m8, err := r.Read8()
	if err != nil {
		return nil, err
	}
	t.Mode = Mode(m8)

	return t, nil
}

func (tv *ROpen) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Qid, err = r.ReadQid()
	if err != nil {
		return nil, err
	}

	tv.IOUnit, err = r.Read32()
	if err != nil {
		return nil, err
	}

	return tv, nil
}

func (t *RError) fillFrom(r TypedReader) (FCall, error) {
	var err error
	t.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	t.Ename, err = r.ReadString()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (tv *TAuth) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.afid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	tv.aname, err = r.ReadString()
	if err != nil {
		return nil, err
	}

	tv.uname, err = r.ReadString()
	if err != nil {
		return nil, err
	}

	return tv, nil
}
func (tv *RAuth) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	tv.aqid, err = r.ReadQid()
	if err != nil {
		return nil, err
	}

	return tv, nil
}
func (tv *TAttach) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	tv.Afid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}

	tv.Uname, err = r.ReadString()
	if err != nil {
		return nil, err
	}

	tv.Aname, err = r.ReadString()
	if err != nil {
		return nil, err
	}
	return tv, nil
}

func (tv *RAttach) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	tv.Qid, err = r.ReadQid()
	if err != nil {
		return nil, err
	}
	return tv, nil
}

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
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}
	tv.Fid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	tv.NewFid, err = r.ReadFid()
	if err != nil {
		return nil, err
	}
	numWname, err := r.Read16()
	if err != nil {
		return nil, err
	}
	tv.WNames = make([]string, numWname)
	for i := 0; i < int(numWname); i++ {
		tv.WNames[i], err = r.ReadString()
		if err != nil {
			return nil, err
		}
	}
	return tv, nil

}

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
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	tv.MSize, err = r.Read32()
	if err != nil {
		return nil, err
	}
	tv.Version, err = r.ReadString()
	return tv, err
}

func (tv *RVersion) fillFrom(r TypedReader) (FCall, error) {
	var err error
	tv.Tag, err = r.ReadTag()
	if err != nil {
		return nil, err
	}

	tv.Msize, err = r.Read32()
	if err != nil {
		return nil, err
	}
	tv.Version, err = r.ReadString()
	return tv, err
}
