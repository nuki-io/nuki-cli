package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetTimeControlEntries(ctx context.Context) ([]blecommands.TimeControlEntry, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	req := &blecommands.RequestTimeControlEntries{
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	var entries []blecommands.TimeControlEntry
	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt time control response: %w", err)
			}
			slog.Debug("Received time control response", "cmd", res.GetCommandCode())
			switch r := res.(type) {
			case *blecommands.TimeControlEntry:
				entries = append(entries, *r)
			case *blecommands.Status:
				if r.Status == blecommands.StatusComplete {
					return entries, nil
				}
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
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
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return 0, fmt.Errorf("failed to decrypt add time control response: %w", err)
			}
			slog.Debug("Received add time control response", "cmd", res.GetCommandCode())
			switch r := res.(type) {
			case *blecommands.TimeControlEntryID:
				return r.EntryID, nil
			case *blecommands.Status:
				if r.Status == blecommands.StatusComplete {
					return 0, nil
				}
			}
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
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
