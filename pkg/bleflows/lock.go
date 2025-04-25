package bleflows

import (
	"fmt"
	"log/slog"
	"slices"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
)

func (f *Flow) PerformLockOperation(id string, action blecommands.Action) error {
	addr, ok := f.ble.GetDeviceAddress(id)
	if !ok {
		return fmt.Errorf("requested device with MAC %s was not discovered", id)
	}

	ctx := &AuthorizeContext{}
	err := ctx.Load(id)
	if err != nil {
		return fmt.Errorf("device is not paired yet. %s", err.Error())
	}

	device, err := f.ble.Connect(*addr)
	if err != nil {
		return fmt.Errorf("cannot connect to device %s. %s", id, err.Error())
	}
	device.DiscoverKeyturnerUsdio()

	crypto := blecommands.NewCrypto(ctx.SharedKey)
	h := blecommands.NewBleHandler(crypto, ctx.AuthId)

	msg := h.ToEncryptedMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandChallenge}, GetNonce24())
	deviceRes := device.WriteUsdio(msg)
	res, err := h.FromEncryptedDeviceResponse(deviceRes)
	if err != nil {
		return fmt.Errorf("failed to get challenge from device: %w", err)
	}
	nonce := res.(*blecommands.Challenge).Nonce

	lock := &blecommands.LockAction{
		Action: action,
		AppId:  ctx.AppId,
		Nonce:  nonce,
	}
	msg = h.ToEncryptedMessage(lock, GetNonce24())
	cb := func(b []byte, c chan int) []byte { return onLockResponse(b, c, h) }
	device.WriteUsdioWithCallback(msg, cb)

	device.Disconnect()
	return nil
}

func onLockResponse(buf []byte, sem chan int, h *blecommands.BleHandler) []byte {
	slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
	res, err := h.FromEncryptedDeviceResponse(buf)
	if err != nil {
		slog.Error("Failed to decrypt response", "err", err)
		<-sem
		return buf
	}
	slog.Info("Received lock action response", "cmd", res.GetCommandCode(), "payload", res)
	if res.GetCommandCode() == blecommands.CommandStatus && slices.Equal(res.GetPayload(), []byte{0x00}) {
		<-sem
	}
	return buf
}
