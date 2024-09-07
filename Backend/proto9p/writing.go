package proto9p

func (tv *TStat) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tstat)
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
	return nil
}
func (tv *RStat) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rstat)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.Write32(uint32(tv.Stat))
	if err != nil {
		return err
	}
	return nil
}
func (tv *TWStat) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Twstat)
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
	err = tw.Write32(uint32(tv.Stat))
	if err != nil {
		return err
	}
	return nil
}
func (tv *RWStat) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rwstat)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TClunk) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tclunk)
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
	return nil
}
func (tv *TRemove) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tremove)
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
	return nil
}

func (tv *RClunk) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rclunk)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	return nil
}
func (tv *RRemove) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rremove)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TCreate) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tcreate)
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
	err = tw.WriteString(tv.Name)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.perm)
	if err != nil {
		return err
	}
	err = tw.Write8(uint8(tv.Mode))
	if err != nil {
		return err
	}
	return nil
}

func (tv *RCreate) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rcreate)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteQid(tv.Qid)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.IOUnit)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TOpen) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Topen)
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
	err = tw.Write8(uint8(tv.Mode))
	if err != nil {
		return err
	}
	return nil
}

func (tv *ROpen) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Ropen)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteQid(tv.Qid)
	if err != nil {
		return err
	}
	err = tw.Write32(tv.IOUnit)
	if err != nil {
		return err
	}
	return nil
}

func (tv *RError) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rerror)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.ename)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TAuth) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tauth)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteFid(tv.afid)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.uname)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.aname)
	if err != nil {
		return err
	}
	return nil
}

func (tv *RAuth) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rauth)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteQid(tv.aqid)
	if err != nil {
		return err
	}
	return nil
}

func (tv *TAttach) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Tattach)
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
	err = tw.WriteFid(tv.afid)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.uname)
	if err != nil {
		return err
	}
	err = tw.WriteString(tv.aname)
	if err != nil {
		return err
	}
	return nil
}
func (tv *RAttach) writeTo(tw TypedWriter) error {
	err := tw.WriteType(Rattach)
	if err != nil {
		return err
	}
	err = tw.WriteTag(tv.Tag)
	if err != nil {
		return err
	}
	err = tw.WriteQid(tv.Qid)
	if err != nil {
		return err
	}
	return nil
}

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
