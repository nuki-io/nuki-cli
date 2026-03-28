package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetConfig(ctx context.Context) (*blecommands.Config, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestConfig{Nonce: nonce}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}

	return res.(*blecommands.Config), nil
}

func (f *Flow) RequestData(ctx context.Context, cmd blecommands.CommandCode) (*blecommands.Response, error) {
	cfg := &blecommands.RequestData{CommandIdentifier: cmd}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to request data from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to request data from device: %w", err)
	}

	return &res, nil
}

func (f *Flow) GetStatus(ctx context.Context) (*blecommands.KeyturnerStates, error) {
	res, err := f.RequestData(ctx, blecommands.CommandKeyturnerStates)
	if err != nil {
		return nil, fmt.Errorf("failed to request keyturner states: %w", err)
	}
	state, ok := (*res).(*blecommands.KeyturnerStates)
	if !ok {
		return nil, fmt.Errorf("failed to cast response to KeyturnerStates: %w", err)
	}
	return state, nil
}

func (f *Flow) GetLogs(ctx context.Context, start int, count int) ([]blecommands.LogEntry, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestLogEntries{
		StartIndex:  uint32(start),
		Count:       uint16(count),
		Nonce:       nonce,
		SortOrder:   blecommands.LogSortOrderDescending,
		TotalCount:  0x00,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	var entries []blecommands.LogEntry
	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt log response: %w", err)
			}
			slog.Debug("Received log entry response", "cmd", res.GetCommandCode(), "payload", res)
			switch r := res.(type) {
			case *blecommands.LogEntry:
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
