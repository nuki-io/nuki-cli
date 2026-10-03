//go:build hwtest

package hwtest

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDisruptive runs at any level, but each test needs its own NUKI_TEST_ALLOW entry.
func testDisruptive(t *testing.T) {
	t.Run("Reboot", testReboot)
	t.Run("SetSecurityPIN", testSetSecurityPIN)
}

func testReboot(t *testing.T) {
	requireAllowed(t, "reboot")
	ctx, flow := newFlow(t)
	require.NoError(t, flow.Reboot(ctx))
	flow.DisconnectDevice()

	time.Sleep(15 * time.Second)
	ctx, flow = newFlow(t)
	_, err := flow.GetStatus(ctx)
	require.NoError(t, err)
}

// testSetSecurityPIN switches to NUKI_TEST_TEMP_PIN and back. If restoring fails,
// the lock keeps the temporary PIN. Assertions avoid printing either PIN.
func testSetSecurityPIN(t *testing.T) {
	requireAllowed(t, "set-pin")
	tempPin := os.Getenv("NUKI_TEST_TEMP_PIN")
	auth, err := store.Load(deviceID)
	require.NoError(t, err)
	origPin := auth.Pin
	_, numErr := strconv.ParseUint(tempPin, 10, 32)
	require.True(t, numErr == nil && len(tempPin) == len(origPin) && tempPin != origPin,
		"NUKI_TEST_TEMP_PIN must be numeric, differ from the current PIN and have the same number of digits")

	ctx, flow := newFlow(t)
	require.NoError(t, flow.SetSecurityPIN(ctx, tempPin))
	t.Cleanup(func() {
		ctx := cleanupCtx(t)
		if assert.NoError(t, flow.SetSecurityPIN(ctx, origPin), "restoring PIN failed, the lock still uses NUKI_TEST_TEMP_PIN") {
			assert.NoError(t, flow.VerifyPIN(ctx))
		}
	})
	require.NoError(t, flow.VerifyPIN(ctx))
}
