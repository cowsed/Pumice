package config

import "fmt"

var CURRENT_VERSION_MAJOR = 0
var CURRENT_VERSION_MINOR = 0
var CURRENT_VERSION_PATCH = 0
var CURRENT_VERSION_COMMENT = "dev"

var VERSION = Version{
	major:   CURRENT_VERSION_MAJOR,
	minor:   CURRENT_VERSION_MINOR,
	patch:   CURRENT_VERSION_PATCH,
	comment: CURRENT_VERSION_COMMENT,
}

type Version struct {
	major   int
	minor   int
	patch   int
	comment string
}

func (v Version) String() string {
	if v.comment != "" {
		return fmt.Sprintf("%d.%d.%d-%s", v.major, v.minor, v.patch, v.comment)
	}
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}
