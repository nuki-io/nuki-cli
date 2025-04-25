package bleflows

import (
	crypto_rand "crypto/rand"
	"fmt"
	"log/slog"
	"os"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
)

func (f *Flow) Authorize(id string) error {
	addr, ok := f.ble.GetDeviceAddress(id)
	if !ok {
		return fmt.Errorf("requested device with MAC %s was not discovered", id)
	}

	device, err := f.ble.Connect(*addr)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to device %s. %s", id, err.Error()))
	}
	device.DiscoverPairing()
	h := blecommands.NewBleHandler(nil, nil)

	ctx := NewAuthorizeContext()
	slog.Info("Requesting public key from smartlock")
	msg := h.ToMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandPublicKey})
	res, err := h.FromDeviceResponse(device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to get public key from device: %w", err)
	}
	ctx.SlPublicKey = res.(*blecommands.PublicKey).PublicKey
	slog.Info("Received public key from smartlock", "pubkey", fmt.Sprintf("%x", ctx.SlPublicKey))

	// pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)
	// ctx.CliPublicKey = pubKey[:]
	// ctx.CliPrivateKey = privKey[:]
	// if err != nil {
	// 	panic(err)
	// }
	ctx.GenerateKeyPair()

	slog.Info("Sending CLI public key", "pubkey", fmt.Sprintf("%x", ctx.CliPublicKey))
	msg = h.ToMessage(&blecommands.PublicKey{PublicKey: ctx.CliPublicKey})
	res, err = h.FromDeviceResponse(device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send public key to device: %w", err)
	}
	challenge := res.(*blecommands.Challenge).Nonce
	slog.Debug("Received challenge", "challenge", fmt.Sprintf("%x", challenge))

	ctx.CalculateSharedKey()
	slog.Info("Calculated shared key", "sharedKey", fmt.Sprintf("%x", ctx.SharedKey))

	authenticator := ctx.GetMessageAuthenticator(ctx.CliPublicKey, ctx.SlPublicKey, challenge)
	slog.Info("Sending authenticator", "authenticator", fmt.Sprintf("%x", authenticator))
	msg = h.ToMessage(&blecommands.AuthorizationAuthenticator{Authenticator: authenticator})
	res, err = h.FromDeviceResponse(device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send authenticator to device: %w", err)
	}
	challenge = res.(*blecommands.Challenge).Nonce
	slog.Debug("Received challenge", "challenge", fmt.Sprintf("%x", challenge))

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown host"
	}
	authData := &blecommands.AuthorizationData{
		IdType: 0x00,
		Id:     ctx.AppId,
		Name:   fmt.Sprintf("Nuki CLI (%s)", hostname),
		Nonce:  GetNonce32(),
	}
	// at this point, the payload will not contain the authenticator as it has 0 length
	authenticator = ctx.GetMessageAuthenticator(authData.GetPayload(), challenge)
	// after the next line the payload contains the authenticator
	authData.Authenticator = authenticator

	msg = h.ToMessage(authData)
	res, err = h.FromDeviceResponse(device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send authorization data: %w", err)
	}
	authId := res.(*blecommands.AuthorizationID)
	ctx.AuthId = authId.AuthId
	slog.Debug("Received authorization data", "authId", ctx.AuthId, "nonceK", authId.Nonce)

	authenticator = ctx.GetMessageAuthenticator(ctx.AuthId, authId.Nonce)
	msg = h.ToMessage(&blecommands.AuthorizationIDConfirmation{Authenticator: authenticator, AuthId: ctx.AuthId})
	res, err = h.FromDeviceResponse(device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send authorization ID confirmation: %w", err)
	}
	status := res.(*blecommands.Status).Status
	slog.Info("Pairing complete", "status", status)
	ctx.Store(id)

	slog.Info("Disconnecting...")
	device.Disconnect()
	return nil
}

func GetNonce32() []byte {
	var buf [32]byte
	crypto_rand.Read(buf[:])
	return buf[:]
}

func GetNonce24() []byte {
	var buf [24]byte
	crypto_rand.Read(buf[:])
	return buf[:]
}
