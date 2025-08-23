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
	id      string
}

// NewAuthenticatedFlow creates a new Flow instance for a Nuki device that was already paired.
func NewAuthenticatedFlow(ble *nukible.NukiBle, id string) *Flow {
	f := &Flow{
		ble: ble,
		id:  id,
	}
	f.loadAuthContext(id)
	f.connect(id)
	f.device.DiscoverKeyturnerUsdio()
	f.initializeHandlerWithCrypto()

	return f
}

// NewAuthenticatedFlow creates a new Flow instance for a Nuki device that was already paired.
func NewUnauthenticatedFlow(ble *nukible.NukiBle, id string) *Flow {
	f := &Flow{
		ble: ble,
		id:  id,
	}
	f.connect(id)
	f.device.DiscoverPairing()
	f.initializeHandler()

	return f
}

func (f *Flow) connect(id string) error {
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
func (f *Flow) loadAuthContext(id string) error {
	f.authCtx = &AuthorizeContext{}
	err := f.authCtx.Load(id)
	if err != nil {
		return fmt.Errorf("device is not paired yet. %s", err.Error())
	}
	return nil
}
func (f *Flow) initializeHandler() {
	f.handler = blecommands.NewBleHandler(nil, nil)
}
func (f *Flow) initializeHandlerWithCrypto() {
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

func (f *Flow) UpdateAuthCtxFromConfig(cfg *blecommands.Config) {
	f.authCtx.Name = cfg.Name
	f.authCtx.NukiId = cfg.NukiID
	f.authCtx.Store(f.id)
}

func (f *Flow) DisconnectDevice() error {
	if f.device == nil {
		return fmt.Errorf("no device connected")
	}
	f.device.Disconnect()
	f.device = nil
	return nil
}
