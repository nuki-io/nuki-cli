package blecommands_test

import (
	"bytes"
	"testing"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/stretchr/testify/require"
)

var (
	testAuthId = []byte{0x02, 0x00, 0x00, 0x00}
	testKey    = bytes.Repeat([]byte{0x42}, 32)
	testNonce  = bytes.Repeat([]byte{0x07}, 24)
)

func newTestHandler() *blecommands.BleHandler {
	return blecommands.NewBleHandler(blecommands.NewCrypto(testKey), testAuthId)
}

// Device responses share the layout of encrypted requests, so ToEncryptedMessage builds them.
func TestDecryptDeviceResponseUnknownCode(t *testing.T) {
	h := newTestHandler()
	msg := h.ToEncryptedMessage(&blecommands.RawRequest{Code: 0x7A01, Payload: []byte{0xDE, 0xAD}}, testNonce)

	code, payload, err := h.DecryptDeviceResponse(msg)
	require.NoError(t, err)
	require.Equal(t, blecommands.CommandCode(0x7A01), code)
	require.Equal(t, []byte{0xDE, 0xAD}, payload)

	_, err = h.FromEncryptedDeviceResponse(msg)
	require.ErrorContains(t, err, "unhandled response command code")
}

func TestFromEncryptedDeviceResponseParsesKnownCode(t *testing.T) {
	h := newTestHandler()
	msg := h.ToEncryptedMessage(&blecommands.RawRequest{Code: blecommands.CommandStatus, Payload: []byte{0x00}}, testNonce)

	res, err := h.FromEncryptedDeviceResponse(msg)
	require.NoError(t, err)
	require.Equal(t, &blecommands.Status{Status: blecommands.StatusComplete}, res)
}

func TestDecryptDeviceResponseInvalid(t *testing.T) {
	h := newTestHandler()
	valid := h.ToEncryptedMessage(&blecommands.RawRequest{Code: blecommands.CommandStatus, Payload: []byte{0x00}}, testNonce)

	// Encrypts pdata as the device would, for plaintexts too short to hold authId, code and CRC.
	sealed := func(pdata []byte) []byte {
		enc, err := blecommands.NewCrypto(testKey).Encrypt(testNonce, pdata)
		require.NoError(t, err)
		return append(append(append([]byte{}, testNonce...), testAuthId...), append([]byte{0, 0}, enc...)...)
	}
	tampered := bytes.Clone(valid)
	tampered[len(tampered)-1] ^= 0xFF

	tests := []struct {
		name string
		msg  []byte
		want string
	}{
		{name: "empty", msg: nil, want: "must be at least 30 bytes"},
		{name: "header only", msg: valid[:29], want: "must be at least 30 bytes"},
		{name: "no ciphertext", msg: valid[:30], want: "failed to decrypt"},
		{name: "short plaintext", msg: sealed([]byte{0x02, 0x00, 0x00, 0x00, 0x0E}), want: "must be at least 8 bytes"},
		{name: "tampered ciphertext", msg: tampered, want: "failed to decrypt"},
		{name: "bad CRC", msg: sealed([]byte{0x02, 0x00, 0x00, 0x00, 0x0E, 0x00, 0x00, 0x00, 0x00}), want: "CRC mismatch"},
		{name: "wrong authId", msg: append(append(bytes.Clone(valid[:24]), 0x09, 0, 0, 0), valid[28:]...), want: "authId mismatch"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := h.DecryptDeviceResponse(tt.msg)
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestFromDeviceResponseInvalid(t *testing.T) {
	h := blecommands.NewBleHandler(nil, nil)
	valid := h.ToMessage(&blecommands.RawRequest{Code: blecommands.CommandStatus, Payload: []byte{0x00}})
	badCRC := bytes.Clone(valid)
	badCRC[len(badCRC)-1] ^= 0xFF

	_, err := h.FromDeviceResponse(valid[:3])
	require.ErrorContains(t, err, "must be at least 4 bytes")
	_, err = h.FromDeviceResponse(badCRC)
	require.ErrorContains(t, err, "CRC mismatch")
	res, err := h.FromDeviceResponse(valid)
	require.NoError(t, err)
	require.Equal(t, &blecommands.Status{Status: blecommands.StatusComplete}, res)
}
