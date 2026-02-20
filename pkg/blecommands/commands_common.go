package blecommands

import (
	"encoding/binary"
	"fmt"
)

type Command interface {
	GetCommandCode() CommandCode
}

type Request interface {
	Command
	GetPayload() []byte
}

type Response interface {
	Command
	FromMessage([]byte) error
}

//go:generate stringer -type=StatusCode -trimprefix=Status
type StatusCode uint8

const (
	StatusComplete StatusCode = 0x00
	StatusAccepted StatusCode = 0x01
)

var _ Request = &RequestData{}

type RequestData struct {
	CommandIdentifier CommandCode
	// AdditionalData    []byte
}

func (c *RequestData) GetCommandCode() CommandCode {
	return CommandRequestData
}

func (c *RequestData) GetPayload() []byte {
	payload := make([]byte, 2)
	binary.LittleEndian.PutUint16(payload, uint16(c.CommandIdentifier))
	return payload
}

var _ Response = &Status{}

type Status struct {
	Status StatusCode
}

func (c *Status) GetCommandCode() CommandCode {
	return CommandStatus
}
func (c *Status) FromMessage(b []byte) error {
	if len(b) != 1 {
		return fmt.Errorf("status length must be exactly 1 byte, got: %d", len(b))
	}
	c.Status = StatusCode(b[0])
	return nil
}
func (c *Status) GetPayload() []byte {
	return []byte{byte(c.Status)}
}

var _ Response = &ErrorReport{}

type ErrorReport struct {
	ErrorCode         byte
	Error             string
	CommandIdentifier CommandCode
}

func (c *ErrorReport) GetCommandCode() CommandCode {
	return CommandErrorReport
}
func (c *ErrorReport) FromMessage(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("error report length must be exactly 3 bytes, got: %d", len(b))
	}
	c.ErrorCode = b[0]
	c.Error = errorCodeNames[c.ErrorCode]
	c.CommandIdentifier = CommandCode(binary.LittleEndian.Uint16(b[1:]))
	return nil
}
func (c *ErrorReport) GetPayload() []byte {
	return []byte{byte(c.ErrorCode), byte(c.CommandIdentifier), byte(c.CommandIdentifier >> 8)}
}

var _ Response = &Challenge{}

type Challenge struct {
	Nonce []byte
}

func (c *Challenge) GetCommandCode() CommandCode {
	return CommandChallenge
}
func (c *Challenge) FromMessage(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("challenge length must be exactly 32 bytes, got: %d", len(b))
	}
	c.Nonce = b
	return nil
}
func (c *Challenge) GetPayload() []byte {
	return c.Nonce
}
