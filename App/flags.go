package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cowsed/Pumice/App/data"
)

type Flags struct {
	VaultPath data.OSPath
}

func parseFlags() Flags {
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("This app requires one command line argument. The directory of the vault to open")
		os.Exit(1)
	}
	return Flags{
		VaultPath: data.OSPath(args[0]),
	}
}
