package bleflows

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetKeypadCodes(ctx context.Context, offset, count uint16) ([]blecommands.KeypadCode, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	return collectResponses[blecommands.KeypadCode](ctx, f, &blecommands.RequestKeypadCodes{
		Offset:      offset,
		Count:       count,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) AddKeypadCode(ctx context.Context, code uint32, name string) (uint16, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get challenge: %w", err)
	}
	req := &blecommands.AddKeypadCode{
		Code:        code,
		Name:        name,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	var codeID uint16
	// The lock sends no StatusComplete after the ID, so stop there.
	err = f.exchange(ctx, req, func(res blecommands.Response) bool {
		r, ok := res.(*blecommands.KeypadCodeID)
		if ok {
			codeID = r.CodeID
		}
		return ok
	})
	return codeID, err
}

func (f *Flow) RemoveKeypadCode(ctx context.Context, codeID uint16) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.RemoveKeypadCode{
		CodeID:      codeID,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}

func (f *Flow) UpdateKeypadCode(ctx context.Context, codeID uint16, code uint32, name string, enabled bool) error {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return fmt.Errorf("failed to get challenge: %w", err)
	}
	return f.performSimpleOp(ctx, &blecommands.UpdateKeypadCode{
		CodeID:      codeID,
		Code:        code,
		Name:        name,
		Enabled:     enabled,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	})
}
