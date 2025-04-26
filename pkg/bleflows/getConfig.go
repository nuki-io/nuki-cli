package bleflows

import (
	"fmt"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
)

func (f *Flow) GetConfig(id string) (*blecommands.Config, error) {
	ctx := &AuthorizeContext{}
	err := ctx.Load(id)
	if err != nil {
		return nil, fmt.Errorf("device is not paired yet. %s", err.Error())
	}

	f.Connect(id)
	f.device.DiscoverKeyturnerUsdio()

	crypto := blecommands.NewCrypto(ctx.SharedKey)
	f.handler = blecommands.NewBleHandler(crypto, ctx.AuthId)

	nonce, err := f.getChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}

	lock := &blecommands.RequestConfig{Nonce: nonce}
	msg := f.handler.ToEncryptedMessage(lock, GetNonce24())
	res, err := f.handler.FromEncryptedDeviceResponse(f.device.WriteUsdio(msg))
	if err != nil {
		return nil, fmt.Errorf("failed to get config from device: %w", err)
	}

	f.device.Disconnect()
	return res.(*blecommands.Config), nil
}
