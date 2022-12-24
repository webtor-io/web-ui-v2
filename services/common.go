package services

import "regexp"

var SHA1R = regexp.MustCompile("(?i)[0-9a-f]{5,40}")
