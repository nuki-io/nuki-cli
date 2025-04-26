package bleflows

import (
	"fmt"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
)

func (f *Flow) GetConfig(id string) (*blecommands.Config, error) {
	f.LoadAuthContext(id)
	f.Connect(id)
	f.device.DiscoverKeyturnerUsdio()
	f.InitializeHandlerWithCrypto()

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

	f.device.Disconnect()
	return res.(*blecommands.Config), nil
}
