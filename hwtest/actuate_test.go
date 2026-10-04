//go:build hwtest

package hwtest

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const motorTimeout = 30 * time.Second

type lockOp func(ctx context.Context, flow *bleflows.Flow, action blecommands.Action) error

func lockAction(ctx context.Context, flow *bleflows.Flow, action blecommands.Action) error {
	return flow.PerformLockOperation(ctx, action)
}

func simpleLockAction(ctx context.Context, flow *bleflows.Flow, action blecommands.Action) error {
	return flow.PerformSimpleLockAction(ctx, action)
}

func testActuate(t *testing.T) {
	t.Run("LockAction", func(t *testing.T) { testLockCycle(t, lockAction) })
	t.Run("SimpleLockAction", func(t *testing.T) { testLockCycle(t, simpleLockAction) })
	t.Run("LockNGo", testLockNGo)
	t.Run("Unlatch", testUnlatch)
}

// testLockCycle moves the bolt away from its initial state and back.
func testLockCycle(t *testing.T, op lockOp) {
	ctx, flow, initial := startActuation(t)
	seq := []blecommands.Action{blecommands.Unlock, blecommands.Lock}
	if initial == blecommands.LockStateUnlocked {
		seq = []blecommands.Action{blecommands.Lock, blecommands.Unlock}
	}
	for _, a := range seq {
		require.NoError(t, op(ctx, flow, a), "%s", a)
		want := blecommands.LockStateLocked
		if a == blecommands.Unlock {
			want = blecommands.LockStateUnlocked
		}
		s, err := waitForLockState(ctx, flow, motorTimeout, want)
		require.NoError(t, err)
		require.EqualValues(t, a, s.LastLockAction, "lastLockAction after %s", a)
		require.Equal(t, blecommands.TriggerSystem, s.LastLockActionTrigger)
	}
}

func testLockNGo(t *testing.T) {
	ctx, flow, _ := startActuation(t)
	adv, err := flow.GetAdvancedConfig(ctx)
	require.NoError(t, err)

	// The device reports completion only after the whole cycle: unlock, timeout, lock.
	start := time.Now()
	require.NoError(t, lockAction(ctx, flow, blecommands.LockAndGo))
	require.GreaterOrEqual(t, time.Since(start), time.Duration(adv.LockNGoTimeout)*time.Second,
		"lock 'n' go completed before its timeout")
	s, err := waitForLockState(ctx, flow, motorTimeout, blecommands.LockStateLocked)
	require.NoError(t, err)
	require.Equal(t, blecommands.LockAndGo, s.LastLockAction)
}

// testUnlatch opens the door, so it needs an explicit opt-in on top of the actuate level.
func testUnlatch(t *testing.T) {
	requireAllowed(t, "unlatch")
	ctx, flow, _ := startActuation(t)
	require.NoError(t, lockAction(ctx, flow, blecommands.Unlatch))
	s, err := waitForLockState(ctx, flow, motorTimeout, blecommands.LockStateUnlatched, blecommands.LockStateUnlocked)
	require.NoError(t, err)
	require.EqualValues(t, blecommands.Unlatch, s.LastLockAction)
}

// startActuation skips unless the bolt can safely move, and restores the initial lock state when t ends.
func startActuation(t *testing.T) (context.Context, *bleflows.Flow, blecommands.LockState) {
	t.Helper()
	ctx, flow := newFlow(t)
	s, err := flow.GetStatus(ctx)
	require.NoError(t, err)
	if s.DoorSensorState == blecommands.DoorSensorOpened {
		t.Skip("door sensor reports the door is open")
	}
	initial := s.LockState
	if initial != blecommands.LockStateLocked && initial != blecommands.LockStateUnlocked {
		t.Skipf("lock is %s, need locked or unlocked", initial)
	}

	t.Cleanup(func() {
		ctx := cleanupCtx(t)
		s, err := flow.GetStatus(ctx)
		if !assert.NoError(t, err) || s.LockState == initial {
			return
		}
		action := blecommands.Lock
		if initial == blecommands.LockStateUnlocked {
			action = blecommands.Unlock
		}
		if assert.NoError(t, flow.PerformLockOperation(ctx, action), "restoring lock state %s", initial) {
			_, err := waitForLockState(ctx, flow, motorTimeout, initial)
			assert.NoError(t, err)
		}
	})
	return ctx, flow, initial
}

// waitForLockState polls because a lock action completes before the motor reaches its end position.
func waitForLockState(ctx context.Context, flow *bleflows.Flow, timeout time.Duration, want ...blecommands.LockState) (*blecommands.KeyturnerStates, error) {
	deadline := time.Now().Add(timeout)
	for {
		s, err := flow.GetStatus(ctx)
		if err != nil {
			return nil, err
		}
		if slices.Contains(want, s.LockState) {
			return s, nil
		}
		if s.LockState == blecommands.LockStateMotorBlocked {
			return s, fmt.Errorf("motor blocked while waiting for %v", want)
		}
		if time.Now().After(deadline) {
			return s, fmt.Errorf("lock state is %s after %s, want %v", s.LockState, timeout, want)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
