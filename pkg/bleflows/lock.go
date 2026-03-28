package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) PerformLockOperation(ctx context.Context, action blecommands.Action) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge from device: %w", err)
	}

	lock := &blecommands.LockAction{
		Action: action,
		AppId:  f.authCtx.AppId,
		Nonce:  nonce,
	}
	msg := f.handler.ToEncryptedMessage(lock, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return fmt.Errorf("failed to decrypt lock response: %w", err)
			}
			slog.Info("Received lock action response", "cmd", res.GetCommandCode(), "payload", res)
			if s, ok := res.(*blecommands.Status); ok && s.Status == blecommands.StatusComplete {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
