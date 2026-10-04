package bleflows

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetAuthorizationEntries(ctx context.Context, offset, count uint16) ([]blecommands.AuthorizationEntry, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	return collectResponses[blecommands.AuthorizationEntry](ctx, f, &blecommands.RequestAuthorizationEntries{
		Offset:      offset,
		Count:       count,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
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
