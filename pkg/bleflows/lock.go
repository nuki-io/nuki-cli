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
	_, err = f.device.WriteUsdioWithCallback(ctx, msg, f.onLockResponse)
	return err
}

func (f *Flow) onLockResponse(buf []byte, sem chan error) []byte {
	slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
	res, err := f.handler.FromEncryptedDeviceResponse(buf)
	if err != nil {
		slog.Error("Failed to decrypt response", "err", err)
		sem <- err
		return buf
	}
	slog.Info("Received lock action response", "cmd", res.GetCommandCode(), "payload", res)
	if s, ok := res.(*blecommands.Status); ok && s.Status == blecommands.StatusComplete {
		sem <- nil
	}
	return buf
}
