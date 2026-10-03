//go:build hwtest

package hwtest

import (
	"context"
	"math/rand/v2"
	"slices"
	"testing"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testWrite(t *testing.T) {
	t.Run("SetConfig", testSetConfig)
	t.Run("SetAdvancedConfig", testSetAdvancedConfig)
	t.Run("UpdateTime", testUpdateTime)
	t.Run("EnableLogging", testEnableLogging)
	t.Run("TimeControlEntry", testTimeControlEntry)
	t.Run("KeypadCode", testKeypadCode)
}

func testSetConfig(t *testing.T) {
	ctx, flow := newFlow(t)
	orig, err := flow.GetConfig(ctx)
	require.NoError(t, err)

	want := orig.ToSetConfig()
	want.LedBrightness = orig.LedBrightness%5 + 1
	t.Cleanup(func() {
		ctx := cleanupCtx(t)
		if !assert.NoError(t, flow.SetConfig(ctx, orig.ToSetConfig()), "restoring config") {
			return
		}
		got, err := flow.GetConfig(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, orig.ToSetConfig(), got.ToSetConfig(), "config not restored")
		}
	})

	expected := *want
	require.NoError(t, flow.SetConfig(ctx, want))
	got, err := flow.GetConfig(ctx)
	require.NoError(t, err)
	require.Equal(t, &expected, got.ToSetConfig())
}

func testSetAdvancedConfig(t *testing.T) {
	ctx, flow := newFlow(t)
	orig, err := flow.GetAdvancedConfig(ctx)
	require.NoError(t, err)

	want := orig.ToSetAdvancedConfig()
	if orig.LockNGoTimeout <= 55 {
		want.LockNGoTimeout = orig.LockNGoTimeout + 5
	} else {
		want.LockNGoTimeout = orig.LockNGoTimeout - 5
	}
	t.Cleanup(func() {
		ctx := cleanupCtx(t)
		if !assert.NoError(t, flow.SetAdvancedConfig(ctx, orig.ToSetAdvancedConfig()), "restoring advanced config") {
			return
		}
		got, err := flow.GetAdvancedConfig(ctx)
		if assert.NoError(t, err) {
			assert.Equal(t, orig, got, "advanced config not restored")
		}
	})

	expected := *want
	require.NoError(t, flow.SetAdvancedConfig(ctx, want))
	got, err := flow.GetAdvancedConfig(ctx)
	require.NoError(t, err)
	require.Equal(t, &expected, got.ToSetAdvancedConfig())
	// Not part of SetAdvancedConfig, so a write must leave them alone.
	require.Equal(t, orig.TotalDegrees, got.TotalDegrees)
	require.Equal(t, orig.MotorSpeed, got.MotorSpeed)
	require.Equal(t, orig.EnableSlowSpeedDuringNightMode, got.EnableSlowSpeedDuringNightMode)
}

func testUpdateTime(t *testing.T) {
	ctx, flow := newFlow(t)
	require.NoError(t, flow.UpdateTime(ctx, time.Now().UTC()))
	s, err := flow.GetStatus(ctx)
	require.NoError(t, err)
	require.WithinDuration(t, time.Now(), s.CurrentTime, 10*time.Second)
}

// testEnableLogging only re-applies the current setting: disabling logging on the
// user's lock is not worth the risk of losing entries.
func testEnableLogging(t *testing.T) {
	ctx, flow := newFlow(t)
	_, before, err := flow.GetLogs(ctx, 0, 1, true)
	require.NoError(t, err)
	require.NotNil(t, before)

	require.NoError(t, flow.EnableLogging(ctx, before.LoggingEnabled))
	_, after, err := flow.GetLogs(ctx, 0, 1, true)
	require.NoError(t, err)
	require.Equal(t, before.LoggingEnabled, after.LoggingEnabled)
}

