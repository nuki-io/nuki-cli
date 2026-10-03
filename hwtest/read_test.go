//go:build hwtest

package hwtest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func testRead(t *testing.T) {
	t.Run("State", testState)
	t.Run("Config", testConfig)
	t.Run("AdvancedConfig", testAdvancedConfig)
	t.Run("BatteryReport", testBatteryReport)
	t.Run("VerifyPIN", testVerifyPIN)
	t.Run("Logs", testLogs)
	t.Run("AuthorizationEntries", testAuthorizationEntries)
	t.Run("TimeControlEntries", testTimeControlEntries)
	t.Run("KeypadCodes", testKeypadCodes)
}

func testState(t *testing.T) {
	ctx, flow := newFlow(t)
	s, err := flow.GetStatus(ctx)
	require.NoError(t, err)
	logJSON(t, "state", s)

	requireKnown(t, s.NukiState, "nukiState")
	requireKnown(t, s.LockState, "lockState")
	requireKnown(t, s.Trigger, "trigger")
	requireKnown(t, s.LastLockAction, "lastLockAction")
	requireKnown(t, s.LastLockActionTrigger, "lastLockActionTrigger")
	requireKnown(t, s.LastLockActionCompletionStatus, "lastLockActionCompletionStatus")
	require.LessOrEqual(t, s.BatteryPercentage, 100)
	requireRecent(t, s.CurrentTime, "currentTime")
}

func testConfig(t *testing.T) {
	ctx, flow := newFlow(t)
	c, err := flow.GetConfig(ctx)
	require.NoError(t, err)
	logJSON(t, "config", c)

	require.NotZero(t, c.NukiID)
	require.NotEmpty(t, c.Name)
	require.NotEqual(t, "0.0.0", c.FirmwareVersion)
	require.LessOrEqual(t, c.LedBrightness, uint8(5))
	require.InDelta(t, 0, c.Latitude, 90)
	require.InDelta(t, 0, c.Longitude, 180)
	requireRecent(t, c.CurrentTime, "currentTime")
}

func testAdvancedConfig(t *testing.T) {
	ctx, flow := newFlow(t)
	c, err := flow.GetAdvancedConfig(ctx)
	require.NoError(t, err)
	logJSON(t, "advanced config", c)

	require.NotZero(t, c.TotalDegrees)
	require.LessOrEqual(t, c.TotalDegrees, uint16(3600))
	require.GreaterOrEqual(t, c.LockNGoTimeout, uint8(5))
	require.LessOrEqual(t, c.LockNGoTimeout, uint8(60))
	require.LessOrEqual(t, c.UnlatchDuration, uint8(30))
	for _, hm := range [][2]byte{c.NightmodeStartTime, c.NightmodeEndTime} {
		require.Less(t, hm[0], byte(24), "nightmode hour")
		require.Less(t, hm[1], byte(60), "nightmode minute")
	}
}

func testBatteryReport(t *testing.T) {
	ctx, flow := newFlow(t)
	b, err := flow.GetBatteryReport(ctx)
	require.NoError(t, err)
	logJSON(t, "battery report", b)

	// Wide range: 4xAA packs sit around 6 V, the Smart Lock 5 Pro pack around 14 V.
	require.Greater(t, b.BatteryVoltage, uint16(2500), "mV")
	require.Less(t, b.BatteryVoltage, uint16(20000), "mV")
	require.LessOrEqual(t, b.LowestVoltage, b.StartVoltage)
	require.Greater(t, b.StartTemperature, int8(-30))
	require.Less(t, b.StartTemperature, int8(70))
	if b.LockAction != 0 {
		requireKnown(t, b.LockAction, "lockAction")
	}
}

func testVerifyPIN(t *testing.T) {
	ctx, flow := newFlow(t)
	require.NoError(t, flow.VerifyPIN(ctx))
}

func testLogs(t *testing.T) {
	const count = 5
	ctx, flow := newFlow(t)
	entries, total, err := flow.GetLogs(ctx, 0, count, true)
	require.NoError(t, err)
	require.NotNil(t, total, "requested total count but got no LogEntryCount")
	logJSON(t, "log count", total)
	logJSON(t, "log entries", entries)

	require.LessOrEqual(t, len(entries), count)
	if total.Count > 0 {
		require.NotEmpty(t, entries)
	}
	for i, e := range entries {
		requireKnown(t, e.Type, "type")
		require.False(t, e.Time.IsZero(), "entry %d has no time", i)
		require.True(t, e.Time.Before(time.Now().Add(time.Hour)), "entry %d is in the future: %s", i, e.Time)
		if i > 0 {
			require.Less(t, e.Index, entries[i-1].Index, "entries should be sorted descending")
		}
	}
}

func testAuthorizationEntries(t *testing.T) {
	ctx, flow := newFlow(t)
	entries, err := flow.GetAuthorizationEntries(ctx, 0, 100)
	require.NoError(t, err)
	logJSON(t, "authorizations", entries)

	own := false
	for _, e := range entries {
		requireKnown(t, e.IDType, "idType")
		require.NotEmpty(t, e.Name)
		own = own || e.AuthorizationID == ownAuthID
	}
	require.True(t, own, "own authorization %d not in list", ownAuthID)
}

func testTimeControlEntries(t *testing.T) {
	ctx, flow := newFlow(t)
	entries, err := flow.GetTimeControlEntries(ctx)
	require.NoError(t, err)
	logJSON(t, "time control entries", entries)

	for _, e := range entries {
		require.Less(t, e.TimeHour, byte(24))
		require.Less(t, e.TimeMinute, byte(60))
		requireKnown(t, e.LockAction, "lockAction")
	}
}

func testKeypadCodes(t *testing.T) {
	requireKeypad(t)
	ctx, flow := newFlow(t)
	codes, err := flow.GetKeypadCodes(ctx, 0, 100)
	require.NoError(t, err)
	t.Logf("%d keypad codes", len(codes))

	for _, c := range codes {
		require.NotEmpty(t, c.Name)
		require.GreaterOrEqual(t, c.Code, uint32(100000))
		require.LessOrEqual(t, c.Code, uint32(999999))
	}
}
