package bleflows

import (
	"fmt"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
	"go.nuki.io/nuki/nukictl/pkg/nukible"
)

type Flow struct {
	ble     *nukible.NukiBle
	handler *blecommands.BleHandler
	device  *nukible.Device
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

func (f *Flow) getChallenge() ([]byte, error) {
	msg := f.handler.ToEncryptedMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandChallenge}, GetNonce24())
	deviceRes := f.device.WriteUsdio(msg)
	res, err := f.handler.FromEncryptedDeviceResponse(deviceRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get challenge from device: %w", err)
	}
	return res.(*blecommands.Challenge).Nonce, nil
}
