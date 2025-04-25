package blecommands

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"slices"
)

type bleHandler struct {
	crypto Crypto
	authId []byte
}

func NewBleHandler(crypto Crypto, authId []byte) bleHandler {
	return bleHandler{
		crypto: crypto,
		authId: authId,
	}
}

func (h *bleHandler) ToMessage(c Command) []byte {
	payload := c.GetPayload()
	res := make([]byte, 2+len(payload))
	binary.LittleEndian.PutUint16(res, uint16(c.GetCommandCode()))
	for i, x := range payload {
		res[i+2] = x
	}
	res = binary.LittleEndian.AppendUint16(res, CRC(res))
	return res
}

func (h *bleHandler) ToEncryptedMessage(c Command, nonce []byte) []byte {
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

func (h *bleHandler) FromDeviceResponse(b []byte) (Command, error) {
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
