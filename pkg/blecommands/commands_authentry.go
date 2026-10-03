package blecommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"time"
)

// AuthorizationEntryCount (0x0027)

var _ Response = &AuthorizationEntryCount{}

type AuthorizationEntryCount struct {
	Count uint16 `json:"count"`
}

func (c *AuthorizationEntryCount) GetCommandCode() CommandCode { return CommandAuthorizationEntryCount }
func (c *AuthorizationEntryCount) FromMessage(b []byte) error {
	if len(b) < 2 {
		return fmt.Errorf("authorization entry count too short: %d bytes", len(b))
	}
	c.Count = binary.LittleEndian.Uint16(b[0:2])
	return nil
}

// AuthorizationEntry (0x000A)
// Layout: AuthID(4) + IDType(1) + Name(32) + Enabled(1) + RemoteAllowed(1) +
//         DateCreated(7) + DateLastActive(7) + LockCount(2) + TimeLimited(1) = 56 bytes base
// If TimeLimited: + AllowedFromDate(7) + AllowedUntilDate(7) + Weekdays(1) + FromTime(2) + UntilTime(2) = +19

var _ Response = &AuthorizationEntry{}

type AuthorizationEntry struct {
	AuthorizationID  uint32            `json:"authorizationId"`
	IDType           AuthorizationType `json:"idType"`
	Name             string            `json:"name"`
	Enabled          bool              `json:"enabled"`
	RemoteAllowed    bool              `json:"remoteAllowed"`
	DateCreated      time.Time         `json:"dateCreated"`
	DateLastActive   time.Time         `json:"dateLastActive"`
	LockCount        uint16            `json:"lockCount"`
	TimeLimited      bool              `json:"timeLimited"`
	AllowedFromDate  time.Time         `json:"allowedFromDate,omitempty"`
	AllowedUntilDate time.Time         `json:"allowedUntilDate,omitempty"`
	AllowedWeekdays  byte              `json:"allowedWeekdays,omitempty"`
	AllowedFromTime  [2]byte           `json:"allowedFromTime,omitempty"`
	AllowedUntilTime [2]byte           `json:"allowedUntilTime,omitempty"`
}

func (c *AuthorizationEntry) GetCommandCode() CommandCode { return CommandAuthorizationEntry }
func (c *AuthorizationEntry) FromMessage(b []byte) error {
	if len(b) < 56 {
		return fmt.Errorf("authorization entry too short: %d bytes, need at least 56", len(b))
	}
	c.AuthorizationID = binary.LittleEndian.Uint32(b[0:4])
	c.IDType = AuthorizationType(b[4])
	c.Name = string(bytes.Trim(b[5:37], "\x00"))
	c.Enabled = b[37] != 0
	c.RemoteAllowed = b[38] != 0
	c.DateCreated = fromNukiTime(b[39:46], time.UTC)
	c.DateLastActive = fromNukiTime(b[46:53], time.UTC)
	c.LockCount = binary.LittleEndian.Uint16(b[53:55])
	c.TimeLimited = b[55] != 0
	if c.TimeLimited && len(b) >= 75 {
		c.AllowedFromDate = fromNukiTime(b[56:63], time.UTC)
		c.AllowedUntilDate = fromNukiTime(b[63:70], time.UTC)
		c.AllowedWeekdays = b[70]
		c.AllowedFromTime = [2]byte{b[71], b[72]}
		c.AllowedUntilTime = [2]byte{b[73], b[74]}
	}
	return nil
}

// RequestAuthorizationEntries (0x0009)

var _ Request = &RequestAuthorizationEntries{}

type RequestAuthorizationEntries struct {
	Offset      uint16
	Count       uint16
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestAuthorizationEntries) GetCommandCode() CommandCode {
	return CommandRequestAuthorizationEntries
}
func (c *RequestAuthorizationEntries) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint16(nil, c.Offset),
		binary.LittleEndian.AppendUint16(nil, c.Count),
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// RemoveAuthorizationEntry (0x0008)

var _ Request = &RemoveAuthorizationEntry{}

type RemoveAuthorizationEntry struct {
	AuthorizationID uint32
	Nonce           []byte
	SecurityPin     Pin
}

func (c *RemoveAuthorizationEntry) GetCommandCode() CommandCode {
	return CommandRemoveAuthorizationEntry
}
func (c *RemoveAuthorizationEntry) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint32(nil, c.AuthorizationID),
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// UpdateAuthorizationEntry (0x0025)
// All time fields are always present in the payload (not conditional on TimeLimited).

var _ Request = &UpdateAuthorizationEntry{}

type UpdateAuthorizationEntry struct {
	AuthorizationID  uint32
	Name             string
	Enabled          bool
	RemoteAllowed    bool
	TimeLimited      bool
	AllowedFromDate  time.Time
	AllowedUntilDate time.Time
	AllowedWeekdays  byte
	AllowedFromTime  [2]byte
	AllowedUntilTime [2]byte
	Nonce            []byte
	SecurityPin      Pin
}

func (c *UpdateAuthorizationEntry) GetCommandCode() CommandCode {
	return CommandUpdateAuthorizationEntry
}
func (c *UpdateAuthorizationEntry) GetPayload() []byte {
	name := [32]byte{}
	copy(name[:], c.Name)
	return slices.Concat(
		binary.LittleEndian.AppendUint32(nil, c.AuthorizationID),
		name[:],
		[]byte{boolToByte(c.Enabled), boolToByte(c.RemoteAllowed), boolToByte(c.TimeLimited)},
		toNukiTime(c.AllowedFromDate),
		toNukiTime(c.AllowedUntilDate),
		[]byte{c.AllowedWeekdays},
		c.AllowedFromTime[:],
		c.AllowedUntilTime[:],
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}
