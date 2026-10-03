package bleflows

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetConfig(ctx context.Context) (*blecommands.Config, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestConfig{Nonce: nonce}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}

	return res.(*blecommands.Config), nil
}

func (f *Flow) RequestData(ctx context.Context, cmd blecommands.CommandCode) (*blecommands.Response, error) {
	cfg := &blecommands.RequestData{CommandIdentifier: cmd}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to request data from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to request data from device: %w", err)
	}

	return &res, nil
}

func (f *Flow) GetStatus(ctx context.Context) (*blecommands.KeyturnerStates, error) {
	res, err := f.RequestData(ctx, blecommands.CommandKeyturnerStates)
	if err != nil {
		return nil, fmt.Errorf("failed to request keyturner states: %w", err)
	}
	state, ok := (*res).(*blecommands.KeyturnerStates)
	if !ok {
		return nil, fmt.Errorf("failed to cast response to KeyturnerStates: %w", err)
	}
	return state, nil
}

func (f *Flow) GetBatteryReport(ctx context.Context) (*blecommands.BatteryReport, error) {
	res, err := f.RequestData(ctx, blecommands.CommandBatteryReport)
	if err != nil {
		return nil, fmt.Errorf("failed to request battery report: %w", err)
	}
	report, ok := (*res).(*blecommands.BatteryReport)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for battery report")
	}
	return report, nil
}

func (f *Flow) GetAdvancedConfig(ctx context.Context) (*blecommands.AdvancedConfig, error) {
	nonce, err := f.getChallenge(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge: %w", err)
	}
	req := &blecommands.RequestAdvancedConfig{Nonce: nonce}
	msg := f.handler.ToEncryptedMessage(req, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get advanced config: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse advanced config: %w", err)
	}
	cfg, ok := res.(*blecommands.AdvancedConfig)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for advanced config")
	}
	return cfg, nil
}
