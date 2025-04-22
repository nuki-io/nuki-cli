package blecommands

import (
	"encoding/binary"
	"fmt"
	"log/slog"
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
	slog.Debug(
		"Received response from smartlock",
		"response", b[:len(b)-2],
		"crcReceived", fmt.Sprintf("%x", crcReceived),
		"crcExpect", fmt.Sprintf("%x", crcExpect))

	r := &BleResponse{
		cmd:      Command(binary.LittleEndian.Uint16(b[0:2])),
		payload:  b[2 : len(b)-2],
		crc:      crcReceived,
		crcMatch: crcReceived == crcExpect,
	}
	return r
}

func (r *BleResponse) GetPayload() []byte {
	return r.payload
}

type BleEncryptedResponse struct {
	authId   []byte
	cmd      Command
	payload  []byte
	crc      uint16
	crcMatch bool
}

func FromEncryptedDeviceResponse(crypto Crypto, b []byte) *BleEncryptedResponse {

	nonce := b[0:24]
	authId := b[24:28]
	// msgLen := b[28:30]

	pdata, _ := crypto.Decrypt(nonce, b[30:])

	crcExpect := CRC(pdata[:len(pdata)-2])
	crcReceived := binary.LittleEndian.Uint16(pdata[len(pdata)-2:])
	slog.Debug(
		"Received response from smartlock",
		"response", pdata[:len(pdata)-2],
		"crcReceived", fmt.Sprintf("%x", crcReceived),
		"crcExpect", fmt.Sprintf("%x", crcExpect))

	r := &BleEncryptedResponse{
		authId:   authId,
		cmd:      Command(binary.LittleEndian.Uint16(pdata[4:6])),
		payload:  pdata[6 : len(pdata)-2],
		crc:      crcReceived,
		crcMatch: crcReceived == crcExpect,
	}
	return r
}

func (r *BleEncryptedResponse) GetCommand() Command {
	return r.cmd
}

func (r *BleEncryptedResponse) GetPayload() []byte {
	return r.payload
}
