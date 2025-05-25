package bleflows

import (
	crypto_rand "crypto/rand"
	"fmt"
	"log/slog"
	"os"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) Authorize(id string) error {
	f.Connect(id)
	f.device.DiscoverPairing()
	f.InitializeHandler()

	f.authCtx = NewAuthorizeContext()
	slog.Info("Requesting public key from smartlock")
	msg := f.handler.ToMessage(&blecommands.RequestData{CommandIdentifier: blecommands.CommandPublicKey})
	res, err := f.handler.FromDeviceResponse(f.device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to get public key from device: %w", err)
	}
	f.authCtx.SlPublicKey = res.(*blecommands.PublicKey).PublicKey
	slog.Info("Received public key from smartlock", "pubkey", fmt.Sprintf("%x", f.authCtx.SlPublicKey))

	f.authCtx.GenerateKeyPair()

	slog.Info("Sending CLI public key", "pubkey", fmt.Sprintf("%x", f.authCtx.CliPublicKey))
	msg = f.handler.ToMessage(&blecommands.PublicKey{PublicKey: f.authCtx.CliPublicKey})
	res, err = f.handler.FromDeviceResponse(f.device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send public key to device: %w", err)
	}
	challenge := res.(*blecommands.Challenge).Nonce
	slog.Debug("Received challenge", "challenge", fmt.Sprintf("%x", challenge))

	f.authCtx.CalculateSharedKey()
	slog.Info("Calculated shared key", "sharedKey", fmt.Sprintf("%x", f.authCtx.SharedKey))

	authenticator := f.authCtx.GetMessageAuthenticator(f.authCtx.CliPublicKey, f.authCtx.SlPublicKey, challenge)
	slog.Info("Sending authenticator", "authenticator", fmt.Sprintf("%x", authenticator))
	msg = f.handler.ToMessage(&blecommands.AuthorizationAuthenticator{Authenticator: authenticator})
	res, err = f.handler.FromDeviceResponse(f.device.WritePairing(msg))
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
		Id:     f.authCtx.AppId,
		Name:   fmt.Sprintf("Nuki CLI (%s)", hostname),
		Nonce:  GetNonce32(),
	}
	// at this point, the payload will not contain the authenticator as it has 0 length
	authenticator = f.authCtx.GetMessageAuthenticator(authData.GetPayload(), challenge)
	// after the next line the payload contains the authenticator
	authData.Authenticator = authenticator

	msg = f.handler.ToMessage(authData)
	res, err = f.handler.FromDeviceResponse(f.device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send authorization data: %w", err)
	}
	authId := res.(*blecommands.AuthorizationID)
	f.authCtx.AuthId = authId.AuthId
	slog.Debug("Received authorization data", "authId", f.authCtx.AuthId, "nonceK", authId.Nonce)

	authenticator = f.authCtx.GetMessageAuthenticator(f.authCtx.AuthId, authId.Nonce)
	msg = f.handler.ToMessage(&blecommands.AuthorizationIDConfirmation{Authenticator: authenticator, AuthId: f.authCtx.AuthId})
	res, err = f.handler.FromDeviceResponse(f.device.WritePairing(msg))
	if err != nil {
		return fmt.Errorf("failed to send authorization ID confirmation: %w", err)
	}
	status := res.(*blecommands.Status).Status
	slog.Info("Pairing complete", "status", status)

	f.device.DiscoverKeyturnerUsdio()
	f.InitializeHandlerWithCrypto()
	nonce, err := f.getChallenge()
	if err != nil {
		return fmt.Errorf("failed to get challenge from device: %w", err)
	}

	slog.Info("Reading config from smartlock")
	cfg := &blecommands.RequestConfig{Nonce: nonce}
	msg = f.handler.ToEncryptedMessage(cfg, GetNonce24())
	res, err = f.handler.FromEncryptedDeviceResponse(f.device.WriteUsdio(msg))
	if err != nil {
		return fmt.Errorf("failed to get config from device: %w", err)
	}
	f.authCtx.Name = res.(*blecommands.Config).Name

	f.authCtx.Store(id)

	slog.Info("Disconnecting...")
	f.device.Disconnect()
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
