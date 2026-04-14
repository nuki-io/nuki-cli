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
	TotalCount  bool
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
		[]byte{byte(c.SortOrder), boolToByte(c.TotalCount)},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

//go:generate stringer -type=LogEntryType
type LogEntryType uint8

func (t LogEntryType) MarshalText() ([]byte, error) { return []byte(t.String()), nil }

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

// LogEntry (0x0032)

var _ Response = &LogEntry{}

type LogEntry struct {
	Index    uint32       `json:"index"`
	Time     time.Time    `json:"time"`
	AuthId   uint32       `json:"authId"`
	AuthName string       `json:"authName"`
	Type     LogEntryType `json:"type"`
	Data     []byte       `json:"data"`
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

// LogEntryCount (0x0033)

var _ Response = &LogEntryCount{}

type LogEntryCount struct {
	LoggingEnabled           bool   `json:"loggingEnabled"`
	Count                    uint16 `json:"count"`
	DoorSensorEnabled        bool   `json:"doorSensorEnabled"`
	DoorSensorLoggingEnabled bool   `json:"doorSensorLoggingEnabled"`
}

func (c *LogEntryCount) GetCommandCode() CommandCode { return CommandLogEntryCount }
func (c *LogEntryCount) FromMessage(b []byte) error {
	if len(b) < 4 {
		return fmt.Errorf("log entry count too short: %d bytes", len(b))
	}
	c.LoggingEnabled = b[0] != 0
	c.Count = binary.LittleEndian.Uint16(b[1:3])
	c.DoorSensorEnabled = b[3] != 0
	if len(b) > 4 {
		c.DoorSensorLoggingEnabled = b[4] != 0
	}
	return nil
}

// EnableLogging (0x0034)

var _ Request = &EnableLogging{}

type EnableLogging struct {
	Enabled     bool
	Nonce       []byte
	SecurityPin Pin
}

func (c *EnableLogging) GetCommandCode() CommandCode { return CommandEnableLogging }
func (c *EnableLogging) GetPayload() []byte {
	return slices.Concat([]byte{boolToByte(c.Enabled)}, c.Nonce, c.SecurityPin.GetPinBytes())
}
