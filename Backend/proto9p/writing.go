package proto9p

// size[4] Tversion tag[2] msize[4] version[s]
func (tv *TVersion) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tversion)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.msize)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.version)
	if err != nil {
		return err
	}
	return nil
}

// size[4] Tversion tag[2] msize[4] version[s]
func (tv *RVersion) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rversion)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.msize)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.version)
	if err != nil {
		return err
	}
	return nil
}

// size[4] Tflush tag[2] oldtag[2]
func (tv *TFlush) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tflush)
	if err != nil {
		return err
	}

	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Oldtag)
	if err != nil {
		return err
	}
	return nil
}

// size[4] Rflush tag[2]
func (tv *RFlush) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rflush)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TWalk) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Twalk)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteFid(tv.Fid)
	if err != nil {
		return err
	}

	err = tw.WriteFid(tv.newFid)
	if err != nil {
		return err
	}

	err = tw.Write16(uint16(len(tv.wnames)))
	if err != nil {
		return err
	}
	for _, wname := range tv.wnames {
		err = tw.WriteString(wname)
		if err != nil {
			return err
		}
	}
	return nil

}
func (tv *RWalk) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rwalk)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.Write16(uint16(len(tv.Wqids)))
	if err != nil {
		return err
	}

	for _, qid := range tv.Wqids {
		err = tw.WriteQid(qid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tv *TWrite) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Twrite)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteFid(tv.Fid)
	if err != nil {
		return err
	}
	err = tw.Write64(tv.Offset)
	if err != nil {
		return err
	}
	err = tw.Write32(uint32(len(tv.Data)))
	if err != nil {
		return err
	}
	err = tw.WriteN(tv.Data)
	if err != nil {
		return err
	}
	return nil
}

func (tv *RWrite) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rwrite)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.Count)
	if err != nil {
		return err
	}

	return nil
}

func (tv *TRead) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tread)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteFid(tv.Fid)
	if err != nil {
		return err
	}
	err = tw.Write64(tv.Offset)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.Count)
	if err != nil {
		return err
	}
	return nil
}

func (tv *RRead) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rread)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.Write32(uint32(len(tv.Data)))
	if err != nil {
		return err
	}
	err = tw.WriteN(tv.Data)
	if err != nil {
		return err
	}
	return nil
}
