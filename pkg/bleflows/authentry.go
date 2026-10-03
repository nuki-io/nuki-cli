package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetAuthorizationEntries(ctx context.Context, offset, count uint16) ([]blecommands.AuthorizationEntry, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	req := &blecommands.RequestAuthorizationEntries{
		Offset:      offset,
		Count:       count,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	var entries []blecommands.AuthorizationEntry
	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt authorization entries response: %w", err)
			}
			slog.Debug("Received authorization entries response", "cmd", res.GetCommandCode())
			switch r := res.(type) {
			case *blecommands.AuthorizationEntry:
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

func (f *Flow) RemoveAuthorizationEntry(ctx context.Context, authID uint32) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.RemoveAuthorizationEntry{
		AuthorizationID: authID,
		Nonce:           nonce,
		SecurityPin:     blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) UpdateAuthorizationEntry(ctx context.Context, req *blecommands.UpdateAuthorizationEntry) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	req.Nonce = nonce
	req.SecurityPin = blecommands.NewPin(f.authCtx.Pin)
	return f.performSimpleOp(ctx, req)
}
