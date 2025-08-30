package blecommands

import (
	"encoding/binary"
	"strconv"
)

type Pin interface {
	GetPinBytes() []byte
	SetPin(pin string)
}

func NewPin(pin string) Pin {
	if len(pin) == 4 {
		p := &FourDigitPin{}
		p.SetPin(pin)
		return p
	} else if len(pin) == 6 {
		p := &SixDigitPin{}
		p.SetPin(pin)
		return p
	}
	return nil
}

var _ Pin = &FourDigitPin{}

type FourDigitPin struct {
	pin uint16
}

func (f *FourDigitPin) GetPinBytes() []byte {
	return binary.LittleEndian.AppendUint16(nil, f.pin)
}

func (f *FourDigitPin) SetPin(pin string) {
	if p, err := strconv.ParseUint(pin, 10, 16); err == nil {
		f.pin = uint16(p)
	}
}

var _ Pin = &SixDigitPin{}

type SixDigitPin struct {
	pin uint32
}

func (s *SixDigitPin) GetPinBytes() []byte {
	return binary.LittleEndian.AppendUint32(nil, s.pin)
}

func (s *SixDigitPin) SetPin(pin string) {
	if p, err := strconv.ParseUint(pin, 10, 32); err == nil {
		s.pin = uint32(p)
	}
}
