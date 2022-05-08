package utils

import (
	"time"
)

var (
	utc *time.Location
)

func init() {
	utc = time.FixedZone("UTC", 0)
}

func UTC() *time.Location {
	return utc
}
