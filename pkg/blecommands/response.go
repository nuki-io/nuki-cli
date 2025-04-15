package blecommands

import (
	"encoding/binary"
	"fmt"
)

type BleResponse struct {
	cmd      Command
	payload  []byte
	crc      uint16
	crcMatch bool
}

func FromDeviceResponse(b []byte) *BleResponse {

	crcExpect := CRC(b[:len(b)-2])
	crcReceived := binary.LittleEndian.Uint16(b[len(b)-2:])
	fmt.Printf("Without CRC: %x\n", b[:len(b)-2])
	fmt.Printf("Expected CRC: %x\n", crcExpect)
	fmt.Printf("Received CRC: %x\n", crcReceived)

	r := &BleResponse{
		cmd:      Command(binary.LittleEndian.Uint16(b[0:2])),
		payload:  b[2 : len(b)-2],
		crc:      crcReceived,
		crcMatch: crcReceived == crcExpect,
	}
	return r
}
