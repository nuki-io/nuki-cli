package bleflows

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetKeypadCodes(ctx context.Context, offset, count uint16) ([]blecommands.KeypadCode, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	req := &blecommands.RequestKeypadCodes{
		Offset:      offset,
		Count:       count,
		Nonce:       nonce,
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	var codes []blecommands.KeypadCode
	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return nil, fmt.Errorf("failed to decrypt keypad codes response: %w", err)
			}
			slog.Debug("Received keypad codes response", "cmd", res.GetCommandCode())
			switch r := res.(type) {
			case *blecommands.KeypadCode:
				codes = append(codes, *r)
			case *blecommands.Status:
				if r.Status == blecommands.StatusComplete {
					return codes, nil
				}
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
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
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	ch, stop := f.device.WriteUsdioStream(ctx, msg)
	defer stop()

	for {
		select {
		case buf := <-ch:
			res, err := f.handler.FromEncryptedDeviceResponse(buf)
			if err != nil {
				return 0, fmt.Errorf("failed to decrypt add keypad code response: %w", err)
			}
			slog.Debug("Received add keypad code response", "cmd", res.GetCommandCode())
			switch r := res.(type) {
			case *blecommands.KeypadCodeID:
				return r.CodeID, nil
			case *blecommands.Status:
				if r.Status == blecommands.StatusComplete {
					return 0, nil
				}
			}
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
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
