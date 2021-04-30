package global

import (
	"math"
	"strconv"
)

/* constant */
const (
	PAGE_SIZE             = 10
	PREVIEW_STRING_LENGTH = 30
)

/* function utils */
func Ceil(n1, n2 int64) int64 {
	return int64(math.Ceil(float64(n1) / float64(n2)))
}

func Int64ToInt(x int64) int {
	strInt64 := strconv.FormatInt(x, 10)
	ret := 0
	if r, e := strconv.Atoi(strInt64); e == nil {
		ret = r
	}
	return ret
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
