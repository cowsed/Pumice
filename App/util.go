package main

import (
	"fmt"
	"path"
)

type OSPath string
type VaultLocation string

var CURRENT_VERSION = Version{
	major:   0,
	minor:   0,
	patch:   0,
	comment: "dev",
}

type Version struct {
	major   int
	minor   int
	patch   int
	comment string
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.patch, v.comment)
}

func ToOSPath(vault OSPath, loc VaultLocation) string {
	return path.Join(string(vault), string(loc))
}

func (vl VaultLocation) Append(file string) VaultLocation {
	return VaultLocation(path.Join(string(vl), file))
}
