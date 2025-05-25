package bleflows

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) PerformLockOperation(id string, action blecommands.Action) error {
	f.LoadAuthContext(id)
	f.Connect(id)
	f.device.DiscoverKeyturnerUsdio()
	f.InitializeHandlerWithCrypto()

	nonce, err := f.getChallenge()
	if err != nil {
		return fmt.Errorf("failed to get challenge from device: %w", err)
	}

	lock := &blecommands.LockAction{
		Action: action,
		AppId:  f.authCtx.AppId,
		Nonce:  nonce,
	}
	msg := f.handler.ToEncryptedMessage(lock, GetNonce24())
	f.device.WriteUsdioWithCallback(msg, f.onLockResponse)

	f.device.Disconnect()
	return nil
}

func (f *Flow) onLockResponse(buf []byte, sem chan int) []byte {
	slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
	res, err := f.handler.FromEncryptedDeviceResponse(buf)
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
