package finance

import (
	"fmt"
	"strconv"
	"strings"
)

func FenToYuan(fen int64) (yuan string) {
	if fen == 0 {
		return "0"
	}
	fl := fen / 100
	fr := fen % 100
	if fr == 0 {
		return strconv.FormatInt(fl, 10)
	}
	return strings.TrimSuffix(fmt.Sprintf("%d.%02d", fl, fr), "0")
}
