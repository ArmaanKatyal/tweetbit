package helpers

import (
	"strconv"
	"time"
)

func ConvertStringToUint(s string) uint {
	u64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(u64)
}

func ConvertStringToTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
