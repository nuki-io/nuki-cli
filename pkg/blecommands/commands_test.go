package blecommands_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/stretchr/testify/require"
)

func TestRequestPublicKeyToMessage(t *testing.T) {
	handler := blecommands.NewBleHandler(nil, nil)
	cmd := &blecommands.RequestData{CommandIdentifier: blecommands.CommandPublicKey}
	got := handler.ToMessage(cmd)
	want := []byte{0x01, 0x00, 0x03, 0x00, 0x27, 0xA7}
	require.Equal(t, want, got)
}

func TestEncryptedRequestChallenge(t *testing.T) {
	authId := []byte{0x02, 0x00, 0x00, 0x00}
	crypto := blecommands.NewCrypto([]byte{0x21, 0x7F, 0xCB, 0x0F, 0x18, 0xCA, 0xF2, 0x84, 0xE9, 0xBD, 0xEA, 0x0B, 0x94, 0xB8, 0x3B, 0x8D, 0x10, 0x86, 0x7E, 0xD7, 0x06, 0xBF, 0xDE, 0xDB, 0xD2, 0x38, 0x1F, 0x4C, 0xB3, 0xB8, 0xF7, 0x30})
	handler := blecommands.NewBleHandler(crypto, authId)
	cmd := &blecommands.RequestData{CommandIdentifier: blecommands.CommandChallenge}
	nonce := []byte{0x88, 0xFD, 0xEF, 0xD7, 0xF9, 0x41, 0xB6, 0x3C, 0x24, 0x2B, 0x7F, 0x84, 0xB3, 0xD7, 0x86, 0x88, 0x63, 0x40, 0xA4, 0xA8, 0xB1, 0xC1, 0xEA, 0xA0}
	msg := handler.ToEncryptedMessage(cmd, nonce)

	want := []byte{
		// nonce 24 bytes
		0x88, 0xFD, 0xEF, 0xD7, 0xF9, 0x41, 0xB6, 0x3C,
		0x24, 0x2B, 0x7F, 0x84, 0xB3, 0xD7, 0x86, 0x88,
		0x63, 0x40, 0xA4, 0xA8, 0xB1, 0xC1, 0xEA, 0xA0,
		// auth id 4 bytes
		0x02, 0x00, 0x00, 0x00,
		// msg length 2 bytes
		0x1A, 0x00,
		// PDATA
		0x06, 0x68, 0x19, 0xA2, 0x95, 0x6E, 0x6A, 0x79,
		0xAF, 0x6E, 0xD6, 0x6D, 0x25, 0x7B, 0x27, 0x67,
		0x15, 0xF5, 0x1F, 0x63, 0xA8, 0xBE, 0xB9, 0xED,
		0x0D, 0x47,
	}
	require.Equal(t, want, msg)
}

// SetAdvancedConfig must mirror the length of the device's AdvancedConfig, or the device
// answers ERROR_BAD_LENGTH.
func TestSetAdvancedConfigMirrorsMotorFields(t *testing.T) {
	nonce := make([]byte, 32)
	for _, tc := range []struct {
		name     string
		response []byte
		wantTail []byte
	}{
		{"without motor fields", make([]byte, 31), []byte{}},
		{"with motor fields", append(make([]byte, 31), 0x02, 0x01), []byte{0x02, 0x01}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &blecommands.AdvancedConfig{}
			require.NoError(t, cfg.FromMessage(tc.response))
			req := cfg.ToSetAdvancedConfig()
			req.Nonce = nonce
			req.SecurityPin = blecommands.NewPin("1234")

			// Same fields minus the read-only TotalDegrees, then nonce and 2-byte PIN.
			payload := req.GetPayload()
			require.Len(t, payload, len(tc.response)-2+32+2)
			require.Equal(t, tc.wantTail, payload[29:len(payload)-34], "motor fields")
		})
	}
}

func TestAuthorizationEntryUnsetLastActive(t *testing.T) {
	b := make([]byte, 56)
	copy(b[39:46], []byte{0xEA, 0x07, 0x05, 0x06, 0x0E, 0x36, 0x2A}) // created 2026-05-06 14:54:42, never active
	e := &blecommands.AuthorizationEntry{}
	require.NoError(t, e.FromMessage(b))

	require.Equal(t, time.Date(2026, 5, 6, 14, 54, 42, 0, time.UTC), e.DateCreated)
	require.True(t, e.DateLastActive.IsZero())
	_, err := json.Marshal(e)
	require.NoError(t, err)
}
