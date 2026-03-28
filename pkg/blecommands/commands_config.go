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
	NukiID           uint32    `json:"nukiId"`
	Name             string    `json:"name"`
	Latitude         float32   `json:"latitude"`
	Longitude        float32   `json:"longitude"`
	AutoUnlatch      bool      `json:"autoUnlatch"`
	PairingEnabled   bool      `json:"pairingEnabled"`
	ButtonEnabled    bool      `json:"buttonEnabled"`
	LedEnabled       bool      `json:"ledEnabled"`
	LedBrightness    uint8     `json:"ledBrightness"`
	CurrentTime      time.Time `json:"currentTime"`
	TimezoneOffset   int16     `json:"timezoneOffset"`
	DstMode          uint8     `json:"dstMode"`
	HasFob           bool      `json:"hasFob"`
	FobAction1       uint8     `json:"fobAction1"`
	FobAction2       uint8     `json:"fobAction2"`
	FobAction3       uint8     `json:"fobAction3"`
	SingleLock       bool      `json:"singleLock"`
	AdvertisingMode  uint8     `json:"advertisingMode"`
	HasKeypad        bool      `json:"hasKeypad"`
	FirmwareVersion  string    `json:"firmwareVersion"`
	HardwareRevision string    `json:"hardwareRevision"`
	HomeKitStatus    uint8     `json:"homeKitStatus"`
	TimezoneID       uint16    `json:"timezoneId"`
	DeviceType       uint8     `json:"deviceType"`
	Capabilities     uint8     `json:"capabilities"`
	HasKeypad2       bool      `json:"hasKeypad2"`
	MatterStatus     uint8     `json:"matterStatus"`
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

func (c *RequestConfig) GetCommandCode() CommandCode {
	return CommandRequestConfig
}

func (c *RequestConfig) GetPayload() []byte {
	return c.Nonce
}
