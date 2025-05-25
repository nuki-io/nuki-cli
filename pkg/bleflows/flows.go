package bleflows

import (
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
)

type Flow struct {
	ble     *nukible.NukiBle
	handler *blecommands.BleHandler
	device  *nukible.Device
	authCtx *AuthorizeContext
}

func NewFlow(ble *nukible.NukiBle) *Flow {
	return &Flow{
		ble: ble,
	}
}

func (f *Flow) Connect(id string) error {
	addr, ok := f.ble.GetDeviceAddress(id)
	if !ok {
		return fmt.Errorf("requested device with MAC %s was not discovered", id)
	}

	device, err := f.ble.Connect(*addr)
	if err != nil {
		return fmt.Errorf("cannot connect to device %s. %s", id, err.Error())
	}
	f.device = device
	return nil
}
func (f *Flow) LoadAuthContext(id string) {
	f.authCtx = &AuthorizeContext{}
	err := f.authCtx.Load(id)
	if err != nil {
		panic(fmt.Errorf("device is not paired yet. %s", err.Error()))
	}
}
func (f *Flow) InitializeHandler() {
	f.handler = blecommands.NewBleHandler(nil, nil)
}
func (f *Flow) InitializeHandlerWithCrypto() {
	crypto := blecommands.NewCrypto(f.authCtx.SharedKey)
	f.handler = blecommands.NewBleHandler(crypto, f.authCtx.AuthId)
}

func (f *Flow) getChallenge() ([]byte, error) {
	msg := f.handler.ToEncryptedMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandChallenge}, GetNonce24())
	deviceRes := f.device.WriteUsdio(msg)
	res, err := f.handler.FromEncryptedDeviceResponse(deviceRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}
	return res.(*blecommands.Challenge).Nonce, nil
}
