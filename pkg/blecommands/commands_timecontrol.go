package blecommands

import (
	"fmt"
	"slices"
)

// TimeControlEntryID (0x003A)

var _ Response = &TimeControlEntryID{}

type TimeControlEntryID struct {
	EntryID byte `json:"entryId"`
}

func (c *TimeControlEntryID) GetCommandCode() CommandCode { return CommandTimeControlEntryID }
func (c *TimeControlEntryID) FromMessage(b []byte) error {
	if len(b) < 1 {
		return fmt.Errorf("time control entry ID too short: %d bytes", len(b))
	}
	c.EntryID = b[0]
	return nil
}

// TimeControlEntryCount (0x003D)

var _ Response = &TimeControlEntryCount{}

type TimeControlEntryCount struct {
	Count byte `json:"count"`
}

func (c *TimeControlEntryCount) GetCommandCode() CommandCode { return CommandTimeControlEntryCount }
func (c *TimeControlEntryCount) FromMessage(b []byte) error {
	if len(b) < 1 {
		return fmt.Errorf("time control entry count too short: %d bytes", len(b))
	}
	c.Count = b[0]
	return nil
}

// TimeControlEntry (0x003E)
// Layout: EntryID(1) + Enabled(1) + Weekdays(1) + Time(2) + LockAction(1) = 6 bytes

var _ Response = &TimeControlEntry{}

type TimeControlEntry struct {
	EntryID    byte   `json:"entryId"`
	Enabled    bool   `json:"enabled"`
	Weekdays   byte   `json:"weekdays"`
	TimeHour   byte   `json:"timeHour"`
	TimeMinute byte   `json:"timeMinute"`
	LockAction Action `json:"lockAction"`
}

func (c *TimeControlEntry) GetCommandCode() CommandCode { return CommandTimeControlEntry }
func (c *TimeControlEntry) FromMessage(b []byte) error {
	if len(b) < 6 {
		return fmt.Errorf("time control entry too short: %d bytes, need 6", len(b))
	}
	c.EntryID = b[0]
	c.Enabled = b[1] != 0
	c.Weekdays = b[2]
	c.TimeHour = b[3]
	c.TimeMinute = b[4]
	c.LockAction = Action(b[5])
	return nil
}

// AddTimeControlEntry (0x0039)

var _ Request = &AddTimeControlEntry{}

type AddTimeControlEntry struct {
	Weekdays    byte
	TimeHour    byte
	TimeMinute  byte
	LockAction  Action
	Nonce       []byte
	SecurityPin Pin
}

func (c *AddTimeControlEntry) GetCommandCode() CommandCode { return CommandAddTimeControlEntry }
func (c *AddTimeControlEntry) GetPayload() []byte {
	return slices.Concat(
		[]byte{c.Weekdays, c.TimeHour, c.TimeMinute, byte(c.LockAction)},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// RemoveTimeControlEntry (0x003B)

var _ Request = &RemoveTimeControlEntry{}

type RemoveTimeControlEntry struct {
	EntryID     byte
	Nonce       []byte
	SecurityPin Pin
}

func (c *RemoveTimeControlEntry) GetCommandCode() CommandCode { return CommandRemoveTimeControlEntry }
func (c *RemoveTimeControlEntry) GetPayload() []byte {
	return slices.Concat(
		[]byte{c.EntryID},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// RequestTimeControlEntries (0x003C)

var _ Request = &RequestTimeControlEntries{}

type RequestTimeControlEntries struct {
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestTimeControlEntries) GetCommandCode() CommandCode {
	return CommandRequestTimeControlEntries
}
func (c *RequestTimeControlEntries) GetPayload() []byte {
	return slices.Concat(c.Nonce, c.SecurityPin.GetPinBytes())
}

// UpdateTimeControlEntry (0x003F)

var _ Request = &UpdateTimeControlEntry{}

type UpdateTimeControlEntry struct {
	EntryID     byte
	Enabled     bool
	Weekdays    byte
	TimeHour    byte
	TimeMinute  byte
	LockAction  Action
	Nonce       []byte
	SecurityPin Pin
}

func (c *UpdateTimeControlEntry) GetCommandCode() CommandCode { return CommandUpdateTimeControlEntry }
func (c *UpdateTimeControlEntry) GetPayload() []byte {
	return slices.Concat(
		[]byte{c.EntryID, boolToByte(c.Enabled), c.Weekdays, c.TimeHour, c.TimeMinute, byte(c.LockAction)},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}
