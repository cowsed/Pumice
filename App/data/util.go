package data

import (
	"fmt"
	"path"
)

type OSPath string
type VaultLocation string

func (v VaultLocation) Name() VaultLocation {
	return VaultLocation(path.Base(string(v)))
}
func (v VaultLocation) Dir() VaultLocation {
	return VaultLocation(path.Dir(string(v)))
}

func ToOSPath(vault OSPath, loc VaultLocation) string {
	return path.Join(string(vault), string(loc))
}
func (op OSPath) Base() string {
	return path.Base(string(op))
}
func (op OSPath) String() string {
	return string(op)
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
func (ts TagSet) StringList() []string {
	l := ts.List()
	s := make([]string, len(l))
	for i, t := range l {
		s[i] = t.String()
	}
	return s
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
