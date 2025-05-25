package bleflows

import (
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) GetConfig() (*blecommands.Config, error) {
	nonce, err := f.getChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	cfg := &blecommands.RequestConfig{Nonce: nonce}
	msg := f.handler.ToEncryptedMessage(cfg, GetNonce24())
	res, err := f.handler.FromEncryptedDeviceResponse(f.device.WriteUsdio(msg))
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}

	return res.(*blecommands.Config), nil
}