func testTimeControlEntry(t *testing.T) {
	ctx, flow := newFlow(t)
	before, err := flow.GetTimeControlEntries(ctx)
	require.NoError(t, err)

	// Lock, never unlock, in case cleanup fails and the entry stays on the device.
	const weekdays, hour, minute = 0x01, 3, 33
	_, err = flow.AddTimeControlEntry(ctx, weekdays, hour, minute, blecommands.Lock)
	require.NoError(t, err)

	added := findNewTimeControlEntry(t, ctx, flow, before)
	removed := false
	t.Cleanup(func() {
		if !removed {
			assert.NoError(t, flow.RemoveTimeControlEntry(cleanupCtx(t), added.EntryID), "removing test entry %d", added.EntryID)
		}
	})
	require.True(t, added.Enabled)
	require.Equal(t, byte(weekdays), added.Weekdays)
	require.Equal(t, byte(hour), added.TimeHour)
	require.Equal(t, byte(minute), added.TimeMinute)
	require.Equal(t, blecommands.Lock, added.LockAction)

	require.NoError(t, flow.UpdateTimeControlEntry(ctx, added.EntryID, false, weekdays, hour, minute+1, blecommands.Lock))
	entries, err := flow.GetTimeControlEntries(ctx)
	require.NoError(t, err)
	i := slices.IndexFunc(entries, func(e blecommands.TimeControlEntry) bool { return e.EntryID == added.EntryID })
	require.GreaterOrEqual(t, i, 0, "updated entry missing")
	require.False(t, entries[i].Enabled)
	require.Equal(t, byte(minute+1), entries[i].TimeMinute)

	require.NoError(t, flow.RemoveTimeControlEntry(ctx, added.EntryID))
	removed = true
	entries, err = flow.GetTimeControlEntries(ctx)
	require.NoError(t, err)
	require.Len(t, entries, len(before))
}

// findNewTimeControlEntry diffs against before, since the flow does not always return the new ID.
func findNewTimeControlEntry(t *testing.T, ctx context.Context, flow *bleflows.Flow, before []blecommands.TimeControlEntry) blecommands.TimeControlEntry {
	t.Helper()
	after, err := flow.GetTimeControlEntries(ctx)
	require.NoError(t, err)
	for _, e := range after {
		if !slices.ContainsFunc(before, func(b blecommands.TimeControlEntry) bool { return b.EntryID == e.EntryID }) {
			return e
		}
	}
	require.FailNow(t, "added time control entry not found")
	return blecommands.TimeControlEntry{}
}

// randomKeypadCode avoids a fixed, publicly known code staying on the lock if cleanup fails.
// Keypad codes are 6 digits without 0 and must not start with 12.
func randomKeypadCode() uint32 {
	code := uint32(rand.IntN(8) + 2)
	for range 5 {
		code = code*10 + uint32(rand.IntN(9)+1)
	}
	return code
}

func testKeypadCode(t *testing.T) {
	requireKeypad(t)
	ctx, flow := newFlow(t)
	before, err := flow.GetKeypadCodes(ctx, 0, 100)
	require.NoError(t, err)

	const name = "hwtest"
	code := randomKeypadCode()
	_, err = flow.AddKeypadCode(ctx, code, name)
	require.NoError(t, err)

	after, err := flow.GetKeypadCodes(ctx, 0, 100)
	require.NoError(t, err)
	i := slices.IndexFunc(after, func(c blecommands.KeypadCode) bool {
		return !slices.ContainsFunc(before, func(b blecommands.KeypadCode) bool { return b.CodeID == c.CodeID })
	})
	require.GreaterOrEqual(t, i, 0, "added keypad code not found")
	added := after[i]
	removed := false
	t.Cleanup(func() {
		if !removed {
			assert.NoError(t, flow.RemoveKeypadCode(cleanupCtx(t), added.CodeID), "removing test code %d", added.CodeID)
		}
	})
	require.Equal(t, name, added.Name)
	require.True(t, added.Code == code, "stored code differs from the one sent")
	require.True(t, added.Enabled)

	require.NoError(t, flow.UpdateKeypadCode(ctx, added.CodeID, code, name, false))
	after, err = flow.GetKeypadCodes(ctx, 0, 100)
	require.NoError(t, err)
	i = slices.IndexFunc(after, func(c blecommands.KeypadCode) bool { return c.CodeID == added.CodeID })
	require.GreaterOrEqual(t, i, 0, "updated keypad code missing")
	require.False(t, after[i].Enabled)

	require.NoError(t, flow.RemoveKeypadCode(ctx, added.CodeID))
	removed = true
	after, err = flow.GetKeypadCodes(ctx, 0, 100)
	require.NoError(t, err)
	require.Len(t, after, len(before))
}
