package dlog

import (
	"strconv"
	"time"
)

func getDelta(prevTime time.Time) (now time.Time, delta time.Duration) {
	now = time.Now()
	delta = now.Sub(prevTime)
	return
}

func humanizeNano(n time.Duration) string {
	var suffix string

	switch {
	case n > 1e9:
		n /= 1e9
		suffix = "s"
	case n > 1e6:
		n /= 1e6
		suffix = "ms"
	case n > 1e3:
		n /= 1e3
		suffix = "us"
	default:
		suffix = "ns"
	}

	return strconv.Itoa(int(n)) + suffix
}
