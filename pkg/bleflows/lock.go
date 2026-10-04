package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

// PerformSimpleLockAction sends a SimpleLockAction (0x0100), which needs no security PIN or app ID.
func (f *Flow) PerformSimpleLockAction(ctx context.Context, action blecommands.Action) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge from device: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.SimpleLockAction{
		Action:  action,
		NonceNK: nonce,
	})
}

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
	return f.exchange(ctx, lock, func(res blecommands.Response) bool {
		slog.Info("Received lock action response", "cmd", res.GetCommandCode(), "payload", res)
		return false
	})
}
