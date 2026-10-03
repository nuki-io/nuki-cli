package blecommands

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
)

// RequestAdvancedConfig (0x0036)

var _ Request = &RequestAdvancedConfig{}

type RequestAdvancedConfig struct {
	Nonce []byte
}

func (c *RequestAdvancedConfig) GetCommandCode() CommandCode { return CommandRequestAdvancedConfig }
func (c *RequestAdvancedConfig) GetPayload() []byte          { return c.Nonce }

// AdvancedConfig (0x0037)
// Byte layout:
//   TotalDegrees(2) + 4×sint16(8) + LockNGoTimeout(1) + SingleButtonAction(1) +
//   DoubleButtonAction(1) + DetachedCylinder(1) + BatteryType(1) + AutoBattDet(1) +
//   UnlatchDuration(1) + AutoLockTimeout(2) + AutoUnlockDisabled(1) + NightmodeEnabled(1) +
//   NightmodeStartTime(2) + NightmodeEndTime(2) + NightmodeAutoLock(1) +
//   NightmodeAutoUnlockDisabled(1) + NightmodeImmediateLock(1) + AutoLockEnabled(1) +
//   ImmediateAutoLockEnabled(1) + AutoUpdateEnabled(1) = 31 bytes base
// + MotorSpeed(1) + EnableSlowSpeedNightMode(1) = Ultra only

var _ Response = &AdvancedConfig{}

type AdvancedConfig struct {
	TotalDegrees                      uint16  `json:"totalDegrees"`
	UnlockedPositionOffsetDegrees     int16   `json:"unlockedPositionOffsetDegrees"`
	LockedPositionOffsetDegrees       int16   `json:"lockedPositionOffsetDegrees"`
	SingleLockedPositionOffsetDegrees int16   `json:"singleLockedPositionOffsetDegrees"`
	UnlockedToLockedTransitionOffset  int16   `json:"unlockedToLockedTransitionOffset"`
	LockNGoTimeout                    uint8   `json:"lockNGoTimeout"`
	SingleButtonPressAction           uint8   `json:"singleButtonPressAction"`
	DoubleButtonPressAction           uint8   `json:"doubleButtonPressAction"`
	DetachedCylinder                  bool    `json:"detachedCylinder"`
	BatteryType                       uint8   `json:"batteryType"`
	AutomaticBatteryTypeDetection     bool    `json:"automaticBatteryTypeDetection"`
	UnlatchDuration                   uint8   `json:"unlatchDuration"`
	AutoLockTimeout                   uint16  `json:"autoLockTimeout"`
	AutoUnlockDisabled                bool    `json:"autoUnlockDisabled"`
	NightmodeEnabled                  bool    `json:"nightmodeEnabled"`
	NightmodeStartTime                [2]byte `json:"nightmodeStartTime"`
	NightmodeEndTime                  [2]byte `json:"nightmodeEndTime"`
	NightmodeAutoLockEnabled          bool    `json:"nightmodeAutoLockEnabled"`
	NightmodeAutoUnlockDisabled       bool    `json:"nightmodeAutoUnlockDisabled"`
	NightmodeImmediateLockOnStart     bool    `json:"nightmodeImmediateLockOnStart"`
	AutoLockEnabled                   bool    `json:"autoLockEnabled"`
	ImmediateAutoLockEnabled          bool    `json:"immediateAutoLockEnabled"`
	AutoUpdateEnabled                 bool    `json:"autoUpdateEnabled"`
	// Smart Lock Ultra only
	MotorSpeed                     uint8 `json:"motorSpeed,omitempty"`
	EnableSlowSpeedDuringNightMode bool  `json:"enableSlowSpeedDuringNightMode,omitempty"`
}

