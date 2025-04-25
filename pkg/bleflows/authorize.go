package bleflows

import (
	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"slices"

	"github.com/spf13/viper"
	"go.nuki.io/nuki/nukictl/pkg/blecommands"
	"golang.org/x/crypto/nacl/box"
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

	pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)
	ctx.CliPublicKey = pubKey[:]
	ctx.CliPrivateKey = privKey[:]
	if err != nil {
		panic(err)
	}

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

type AuthorizeContext struct {
	CliPublicKey  []byte
	CliPrivateKey []byte
	SlPublicKey   []byte
	SharedKey     []byte
	AuthId        []byte
	AppId         []byte
}

func (ac *AuthorizeContext) toStorage() *authorizeContextStorage {
	return &authorizeContextStorage{
		CliPublicKey:  fmt.Sprintf("%x", ac.CliPublicKey),
		CliPrivateKey: fmt.Sprintf("%x", ac.CliPrivateKey),
		SlPublicKey:   fmt.Sprintf("%x", ac.SlPublicKey),
		SharedKey:     fmt.Sprintf("%x", ac.SharedKey),
		AuthId:        fmt.Sprintf("%x", ac.AuthId),
		AppId:         fmt.Sprintf("%x", ac.AppId),
	}
}
func (ac *AuthorizeContext) fromStorage(s *authorizeContextStorage) {
	if v, err := hex.DecodeString(s.CliPublicKey); err == nil {
		ac.CliPublicKey = v
	}
	if v, err := hex.DecodeString(s.CliPrivateKey); err == nil {
		ac.CliPrivateKey = v
	}
	if v, err := hex.DecodeString(s.SlPublicKey); err == nil {
		ac.SlPublicKey = v
	}
	if v, err := hex.DecodeString(s.SharedKey); err == nil {
		ac.SharedKey = v
	}
	if v, err := hex.DecodeString(s.AuthId); err == nil {
		ac.AuthId = v
	}
	if v, err := hex.DecodeString(s.AppId); err == nil {
		ac.AppId = v
	}
}

type authorizeContextStorage struct {
	CliPublicKey  string
	CliPrivateKey string
	SlPublicKey   string
	SharedKey     string
	AuthId        string
	AppId         string
}

func NewAuthorizeContext() *AuthorizeContext {
	ctx := &AuthorizeContext{}
	ctx.AppId = binary.LittleEndian.AppendUint32(ctx.AppId, rand.Uint32())
	return ctx
}

func (ac *AuthorizeContext) Load(id string) error {
	cfgKey := fmt.Sprintf("authorizations.%s", id)
	if !viper.IsSet(cfgKey) {
		return fmt.Errorf("no authorization for device with id %s found", id)
	}
	s := &authorizeContextStorage{}
	viper.UnmarshalKey(cfgKey, s)
	ac.fromStorage(s)
	return nil
}

func (ac *AuthorizeContext) Store(id string) {
	cfgKey := fmt.Sprintf("authorizations.%s", id)
	viper.Set(cfgKey, ac.toStorage())
}

func (ac *AuthorizeContext) CalculateSharedKey() {
	sharedKey := [32]byte{}
	box.Precompute(&sharedKey, (*[32]byte)(ac.SlPublicKey), (*[32]byte)(ac.CliPrivateKey))
	ac.SharedKey = sharedKey[:]
}

func (ac *AuthorizeContext) GetMessageAuthenticator(parts ...[]byte) []byte {
	h := hmac.New(sha256.New, ac.SharedKey)
	h.Write(slices.Concat(parts...))
	return h.Sum(nil)
}
