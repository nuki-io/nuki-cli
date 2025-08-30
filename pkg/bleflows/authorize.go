package bleflows

import (
	crypto_rand "crypto/rand"
	"fmt"
	"log/slog"
	"os"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
)

func (f *Flow) Authorize(pin string) error {
	f.authCtx = NewAuthorizeContext()
	f.authCtx.Pin = pin
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

	if _, ok := res.(*blecommands.AuthorizationInfo); ok {
		err = f.auth5G(res)
	} else {
		err = f.authPre5G(res)
	}
	if err != nil {
		return err
	}

	f.device.DiscoverKeyturnerUsdio()
	f.initializeHandlerWithCrypto()
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
	f.UpdateAuthCtxFromConfig(res.(*blecommands.Config))

	return nil
}

func (f *Flow) authPre5G(res blecommands.Command) error {
	challenge := res.(*blecommands.Challenge).Nonce
	slog.Debug("Received challenge", "challenge", fmt.Sprintf("%x", challenge))

	authData := &blecommands.AuthorizationData{
		IdType: 0x00,
		Id:     f.authCtx.AppId,
		Name:   getAuthName(),
		Nonce:  GetNonce32(),
	}
	// at this point, the payload will not contain the authenticator as it has 0 length
	authenticator := f.authCtx.GetMessageAuthenticator(authData.GetPayload(), challenge)
	// after the next line the payload contains the authenticator
	authData.Authenticator = authenticator

	msg := f.handler.ToMessage(authData)
	res, err := f.handler.FromDeviceResponse(f.device.WritePairing(msg))
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
	return nil
}

func (f *Flow) auth5G(res blecommands.Command) error {
	// TODO: if security pin is not set, should we also not send in the auth info?
	_ = res.(*blecommands.AuthorizationInfo)

	// temporarily set to fixed AuthID 0x7FFFFFFF
	// see Nuki BLE Spec, Example for 5G Authorization
	f.authCtx.AuthId = []byte{0x7F, 0xFF, 0xFF, 0xFF}
	f.initializeHandlerWithCrypto()

	authData := &blecommands.AuthorizationData5G{
		Id:          f.authCtx.AppId,
		Name:        getAuthName(),
		SecurityPin: blecommands.NewPin(f.authCtx.Pin),
	}
	msg := f.handler.ToEncryptedMessage(authData, GetNonce24())
	res, err := f.handler.FromEncryptedDeviceResponse(f.device.WritePairing(msg))

	if err != nil {
		return fmt.Errorf("failed to send authorization data: %w", err)
	}
	authId := res.(*blecommands.AuthorizationID)
	f.authCtx.AuthId = authId.AuthId
	slog.Debug("Received authorization data", "authId", f.authCtx.AuthId)

	slog.Info("Pairing complete")

	return nil
}

func getAuthName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown host"
	}
	return fmt.Sprintf("Nuki CLI (%s)", hostname)
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
