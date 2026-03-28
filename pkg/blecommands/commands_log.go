package blecommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"time"
)

//go:generate stringer -type=LogSortOrder -trimprefix=LogSortOrder
type LogSortOrder byte

const (
	LogSortOrderAscending  LogSortOrder = 0x00
	LogSortOrderDescending LogSortOrder = 0x01
)

var _ Request = &RequestLogEntries{}

type RequestLogEntries struct {
	StartIndex  uint32
	Count       uint16
	SortOrder   LogSortOrder
	TotalCount  byte
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestLogEntries) GetCommandCode() CommandCode {
	return CommandRequestLogEntries
}

func (c *RequestLogEntries) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint32(nil, c.StartIndex),
		binary.LittleEndian.AppendUint16(nil, c.Count),
		[]byte{byte(c.SortOrder), c.TotalCount},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

//go:generate stringer -type=LogEntryType
type LogEntryType uint8

const (
	LoggingEnabledDisabled           LogEntryType = 0x01
	LogLockAction                    LogEntryType = 0x02
	LogCalibration                   LogEntryType = 0x03
	LogInitializationRun             LogEntryType = 0x04
	LogKeypadAction                  LogEntryType = 0x05
	LogDoorSensor                    LogEntryType = 0x06
	DoorSensorLoggingEnabledDisabled LogEntryType = 0x07
	LogFirmwareUpdate                LogEntryType = 0x0A
)

var _ Response = &LogEntry{}

type LogEntry struct {
	Index    uint32
	Time     time.Time
	AuthId   uint32
	AuthName string
	Type     LogEntryType
	Data     []byte
}

func (c *LogEntry) FromMessage(b []byte) error {
	if len(b) < 48 {
		return fmt.Errorf("log entry length must be at least 48 bytes, got: %d", len(b))
	}
	c.Index = binary.LittleEndian.Uint32(b[0:4])
	c.Time = fromNukiTime(b[4:11], time.UTC)
	c.AuthId = binary.LittleEndian.Uint32(b[11:15])
	c.AuthName = string(bytes.Trim(b[15:47], "\x00"))
	c.Type = LogEntryType(b[47])
	c.Data = b[48:]
	return nil
}

func (c *LogEntry) GetCommandCode() CommandCode {
	return CommandLogEntry
}

func (c *LogEntry) String() string {
	switch c.Type {
	case LoggingEnabledDisabled:
		if c.Data[0] == 0 {
			return "Logging Disabled"
		} else {
			return "Logging Enabled"
		}
	case LogLockAction, LogCalibration, LogInitializationRun:
		var suffix string
		if t := Trigger(c.Data[1]).String(); t != "" {
			suffix = fmt.Sprintf(" (%s)", t)
		}
		return fmt.Sprintf("%v%s", Action(c.Data[0]), suffix)
	}
	if c.Type == LogFirmwareUpdate {
		return fmt.Sprintf("Firmware Update (%d.%d.%d)", c.Data[0], c.Data[1], c.Data[2])
	}
	return fmt.Sprintf("%s by %s (ID: %d)", c.Type.String(), c.AuthName, c.AuthId)
}
