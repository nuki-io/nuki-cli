package bleflows

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
)

type Flow struct {
	ble     *nukible.NukiBle
	handler *blecommands.BleHandler
	device  *nukible.Device
	authCtx *AuthorizeContext
	store   AuthStore
	id      string
}

// NewAuthenticatedFlow creates a new Flow instance for a Nuki device that was already paired.
func NewAuthenticatedFlow(ble *nukible.NukiBle, id string, store AuthStore) (*Flow, error) {
	f := &Flow{
		ble:   ble,
		id:    id,
		store: store,
	}
	err := f.loadAuthContext(id)
	if err != nil {
		return nil, err
	}
	err = f.connect(id)
	if err != nil {
		return nil, err
	}
	err = f.device.DiscoverKeyturnerUsdio()
	if err != nil {
		return nil, err
	}
	f.initializeHandlerWithCrypto()

	return f, nil
}

// NewUnauthenticatedFlow creates a new Flow instance for a Nuki device that has not been paired yet.
func NewUnauthenticatedFlow(ble *nukible.NukiBle, id string, store AuthStore) (*Flow, error) {
	f := &Flow{
		ble:   ble,
		id:    id,
		store: store,
	}
	err := f.connect(id)
	if err != nil {
		return nil, err
	}
	err = f.device.DiscoverPairing()
	if err != nil {
		return nil, err
	}
	f.initializeHandler()

	return f, nil
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
	ctx, err := f.store.Load(id)
	if err != nil {
		return fmt.Errorf("device is not paired yet. %s", err.Error())
	}
	f.authCtx = ctx
	return nil
}

func (f *Flow) initializeHandler() {
	f.handler = blecommands.NewBleHandler(nil, nil)
}

func (f *Flow) initializeHandlerWithCrypto() {
	crypto := blecommands.NewCrypto(f.authCtx.SharedKey)
	f.handler = blecommands.NewBleHandler(crypto, f.authCtx.AuthId)
}

func (f *Flow) getChallenge(ctx context.Context) ([]byte, error) {
	msg := f.handler.ToEncryptedMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandChallenge}, GetNonce24())
	raw, err := f.device.WriteUsdio(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}
	res, err := f.handler.FromEncryptedDeviceResponse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}
	return res.(*blecommands.Challenge).Nonce, nil
}

func (f *Flow) UpdateAuthCtxFromConfig(cfg *blecommands.Config) {
	f.authCtx.Name = cfg.Name
	f.authCtx.NukiId = cfg.NukiID
	f.store.Store(f.id, f.authCtx)
}

func (f *Flow) DisconnectDevice() error {
	if f.device == nil {
		return fmt.Errorf("no device connected")
	}
	f.device.Disconnect()
	f.device = nil
	return nil
}
