package blecommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

var _ Response = &Config{}

// Config Command 0x0015
type Config struct {
	NukiID           uint32
	Name             string
	Latitude         float32
	Longitude        float32
	AutoUnlatch      bool
	PairingEnabled   bool
	ButtonEnabled    bool
	LedEnabled       bool
	LedBrightness    uint8
	CurrentTime      time.Time
	TimezoneOffset   int16
	DstMode          uint8
	HasFob           bool
	FobAction1       uint8
	FobAction2       uint8
	FobAction3       uint8
	SingleLock       bool
	AdvertisingMode  uint8
	HasKeypad        bool
	FirmwareVersion  string
	HardwareRevision string
	HomeKitStatus    uint8
	TimezoneID       uint16
	DeviceType       uint8
	Capabilities     uint8
	HasKeypad2       bool
	MatterStatus     uint8
}

func (c *Config) GetTimezoneLocation() *time.Location {
	tz, err := time.LoadLocation(timezoneMap[c.TimezoneID])
	if err != nil {
		return time.UTC
	}
	return tz
}
func (c *Config) FromMessage(b []byte) error {
	if len(b) < 76 {
		return fmt.Errorf("invalid Config message length")
	}

	c.NukiID = binary.LittleEndian.Uint32(b[0:4])
	c.Name = string(bytes.Trim(b[4:36], "\x00"))

	c.Latitude = math.Float32frombits(binary.LittleEndian.Uint32(b[36:40]))
	c.Longitude = math.Float32frombits(binary.LittleEndian.Uint32(b[40:44]))
	c.AutoUnlatch = byteToBool(b[44])
	c.PairingEnabled = byteToBool(b[45])
	c.ButtonEnabled = byteToBool(b[46])
	c.LedEnabled = byteToBool(b[47])
	c.LedBrightness = b[48]

	c.TimezoneID = binary.LittleEndian.Uint16(b[72:74])
	tz := c.GetTimezoneLocation()
	c.CurrentTime = fromNukiTime(b[49:56], tz)
	c.TimezoneOffset = int16(binary.LittleEndian.Uint16(b[56:58]))
	c.DstMode = b[58]
	c.HasFob = byteToBool(b[59])
	c.FobAction1 = b[60]
	c.FobAction2 = b[61]
	c.FobAction3 = b[62]
	c.SingleLock = byteToBool(b[63])
	c.AdvertisingMode = b[64]
	c.HasKeypad = byteToBool(b[65])

	c.FirmwareVersion = fmt.Sprintf("%d.%d.%d", b[66], b[67], b[68])
	c.HardwareRevision = fmt.Sprintf("%d.%d", b[69], b[70])

	c.HomeKitStatus = b[71]
	// timezoneID is set above, next to the current time
	c.DeviceType = b[74]
	c.Capabilities = b[75]
	if len(b) > 76 { // MatterStatus is optional? TODO: verify
		c.HasKeypad2 = byteToBool(b[76])
	}
	if len(b) > 77 { // MatterStatus is optional? TODO: verify
		c.MatterStatus = b[77]
	}
	return nil
}

func (c *Config) GetCommandCode() CommandCode {
	return CommandConfig
}

var _ Request = &RequestConfig{}

type RequestConfig struct {
	Nonce []byte
}

func (c *RequestConfig) FromMessage(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("invalid RequestConfig message length")
	}
	c.Nonce = b
	return nil
}

func (c *RequestConfig) GetCommandCode() CommandCode {
	return CommandRequestConfig
}

func (c *RequestConfig) GetPayload() []byte {
	return c.Nonce
}
