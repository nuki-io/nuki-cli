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

	ctx := NewAuthorizeContext()
	slog.Info("Requesting public key from smartlock")
	cmd := blecommands.NewUnencryptedRequestData(blecommands.PublicKey)
	res := blecommands.FromDeviceResponse(device.WritePairing(cmd.ToMessage()))
	ctx.SlPublicKey = res.GetPayload()
	slog.Info("Received public key from smartlock", "pubkey", fmt.Sprintf("%x", ctx.SlPublicKey))

	pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)
	ctx.CliPublicKey = pubKey[:]
	ctx.CliPrivateKey = privKey[:]
	if err != nil {
		panic(err)
	}

	slog.Info("Sending CLI public key", "pubkey", fmt.Sprintf("%x", ctx.CliPublicKey))
	cmd = blecommands.NewUnencryptedCommand(blecommands.PublicKey, ctx.CliPublicKey)
	res = blecommands.FromDeviceResponse(device.WritePairing(cmd.ToMessage()))
	challenge := res.GetPayload()
	slog.Debug("Received challenge", "challenge", fmt.Sprintf("%x", challenge))

	ctx.CalculateSharedKey()
	slog.Info("Calculated shared key", "sharedKey", fmt.Sprintf("%x", ctx.SharedKey))

	authenticator := ctx.GetMessageAuthenticator(ctx.CliPublicKey, ctx.SlPublicKey, challenge)
	slog.Info("Sending authenticator", "authenticator", fmt.Sprintf("%x", authenticator))
	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationAuthenticator,
		authenticator)
	res = blecommands.FromDeviceResponse(device.WritePairing(cmd.ToMessage()))
	challenge = res.GetPayload()
	slog.Debug("Received challenge", "challenge", fmt.Sprintf("%x", challenge))

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown host"
	}
	appName := [32]byte{}
	copy(appName[:], fmt.Appendf(nil, "Nuki CLI (%s)", hostname))
	payload := slices.Concat(
		[]byte{0x00}, // App,
		ctx.AppId,
		appName[:],
		GetNonce32(),
	)
	authenticator = ctx.GetMessageAuthenticator(payload, challenge)

	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationData,
		slices.Concat(
			authenticator,
			payload,
		),
	)
	res = blecommands.FromDeviceResponse(device.WritePairing(cmd.ToMessage()))
	ctx.AuthId = res.GetPayload()[32 : 32+4]
	nonceK := res.GetPayload()[52:]
	slog.Debug("Received authorization data", "authId", ctx.AuthId, "nonceK", nonceK)

	authenticator = ctx.GetMessageAuthenticator(ctx.AuthId, nonceK)
	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationIDConfirmation,
		slices.Concat(
			authenticator,
			ctx.AuthId,
		),
	)
	res = blecommands.FromDeviceResponse(device.WritePairing(cmd.ToMessage()))
	complete := res.GetPayload()
	slog.Debug("Pairing complete", "complete", fmt.Sprintf("%x", complete))
	ctx.Store(id)

	slog.Debug("Disconnecting...")
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
