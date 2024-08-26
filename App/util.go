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
func (op OSPath) Base() string {
	return path.Base(string(op))
}

func (vl VaultLocation) Append(file string) VaultLocation {
	return VaultLocation(path.Join(string(vl), file))
}

type TagSet struct {
	internal map[Tag]struct{}
}

func (ts TagSet) String() string {
	return fmt.Sprint(ts.internal)
}
func (ts TagSet) List() []Tag {
	keys := make([]Tag, 0, len(ts.internal))
	for k := range ts.internal {
		keys = append(keys, k)
	}
	return keys
}

func (ts TagSet) Contains(t Tag) bool {
	_, contains := ts.internal[t]
	return contains
}

func (ts TagSet) Len() int { return len(ts.internal) }

func NewTagSet() TagSet {
	return TagSet{
		internal: map[Tag]struct{}{},
	}
}
func (ts *TagSet) Add(tag Tag) {
	ts.internal[tag] = struct{}{}
}
