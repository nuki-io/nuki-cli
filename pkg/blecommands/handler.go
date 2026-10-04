package blecommands

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"slices"
)

var authId5GPairing = []byte{0x7F, 0xFF, 0xFF, 0xFF}

type BleHandler struct {
	crypto   Crypto
	authId   []byte
	recorder func(cmd CommandCode, payload []byte)
}

func NewBleHandler(crypto Crypto, authId []byte) *BleHandler {
	return &BleHandler{
		crypto: crypto,
		authId: authId,
	}
}

// SetRecorder registers fn to receive every CRC-valid response payload before it is parsed.
func (h *BleHandler) SetRecorder(fn func(cmd CommandCode, payload []byte)) {
	h.recorder = fn
}

func (h *BleHandler) record(cmd CommandCode, payload []byte) {
	if h.recorder != nil {
		h.recorder(cmd, slices.Clone(payload))
	}
}

func (h *BleHandler) ToMessage(c Request) []byte {
	payload := c.GetPayload()
	res := make([]byte, 2+len(payload))
	binary.LittleEndian.PutUint16(res, uint16(c.GetCommandCode()))
	for i, x := range payload {
		res[i+2] = x
	}
	res = binary.LittleEndian.AppendUint16(res, CRC(res))
	return res
}

func (h *BleHandler) ToEncryptedMessage(c Request, nonce []byte) []byte {
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
	cmdCode, payload, err := splitMessage(b, 0)
	if err != nil {
		return nil, err
	}
	h.record(cmdCode, payload)
	return ParseResponse(cmdCode, payload)
}

func (h *BleHandler) FromEncryptedDeviceResponse(b []byte) (Response, error) {
	cmdCode, payload, err := h.DecryptDeviceResponse(b)
	if err != nil {
		return nil, err
	}
	return ParseResponse(cmdCode, payload)
}

// DecryptDeviceResponse decrypts and CRC-checks b without parsing the payload,
// so it also works for command codes this package has no type for.
func (h *BleHandler) DecryptDeviceResponse(b []byte) (CommandCode, []byte, error) {
	// nonce[24] + authId[4] + length[2] + ciphertext
	if len(b) < 30 {
		return 0, nil, fmt.Errorf("invalid encrypted response length: %d. must be at least 30 bytes", len(b))
	}
	nonce := b[0:24]
	authId := b[24:28]
	if !slices.Equal(h.authId, authId5GPairing) && !slices.Equal(authId, h.authId) {
		return 0, nil, fmt.Errorf("authId mismatch: expected %x, got %x", h.authId, authId)
	}

	pdata, err := h.crypto.Decrypt(nonce, b[30:])
	if err != nil {
		return 0, nil, fmt.Errorf("failed to decrypt response: %w", err)
	}
	// pdata = authId[4] + command[2] + payload + crc[2]
	cmdCode, payload, err := splitMessage(pdata, 4)
	if err != nil {
		return 0, nil, err
	}
	h.record(cmdCode, payload)
	return cmdCode, payload, nil
}

// splitMessage checks the trailing CRC of b, which covers everything before it,
// and returns the command code at offset and the payload between it and the CRC.
func splitMessage(b []byte, offset int) (CommandCode, []byte, error) {
	if len(b) < offset+4 {
		return 0, nil, fmt.Errorf("invalid response length: %d. must be at least %d bytes", len(b), offset+4)
	}
	crcExpect := CRC(b[:len(b)-2])
	crcReceived := binary.LittleEndian.Uint16(b[len(b)-2:])
	cmdCode := CommandCode(binary.LittleEndian.Uint16(b[offset : offset+2]))
	payload := b[offset+2 : len(b)-2]

	slog.Debug(
		"Received response from smartlock",
		"cmd", cmdCode.String(),
		"payload", fmt.Sprintf("%x", payload),
		"crcReceived", fmt.Sprintf("%x", crcReceived),
		"crcExpect", fmt.Sprintf("%x", crcExpect))

	if crcReceived != crcExpect {
		return 0, nil, fmt.Errorf("CRC mismatch: expected %x, got %x", crcExpect, crcReceived)
	}
	return cmdCode, payload, nil
}

// ParseResponse decodes a decrypted, CRC-checked payload into its typed response.
// An ErrorReport is returned together with a non-nil error.
func ParseResponse(cmdCode CommandCode, payload []byte) (Response, error) {
	cmdImpl, ok := responseImplMap[cmdCode]
	if !ok {
		return nil, fmt.Errorf("unhandled response command code: %x, name: %s", int(cmdCode), cmdCode)
	}
	cmd := cmdImpl()
	err := cmd.FromMessage(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command: %w", err)
	}
	if e, ok := cmd.(*ErrorReport); ok {
		return cmd, fmt.Errorf("%s, command: %s", e.Error, e.CommandIdentifier)
	}
	return cmd, nil
}
