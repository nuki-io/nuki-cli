package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) EnableLogging(ctx context.Context, enabled bool) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.EnableLogging{
		Enabled:     enabled,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) GetLogs(ctx context.Context, start int, count int, withCount bool) ([]blecommands.LogEntry, *blecommands.LogEntryCount, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestLogEntries{
		StartIndex:  uint32(start),
		Count:       uint16(count),
		Nonce:       nonce,
		SortOrder:   blecommands.LogSortOrderDescending,
		TotalCount:  withCount,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	var entries []blecommands.LogEntry
	var logCount *blecommands.LogEntryCount
	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to decrypt log response: %w", err)
			}
			slog.Debug("Received log entry response", "cmd", res.GetCommandCode(), "payload", res)
			switch r := res.(type) {
			case *blecommands.LogEntry:
				entries = append(entries, *r)
			case *blecommands.LogEntryCount:
				logCount = r
			case *blecommands.Status:
				if r.Status == blecommands.StatusComplete {
					return entries, logCount, nil
				}
			}
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
}
