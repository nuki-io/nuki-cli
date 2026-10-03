package bleflows

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) SetConfig(ctx context.Context, cfg *blecommands.SetConfig) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	cfg.Nonce = nonce
	cfg.SecurityPin = blecommands.NewPin(f.authCtx.Pin)
	return f.performSimpleOp(ctx, cfg)
}

func (f *Flow) SetAdvancedConfig(ctx context.Context, cfg *blecommands.SetAdvancedConfig) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	cfg.Nonce = nonce
	cfg.SecurityPin = blecommands.NewPin(f.authCtx.Pin)
	return f.performSimpleOp(ctx, cfg)
}
