package server

import (
	"PumiceBackend/proto9p"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
)

var ErrServerReceivedRMessage = errors.New("server received R message")

func NewErrServerRMessage(fc proto9p.FCall) error {
	return fmt.Errorf("got %v (type %T): %w", fc, fc, ErrServerReceivedRMessage)
}

type Conn interface {
}

type Server interface {
	NewConnection() Conn
	Version(Conn, *proto9p.TVersion) (proto9p.FCall, error)
	Auth(Conn, *proto9p.TAuth) (proto9p.FCall, error)
	Attach(Conn, *proto9p.TAttach) (proto9p.FCall, error)
	Walk(Conn, *proto9p.TWalk) (proto9p.FCall, error)
	Open(Conn, *proto9p.TOpen) (proto9p.FCall, error)
	Create(Conn, *proto9p.TCreate) (proto9p.FCall, error)
	Flush(Conn, *proto9p.TFlush) (proto9p.FCall, error)
	Read(Conn, *proto9p.TRead) (proto9p.FCall, error)
	Write(Conn, *proto9p.TWrite) (proto9p.FCall, error)
	Clunk(Conn, *proto9p.TClunk) (proto9p.FCall, error)
	Remove(Conn, *proto9p.TRemove) (proto9p.FCall, error)
	Stat(Conn, *proto9p.TStat) (proto9p.FCall, error)
	Wstat(Conn, *proto9p.TWStat) (proto9p.FCall, error)
}

func writeToWire(w io.Writer, toWire chan proto9p.FCall, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(toWire)

	for call := range toWire {
		slog.Debug("Writing fcall", "fcall", call)

		bs, err := proto9p.WriteFCall(call)
		if err != nil {
			slog.Error("Error serializing packet", "err", err)
			continue
		}
		n, err := w.Write(bs)
		if n < len(bs) {
			slog.Error("Could not write entire packet to wire", "err", err)
			continue
		}
		if err != nil {
			slog.Error("Could not write packet to wire", "err", err)
			continue
		}
	}
}

func takeCalls(fromWire, toWire chan proto9p.FCall, server Server, conn Conn, wg *sync.WaitGroup) {
	defer close(toWire)
	defer wg.Done()

	for call := range fromWire {
		call, err := handleFCall(call, server, conn)
		if err != nil {
			slog.Error("Error handling fcall", "err", err)
		}
		toWire <- call
	}

}

// Serve a single channel of data
// this may be one of many channels of data if server can handle multiple requests
// the server needs to ne handle the multiple clients (identified by Conn) if is is used in this way
func Serve(r io.Reader, w io.Writer, server Server) error {
	readChan := make(chan proto9p.FCall, 10)
	writeChan := make(chan proto9p.FCall, 10)
	conn := server.NewConnection()
	// reading thread
	var wg sync.WaitGroup
	wg.Add(2)
	defer func() { wg.Wait() }()

	go takeCalls(readChan, writeChan, server, conn, &wg)
	go writeToWire(w, writeChan, &wg)

	defer close(readChan)

	for {
		call, err := proto9p.ParseFCall(r)
		if err != nil {
			slog.Error("Reading fcall", "err", err)
			return err
		}
		slog.Debug("read call", "call", call)
		readChan <- call
	}

}

func handleFCall(call proto9p.FCall, srv Server, conn Conn) (proto9p.FCall, error) {
	switch conc := call.(type) {

	case *proto9p.TAttach:
		return srv.Attach(conn, conc)
	case *proto9p.TAuth:
		return srv.Auth(conn, conc)
	case *proto9p.TClunk:
		return srv.Clunk(conn, conc)
	case *proto9p.TCreate:
		return srv.Create(conn, conc)
	case *proto9p.TFlush:
		return srv.Flush(conn, conc)
	case *proto9p.TOpen:
		return srv.Open(conn, conc)
	case *proto9p.TRead:
		return srv.Read(conn, conc)
	case *proto9p.TRemove:
		return srv.Remove(conn, conc)
	case *proto9p.TStat:
		return srv.Stat(conn, conc)
	case *proto9p.TVersion:
		return srv.Version(conn, conc)
	case *proto9p.TWStat:
		return srv.Wstat(conn, conc)
	case *proto9p.TWalk:
		return srv.Walk(conn, conc)
	case *proto9p.TWrite:
		return srv.Write(conn, conc)
	default:
		return nil, NewErrServerRMessage(conc)
	}
}
