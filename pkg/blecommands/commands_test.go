package blecommands_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.nuki.io/nuki/nukictl/pkg/blecommands"
)

func TestRequestPublicKeyToMessage(t *testing.T) {
	cmd := blecommands.NewUnencryptedRequestData(blecommands.PublicKey)
	got := cmd.ToMessage()
	want := []byte{0x01, 0x00, 0x03, 0x00, 0x27, 0xA7}
	require.Equal(t, want, got)
}