func (c *AdvancedConfig) GetCommandCode() CommandCode { return CommandAdvancedConfig }
func (c *AdvancedConfig) FromMessage(b []byte) error {
	if len(b) < 25 {
		return fmt.Errorf("advanced config too short: %d bytes, need at least 25", len(b))
	}
	c.TotalDegrees = binary.LittleEndian.Uint16(b[0:2])
	c.UnlockedPositionOffsetDegrees = int16(binary.LittleEndian.Uint16(b[2:4]))
	c.LockedPositionOffsetDegrees = int16(binary.LittleEndian.Uint16(b[4:6]))
	c.SingleLockedPositionOffsetDegrees = int16(binary.LittleEndian.Uint16(b[6:8]))
	c.UnlockedToLockedTransitionOffset = int16(binary.LittleEndian.Uint16(b[8:10]))
	c.LockNGoTimeout = b[10]
	c.SingleButtonPressAction = b[11]
	c.DoubleButtonPressAction = b[12]
	c.DetachedCylinder = b[13] != 0
	c.BatteryType = b[14]
	c.AutomaticBatteryTypeDetection = b[15] != 0
	c.UnlatchDuration = b[16]
	c.AutoLockTimeout = binary.LittleEndian.Uint16(b[17:19])
	c.AutoUnlockDisabled = b[19] != 0
	c.NightmodeEnabled = b[20] != 0
	c.NightmodeStartTime = [2]byte{b[21], b[22]}
	c.NightmodeEndTime = [2]byte{b[23], b[24]}
	if len(b) > 25 {
		c.NightmodeAutoLockEnabled = b[25] != 0
	}
	if len(b) > 26 {
		c.NightmodeAutoUnlockDisabled = b[26] != 0
	}
	if len(b) > 27 {
		c.NightmodeImmediateLockOnStart = b[27] != 0
	}
	if len(b) > 28 {
		c.AutoLockEnabled = b[28] != 0
	}
	if len(b) > 29 {
		c.ImmediateAutoLockEnabled = b[29] != 0
	}
	if len(b) > 30 {
		c.AutoUpdateEnabled = b[30] != 0
	}
	if len(b) > 31 {
		c.MotorSpeed = b[31]
	}
	if len(b) > 32 {
		c.EnableSlowSpeedDuringNightMode = b[32] != 0
	}
	return nil
}

// SetAdvancedConfig (0x0035)
// Same fields as AdvancedConfig minus TotalDegrees (read-only), plus Nonce and SecurityPIN.

var _ Request = &SetAdvancedConfig{}

type SetAdvancedConfig struct {
	UnlockedPositionOffsetDegrees     int16
	LockedPositionOffsetDegrees       int16
	SingleLockedPositionOffsetDegrees int16
	UnlockedToLockedTransitionOffset  int16
	LockNGoTimeout                    uint8
	SingleButtonPressAction           uint8
	DoubleButtonPressAction           uint8
	DetachedCylinder                  bool
	BatteryType                       uint8
	AutomaticBatteryTypeDetection     bool
	UnlatchDuration                   uint8
	AutoLockTimeout                   uint16
	AutoUnlockDisabled                bool
	NightmodeEnabled                  bool
	NightmodeStartTime                [2]byte
	NightmodeEndTime                  [2]byte
	NightmodeAutoLockEnabled          bool
	NightmodeAutoUnlockDisabled       bool
	NightmodeImmediateLockOnStart     bool
	AutoLockEnabled                   bool
	ImmediateAutoLockEnabled          bool
	AutoUpdateEnabled                 bool
	Nonce                             []byte
	SecurityPin                       Pin
}

// ToSetAdvancedConfig returns a request that writes back c unchanged. Nonce and PIN are left empty.
func (c *AdvancedConfig) ToSetAdvancedConfig() *SetAdvancedConfig {
	return &SetAdvancedConfig{
		UnlockedPositionOffsetDegrees:     c.UnlockedPositionOffsetDegrees,
		LockedPositionOffsetDegrees:       c.LockedPositionOffsetDegrees,
		SingleLockedPositionOffsetDegrees: c.SingleLockedPositionOffsetDegrees,
		UnlockedToLockedTransitionOffset:  c.UnlockedToLockedTransitionOffset,
		LockNGoTimeout:                    c.LockNGoTimeout,
		SingleButtonPressAction:           c.SingleButtonPressAction,
		DoubleButtonPressAction:           c.DoubleButtonPressAction,
		DetachedCylinder:                  c.DetachedCylinder,
		BatteryType:                       c.BatteryType,
		AutomaticBatteryTypeDetection:     c.AutomaticBatteryTypeDetection,
		UnlatchDuration:                   c.UnlatchDuration,
		AutoLockTimeout:                   c.AutoLockTimeout,
		AutoUnlockDisabled:                c.AutoUnlockDisabled,
		NightmodeEnabled:                  c.NightmodeEnabled,
		NightmodeStartTime:                c.NightmodeStartTime,
		NightmodeEndTime:                  c.NightmodeEndTime,
		NightmodeAutoLockEnabled:          c.NightmodeAutoLockEnabled,
		NightmodeAutoUnlockDisabled:       c.NightmodeAutoUnlockDisabled,
		NightmodeImmediateLockOnStart:     c.NightmodeImmediateLockOnStart,
		AutoLockEnabled:                   c.AutoLockEnabled,
		ImmediateAutoLockEnabled:          c.ImmediateAutoLockEnabled,
		AutoUpdateEnabled:                 c.AutoUpdateEnabled,
	}
}

