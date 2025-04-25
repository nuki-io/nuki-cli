package blecommands

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"slices"
)

type BleHandler struct {
	crypto Crypto
	authId []byte
}

func NewBleHandler(crypto Crypto, authId []byte) *BleHandler {
	return &BleHandler{
		crypto: crypto,
		authId: authId,
	}
}

func (h *BleHandler) ToMessage(c Command) []byte {
	payload := c.GetPayload()
	res := make([]byte, 2+len(payload))
	binary.LittleEndian.PutUint16(res, uint16(c.GetCommandCode()))
	for i, x := range payload {
		res[i+2] = x
	}
	res = binary.LittleEndian.AppendUint16(res, CRC(res))
	return res
}

func (h *BleHandler) ToEncryptedMessage(c Command, nonce []byte) []byte {
	payload := c.GetPayload()
	// length = authId + command + payload length + CRC
	pdata := make([]byte, 0, 4+2+len(payload)+2)
	pdata = append(pdata, h.authId...)
	pdata = binary.LittleEndian.AppendUint16(pdata, uint16(c.GetCommandCode()))
	pdata = append(pdata, payload...)
	pdata = binary.LittleEndian.AppendUint16(pdata, CRC(pdata))

	pdataEnc, _ := h.crypto.Encrypt(nonce, pdata)

	// length = nonce + authId + encrypted message length
	adata := make([]byte, 0, 24+4+2)
	adata = append(adata, nonce...)
	adata = append(adata, h.authId...)
	adata = binary.LittleEndian.AppendUint16(adata, uint16(len(pdataEnc)))
	return slices.Concat(adata, pdataEnc)
}

func (h *BleHandler) FromDeviceResponse(b []byte) (Command, error) {
	if len(b) < 4 {
		return nil, fmt.Errorf("invalid response length: %d. must be at least 4 bytes", len(b))
	}
	crcExpect := CRC(b[:len(b)-2])
	crcReceived := binary.LittleEndian.Uint16(b[len(b)-2:])
	cmdCode := CommandCode(binary.LittleEndian.Uint16(b[0:2]))
	payload := b[2 : len(b)-2]

	slog.Debug(
		"Received response from smartlock",
		"cmd", cmdCode.String(),
		"payload", fmt.Sprintf("%x", payload),
		"crcReceived", fmt.Sprintf("%x", crcReceived),
		"crcExpect", fmt.Sprintf("%x", crcExpect))

	if crcReceived != crcExpect {
		return nil, fmt.Errorf("CRC mismatch: expected %x, got %x", crcExpect, crcReceived)
	}
	cmdImpl, ok := cmdImplMap[cmdCode]
	if !ok {
		return nil, fmt.Errorf("unhandled response command code: %x, name: %s", int(cmdCode), cmdCode)
	}
	cmd := cmdImpl()
	cmd.FromMessage(payload)
	if cmdCode == CommandErrorReport {
		return cmd, fmt.Errorf("error report: %x", payload)
	}
	return cmd, nil
}

func (h *BleHandler) FromEncryptedDeviceResponse(b []byte) (Command, error) {
	nonce := b[0:24]
	authId := b[24:28]
	// msgLen := b[28:30]
	if !slices.Equal(authId, h.authId) {
		return nil, fmt.Errorf("authId mismatch: expected %x, got %x", h.authId, authId)
	}

	pdata, err := h.crypto.Decrypt(nonce, b[30:])
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt response: %w", err)
	}

	// TODO: the code below is mostly the same as the one in FromDeviceResponse but with a different offset for the CRC calculation
	// pdata = authId[4] + command[2] + payload + crc[2]
	crcExpect := CRC(pdata[:len(pdata)-2])
	crcReceived := binary.LittleEndian.Uint16(pdata[len(pdata)-2:])
	cmdCode := CommandCode(binary.LittleEndian.Uint16(pdata[4:6]))
	payload := pdata[6 : len(pdata)-2]

	slog.Debug(
		"Received encrypted response from smartlock",
		"cmd", cmdCode.String(),
		"payload", fmt.Sprintf("%x", payload),
		"crcReceived", fmt.Sprintf("%x", crcReceived),
		"crcExpect", fmt.Sprintf("%x", crcExpect))

	if crcReceived != crcExpect {
		return nil, fmt.Errorf("CRC mismatch: expected %x, got %x", crcExpect, crcReceived)
	}
	cmdImpl, ok := cmdImplMap[cmdCode]
	if !ok {
		return nil, fmt.Errorf("unhandled response command code: %x, name: %s", int(cmdCode), cmdCode)
	}
	cmd := cmdImpl()
	if cmdCode == CommandErrorReport {
		return cmd, fmt.Errorf("error report: %x", payload)
	}
	err = cmd.FromMessage(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}
	return cmd, nil

}
