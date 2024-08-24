package main

import (
	"flag"
	"fmt"
	"os"
)

type Flags struct {
	VaultPath OSPath
}

func parseFlags() Flags {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("This app requires one command line argument. The directory of the vault to open")
		os.Exit(1)
	}
	return Flags{
		VaultPath: OSPath(args[0]),
	}
}
