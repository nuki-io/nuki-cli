package bleflows

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetTimeControlEntries(ctx context.Context) ([]blecommands.TimeControlEntry, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	return collectResponses[blecommands.TimeControlEntry](ctx, f, &blecommands.RequestTimeControlEntries{
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) AddTimeControlEntry(ctx context.Context, weekdays, hour, minute byte, action blecommands.Action) (byte, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get challenge: %w", err)
	}
	req := &blecommands.AddTimeControlEntry{
		Weekdays:    weekdays,
		TimeHour:    hour,
		TimeMinute:  minute,
		LockAction:  action,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	var entryID byte
	// Unlike AddKeypadCode, a StatusComplete follows the ID and must be consumed,
	// or the next command reads it as its own response.
	err = f.exchange(ctx, req, func(res blecommands.Response) bool {
		if r, ok := res.(*blecommands.TimeControlEntryID); ok {
			entryID = r.EntryID
		}
		return false
	})
	return entryID, err
}

func (f *Flow) RemoveTimeControlEntry(ctx context.Context, entryID byte) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.RemoveTimeControlEntry{
		EntryID:     entryID,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) UpdateTimeControlEntry(ctx context.Context, entryID byte, enabled bool, weekdays, hour, minute byte, action blecommands.Action) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.UpdateTimeControlEntry{
		EntryID:     entryID,
		Enabled:     enabled,
		Weekdays:    weekdays,
		TimeHour:    hour,
		TimeMinute:  minute,
		LockAction:  action,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}
