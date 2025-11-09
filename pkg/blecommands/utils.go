package blecommands

import (
	"encoding/binary"
	"time"
)

func byteToBool(b byte) bool {
	return b != 0
}

func fromNukiTime(b []byte, tz *time.Location) time.Time {
	if len(b) != 7 {
		return time.Time{}
	}
	year := int(binary.LittleEndian.Uint16(b[0:2]))
	month := int(b[2])
	day := int(b[3])
	hour := int(b[4])
	minute := int(b[5])
	second := int(b[6])

	time := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC).In(tz)
	return time
}
