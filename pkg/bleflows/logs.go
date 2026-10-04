package bleflows

import (
	"context"
	"fmt"

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
	var entries []blecommands.LogEntry
	var logCount *blecommands.LogEntryCount
	collectEntries := collectInto(&entries)
	err = f.exchange(ctx, cfg, func(res blecommands.Response) bool {
		if r, ok := res.(*blecommands.LogEntryCount); ok {
			logCount = r
			return false
		}
		return collectEntries(res)
	})
	if err != nil {
		return nil, nil, err
	}
	return entries, logCount, nil
}
