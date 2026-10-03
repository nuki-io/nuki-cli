package blecommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"time"
)

// KeypadCodeID (0x0042)

var _ Response = &KeypadCodeID{}

type KeypadCodeID struct {
	CodeID      uint16    `json:"codeId"`
	DateCreated time.Time `json:"dateCreated"`
}

func (c *KeypadCodeID) GetCommandCode() CommandCode { return CommandKeypadCodeID }
func (c *KeypadCodeID) FromMessage(b []byte) error {
	if len(b) < 9 {
		return fmt.Errorf("keypad code ID too short: %d bytes, need 9", len(b))
	}
	c.CodeID = binary.LittleEndian.Uint16(b[0:2])
	c.DateCreated = fromNukiTime(b[2:9], time.UTC)
	return nil
}

// KeypadCodeCount (0x0044)

var _ Response = &KeypadCodeCount{}

type KeypadCodeCount struct {
	Count uint16 `json:"count"`
}

func (c *KeypadCodeCount) GetCommandCode() CommandCode { return CommandKeypadCodeCount }
func (c *KeypadCodeCount) FromMessage(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("keypad code count too short: %d bytes", len(b))
	}
	c.Count = binary.LittleEndian.Uint16(b[0:2])
	return nil
}

// KeypadCode (0x0045)
// Layout: CodeID(2) + Code(4) + Name(20) + Enabled(1) + DateCreated(7) + DateLastActive(7) +
//         LockCount(2) + TimeLimited(1) = 44 bytes base
// If TimeLimited: + AllowedFromDate(7) + AllowedUntilDate(7) + Weekdays(1) + FromTime(2) + UntilTime(2) = +19

var _ Response = &KeypadCode{}

type KeypadCode struct {
	CodeID           uint16    `json:"codeId"`
	Code             uint32    `json:"code"`
	Name             string    `json:"name"`
	Enabled          bool      `json:"enabled"`
	DateCreated      time.Time `json:"dateCreated"`
	DateLastActive   time.Time `json:"dateLastActive"`
	LockCount        uint16    `json:"lockCount"`
	TimeLimited      bool      `json:"timeLimited"`
	AllowedFromDate  time.Time `json:"allowedFromDate,omitempty"`
	AllowedUntilDate time.Time `json:"allowedUntilDate,omitempty"`
	AllowedWeekdays  byte      `json:"allowedWeekdays,omitempty"`
	AllowedFromTime  [2]byte   `json:"allowedFromTime,omitempty"`
	AllowedUntilTime [2]byte   `json:"allowedUntilTime,omitempty"`
}

func (c *KeypadCode) GetCommandCode() CommandCode { return CommandKeypadCode }
func (c *KeypadCode) FromMessage(b []byte) error {
	if len(b) < 44 {
		return fmt.Errorf("keypad code too short: %d bytes, need at least 44", len(b))
	}
	c.CodeID = binary.LittleEndian.Uint16(b[0:2])
	c.Code = binary.LittleEndian.Uint32(b[2:6])
	c.Name = string(bytes.Trim(b[6:26], "\x00"))
	c.Enabled = b[26] != 0
	c.DateCreated = fromNukiTime(b[27:34], time.UTC)
	c.DateLastActive = fromNukiTime(b[34:41], time.UTC)
	c.LockCount = binary.LittleEndian.Uint16(b[41:43])
	c.TimeLimited = b[43] != 0
	if c.TimeLimited && len(b) >= 63 {
		c.AllowedFromDate = fromNukiTime(b[44:51], time.UTC)
		c.AllowedUntilDate = fromNukiTime(b[51:58], time.UTC)
		c.AllowedWeekdays = b[58]
		c.AllowedFromTime = [2]byte{b[59], b[60]}
		c.AllowedUntilTime = [2]byte{b[61], b[62]}
	}
	return nil
}

// RequestKeypadCodes (0x0043)

var _ Request = &RequestKeypadCodes{}

type RequestKeypadCodes struct {
	Offset      uint16
	Count       uint16
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestKeypadCodes) GetCommandCode() CommandCode { return CommandRequestKeypadCodes }
func (c *RequestKeypadCodes) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint16(nil, c.Offset),
		binary.LittleEndian.AppendUint16(nil, c.Count),
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// AddKeypadCode (0x0041)
// Time fields are always present in the payload (fixed-size protocol message).

var _ Request = &AddKeypadCode{}

type AddKeypadCode struct {
	Code             uint32
	Name             string // max 20 chars
	TimeLimited      bool
	AllowedFromDate  time.Time
	AllowedUntilDate time.Time
	AllowedWeekdays  byte
	AllowedFromTime  [2]byte
	AllowedUntilTime [2]byte
	Nonce            []byte
	SecurityPin      Pin
}

func (c *AddKeypadCode) GetCommandCode() CommandCode { return CommandAddKeypadCode }
func (c *AddKeypadCode) GetPayload() []byte {
	name := [20]byte{}
	copy(name[:], c.Name)
	return slices.Concat(
		binary.LittleEndian.AppendUint32(nil, c.Code),
		name[:],
		[]byte{boolToByte(c.TimeLimited)},
		toNukiTime(c.AllowedFromDate),
		toNukiTime(c.AllowedUntilDate),
		[]byte{c.AllowedWeekdays},
		c.AllowedFromTime[:],
		c.AllowedUntilTime[:],
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// UpdateKeypadCode (0x0046)

var _ Request = &UpdateKeypadCode{}

type UpdateKeypadCode struct {
	CodeID           uint16
	Code             uint32
	Name             string // max 20 chars
	Enabled          bool
	TimeLimited      bool
	AllowedFromDate  time.Time
	AllowedUntilDate time.Time
	AllowedWeekdays  byte
	AllowedFromTime  [2]byte
	AllowedUntilTime [2]byte
	Nonce            []byte
	SecurityPin      Pin
}

func (c *UpdateKeypadCode) GetCommandCode() CommandCode { return CommandUpdateKeypadCode }
func (c *UpdateKeypadCode) GetPayload() []byte {
	name := [20]byte{}
	copy(name[:], c.Name)
	return slices.Concat(
		binary.LittleEndian.AppendUint16(nil, c.CodeID),
		binary.LittleEndian.AppendUint32(nil, c.Code),
		name[:],
		[]byte{boolToByte(c.Enabled), boolToByte(c.TimeLimited)},
		toNukiTime(c.AllowedFromDate),
		toNukiTime(c.AllowedUntilDate),
		[]byte{c.AllowedWeekdays},
		c.AllowedFromTime[:],
		c.AllowedUntilTime[:],
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// RemoveKeypadCode (0x0047)

var _ Request = &RemoveKeypadCode{}

type RemoveKeypadCode struct {
	CodeID      uint16
	Nonce       []byte
	SecurityPin Pin
}

func (c *RemoveKeypadCode) GetCommandCode() CommandCode { return CommandRemoveKeypadCode }
func (c *RemoveKeypadCode) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint16(nil, c.CodeID),
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}
