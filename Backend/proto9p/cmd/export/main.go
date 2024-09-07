// export a directory to the given port
package main

import (
	"PumiceBackend/proto9p/fs"
	"PumiceBackend/server"
	"bufio"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
)

type Config struct {
	port      int
	directory string
}

func parseArgs() Config {
	portPtr := flag.Int("port", 9000, "port to serve on")
	dirPtr := flag.String("directory", "", "directory to save")
	flag.Parse()

	var c Config
	c.port = *portPtr
	c.directory = *dirPtr
	return c
}
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	cfg := parseArgs()
	if cfg.directory == "" {
		fmt.Println("choose a directory to serve")
		os.Exit(1)
	}

	fileServer, err := fs.NewServer(cfg.directory)
	if err != nil {
		log.Fatal("Couldnt create file server: ", err)
	}

	addr := fmt.Sprintf(":%d", cfg.port)
	slog.Debug("listening", "addr", addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		a, err := listener.Accept()
		slog.Debug("Accepted connection", "local_addr", a.LocalAddr())
		if err != nil {
			panic(err)
		}
		go func(nc net.Conn, srv server.Server) {
			defer nc.Close()
			read := bufio.NewReader(nc)
			err := server.Serve(read, nc, srv)
			if err != nil {
				log.Printf("%v\n", err)
			}
		}(a, &fileServer)
	}
}
