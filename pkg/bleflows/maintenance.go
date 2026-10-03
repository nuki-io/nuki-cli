package bleflows

import (
	"context"
	"fmt"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

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
