package blecommands

import (
	"slices"
	"time"
)

// SetSecurityPIN (0x0019)

var _ Request = &SetSecurityPIN{}

type SetSecurityPIN struct {
	NewPin      Pin
	Nonce       []byte
	SecurityPin Pin
}

func (c *SetSecurityPIN) GetCommandCode() CommandCode { return CommandSetSecurityPIN }
func (c *SetSecurityPIN) GetPayload() []byte {
	return slices.Concat(c.NewPin.GetPinBytes(), c.Nonce, c.SecurityPin.GetPinBytes())
}

// VerifySecurityPIN (0x0020)

var _ Request = &VerifySecurityPIN{}

type VerifySecurityPIN struct {
	Nonce       []byte
	SecurityPin Pin
}

func (c *VerifySecurityPIN) GetCommandCode() CommandCode { return CommandVerifySecurityPIN }
func (c *VerifySecurityPIN) GetPayload() []byte {
	return slices.Concat(c.Nonce, c.SecurityPin.GetPinBytes())
}

// RequestCalibration (0x001A)

var _ Request = &RequestCalibration{}

type RequestCalibration struct {
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestCalibration) GetCommandCode() CommandCode { return CommandRequestCalibration }
func (c *RequestCalibration) GetPayload() []byte {
	return slices.Concat(c.Nonce, c.SecurityPin.GetPinBytes())
}

// RequestReboot (0x001D)

var _ Request = &RequestReboot{}

type RequestReboot struct {
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestReboot) GetCommandCode() CommandCode { return CommandRequestReboot }
func (c *RequestReboot) GetPayload() []byte {
	return slices.Concat(c.Nonce, c.SecurityPin.GetPinBytes())
}

// UpdateTime (0x0021)

var _ Request = &UpdateTime{}

type UpdateTime struct {
	Time        time.Time
	Nonce       []byte
	SecurityPin Pin
}

func (c *UpdateTime) GetCommandCode() CommandCode { return CommandUpdateTime }
func (c *UpdateTime) GetPayload() []byte {
	return slices.Concat(toNukiTime(c.Time), c.Nonce, c.SecurityPin.GetPinBytes())
}
