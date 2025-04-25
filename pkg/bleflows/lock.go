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

	cmd := blecommands.NewEncryptedRequestData(crypto, ctx.AuthId, blecommands.CommandChallenge)
	res := blecommands.FromEncryptedDeviceResponse(crypto, device.WriteUsdio(cmd.ToMessage(GetNonce24())))
	nonce := res.GetPayload()

	cmd = blecommands.NewEncryptedCommand(
		crypto,
		ctx.AuthId,
		blecommands.CommandLockAction,
		slices.Concat(
			[]byte{byte(action)},
			ctx.AppId,
			[]byte{0x00},
			nonce,
		))
	device.WriteUsdioWithCallback(
		cmd.ToMessage(GetNonce24()),
		func(b []byte, c chan int) []byte { return onLockResponse(b, c, crypto) },
	)

	device.Disconnect()
	return nil
}

func onLockResponse(buf []byte, sem chan int, crypto blecommands.Crypto) []byte {
	slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
	res := blecommands.FromEncryptedDeviceResponse(crypto, buf)
	slog.Info("Received lock action response", "cmd", res.GetCommandCode(), "payload", res.GetPayload())
	if (res.GetCommandCode() == blecommands.CommandStatus && slices.Equal(res.GetPayload(), []byte{0x00})) || res.GetCommandCode() == blecommands.CommandErrorReport {
		<-sem
	}
	return buf
}