func (c *SetAdvancedConfig) GetCommandCode() CommandCode { return CommandSetAdvancedConfig }
func (c *SetAdvancedConfig) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint16(nil, uint16(c.UnlockedPositionOffsetDegrees)),
		binary.LittleEndian.AppendUint16(nil, uint16(c.LockedPositionOffsetDegrees)),
		binary.LittleEndian.AppendUint16(nil, uint16(c.SingleLockedPositionOffsetDegrees)),
		binary.LittleEndian.AppendUint16(nil, uint16(c.UnlockedToLockedTransitionOffset)),
		[]byte{c.LockNGoTimeout, c.SingleButtonPressAction, c.DoubleButtonPressAction,
			boolToByte(c.DetachedCylinder), c.BatteryType, boolToByte(c.AutomaticBatteryTypeDetection),
			c.UnlatchDuration},
		binary.LittleEndian.AppendUint16(nil, c.AutoLockTimeout),
		[]byte{boolToByte(c.AutoUnlockDisabled), boolToByte(c.NightmodeEnabled)},
		c.NightmodeStartTime[:],
		c.NightmodeEndTime[:],
		[]byte{
			boolToByte(c.NightmodeAutoLockEnabled),
			boolToByte(c.NightmodeAutoUnlockDisabled),
			boolToByte(c.NightmodeImmediateLockOnStart),
			boolToByte(c.AutoLockEnabled),
			boolToByte(c.ImmediateAutoLockEnabled),
			boolToByte(c.AutoUpdateEnabled),
		},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

// SetConfig (0x0013)
// Sends configuration fields back to the device. Mirrors Config (0x0015) minus read-only
// fields (NukiID, firmware, hardware, HomeKit, DeviceType, Capabilities, Keypad flags, MatterStatus).

var _ Request = &SetConfig{}

type SetConfig struct {
	Name            string
	Latitude        float32
	Longitude       float32
	AutoUnlatch     bool
	PairingEnabled  bool
	ButtonEnabled   bool
	LedEnabled      bool
	LedBrightness   uint8
	TimezoneOffset  int16
	DstMode         uint8
	FobAction1      uint8
	FobAction2      uint8
	FobAction3      uint8
	SingleLock      bool
	AdvertisingMode uint8
	TimezoneID      uint16
	Nonce           []byte
	SecurityPin     Pin
}

// ToSetConfig returns a request that writes back c unchanged. Nonce and PIN are left empty.
func (c *Config) ToSetConfig() *SetConfig {
	return &SetConfig{
		Name:            c.Name,
		Latitude:        c.Latitude,
		Longitude:       c.Longitude,
		AutoUnlatch:     c.AutoUnlatch,
		PairingEnabled:  c.PairingEnabled,
		ButtonEnabled:   c.ButtonEnabled,
		LedEnabled:      c.LedEnabled,
		LedBrightness:   c.LedBrightness,
		TimezoneOffset:  c.TimezoneOffset,
		DstMode:         c.DstMode,
		FobAction1:      c.FobAction1,
		FobAction2:      c.FobAction2,
		FobAction3:      c.FobAction3,
		SingleLock:      c.SingleLock,
		AdvertisingMode: c.AdvertisingMode,
		TimezoneID:      c.TimezoneID,
	}
}

func (c *SetConfig) GetCommandCode() CommandCode { return CommandSetConfig }
func (c *SetConfig) GetPayload() []byte {
	name := [32]byte{}
	copy(name[:], c.Name)
	return slices.Concat(
		name[:],
		binary.LittleEndian.AppendUint32(nil, math.Float32bits(c.Latitude)),
		binary.LittleEndian.AppendUint32(nil, math.Float32bits(c.Longitude)),
		[]byte{
			boolToByte(c.AutoUnlatch),
			boolToByte(c.PairingEnabled),
			boolToByte(c.ButtonEnabled),
			boolToByte(c.LedEnabled),
			c.LedBrightness,
		},
		binary.LittleEndian.AppendUint16(nil, uint16(c.TimezoneOffset)),
		[]byte{c.DstMode, c.FobAction1, c.FobAction2, c.FobAction3, boolToByte(c.SingleLock), c.AdvertisingMode},
		binary.LittleEndian.AppendUint16(nil, c.TimezoneID),
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}
