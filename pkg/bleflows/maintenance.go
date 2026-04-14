package bleflows

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

// performSimpleOp sends a command that requires a challenge+PIN and waits for StatusComplete.
// The caller provides an already-built request (with nonce and pin already set).
func (f *Flow) performSimpleOp(ctx context.Context, req blecommands.Request) error {
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return fmt.Errorf("failed to decrypt response: %w", err)
			}
			slog.Debug("Received response", "cmd", res.GetCommandCode())
			if s, ok := res.(*blecommands.Status); ok && s.Status == blecommands.StatusComplete {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (f *Flow) Calibrate(ctx context.Context) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.RequestCalibration{
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) Reboot(ctx context.Context) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.RequestReboot{
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) SetSecurityPIN(ctx context.Context, newPin string) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	err = f.performSimpleOp(ctx, &blecommands.SetSecurityPIN{
		NewPin:      blecommands.NewPin(newPin),
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
	if err != nil {
		return err
	}
	f.authCtx.Pin = newPin
	f.store.Store(f.id, f.authCtx)
	return nil
}

func (f *Flow) UpdateTime(ctx context.Context, t time.Time) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.UpdateTime{
		Time:        t,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) VerifyPIN(ctx context.Context) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.VerifySecurityPIN{
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

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
