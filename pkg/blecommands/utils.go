package blecommands

import (
	"encoding/binary"
	"time"
)

func byteToBool(b byte) bool {
	return b != 0
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func toNukiTime(t time.Time) []byte {
	b := make([]byte, 7)
	binary.LittleEndian.PutUint16(b[0:2], uint16(t.Year()))
	b[2] = byte(t.Month())
	b[3] = byte(t.Day())
	b[4] = byte(t.Hour())
	b[5] = byte(t.Minute())
	b[6] = byte(t.Second())
	return b
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
