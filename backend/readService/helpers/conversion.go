package helpers

import "strconv"

func UintToString(input uint) string {
	return strconv.FormatUint(uint64(input), 10)
}