package bleflows

import (
	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"slices"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
	"golang.org/x/crypto/nacl/box"
)

func (f *Flow) Authorize(mac string) error {
	addr, ok := f.ble.GetDeviceAddress(mac)
	if !ok {
		return fmt.Errorf("requested device with MAC %s was not discovered", mac)
	}

	device, err := f.ble.Connect(*addr)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to device %s. %s", mac, err.Error()))
	}
	device.DiscoverPairing()

	ctx := &AuthorizeContext{}
	fmt.Println("Requesting SL public key")
	cmd := blecommands.NewUnencryptedRequestData(blecommands.PublicKey)
	res := blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	ctx.slPublicKey = res.GetPayload()
	fmt.Printf("SL public key: %x\n", ctx.slPublicKey)

	pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)
	ctx.cliPublicKey = pubKey[:]
	ctx.cliPrivateKey = privKey[:]
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending public key:", ctx.cliPublicKey)
	cmd = blecommands.NewUnencryptedCommand(blecommands.PublicKey, ctx.cliPublicKey)
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	challenge := res.GetPayload()
	fmt.Printf("Received challenge: %x\n", challenge)

	ctx.CalculateSharedKey()
	fmt.Printf("Calculated shared key: %x\n", ctx.sharedKey)

	authenticator := ctx.GetMessageAuthenticator(ctx.cliPublicKey, ctx.slPublicKey, challenge)
	fmt.Printf("Sending authenticator: %x\n", authenticator)
	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationAuthenticator,
		authenticator)
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	challenge = res.GetPayload()
	fmt.Printf("Received challenge: %x\n", challenge)

	appName := [32]byte{}
	copy(appName[:], "Nuki CLI")
	payload := slices.Concat(
		[]byte{0x00},                   // App,
		[]byte{0x27, 0xED, 0x7E, 0x18}, // From the example. should be random for each app
		appName[:],
		GetNonce(),
	)
	authenticator = ctx.GetMessageAuthenticator(payload, challenge)

	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationData,
		slices.Concat(
			authenticator,
			payload,
		),
	)
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	authId := res.GetPayload()[32 : 32+4]
	fmt.Printf("Received AuthId: %x\n", authId)
	nonceK := res.GetPayload()[52:]
	fmt.Printf("Received nonceK: %x\n", nonceK)

	authenticator = ctx.GetMessageAuthenticator(authId, nonceK)
	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationIDConfirmation,
		slices.Concat(
			authenticator,
			authId,
		),
	)
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	complete := res.GetPayload()
	fmt.Printf("Complete: %x\n", complete)

	fmt.Println("Disconnecting...")
	device.Disconnect()
	return nil
}

func GetNonce() []byte {
	var buf [32]byte
	crypto_rand.Read(buf[:])
	return buf[:]
}

type AuthorizeContext struct {
	cliPublicKey  []byte
	cliPrivateKey []byte
	slPublicKey   []byte
	sharedKey     []byte
}

func (ac *AuthorizeContext) SetCliKeys(priv []byte, pub []byte) {
	ac.cliPrivateKey = priv
	ac.cliPublicKey = pub
}
func (ac *AuthorizeContext) SetSmartlockPublicKey(pub []byte) {
	ac.slPublicKey = pub
}
func (ac *AuthorizeContext) GetSharedKey() []byte {
	return ac.sharedKey
}
func (ac *AuthorizeContext) GetCliPublicKey() []byte {
	return ac.cliPublicKey
}
func (ac *AuthorizeContext) GetCliPrivateKey() []byte {
	return ac.cliPrivateKey
}
func (ac *AuthorizeContext) GetSmartlockPublicKey() []byte {
	return ac.slPublicKey
}
func (ac *AuthorizeContext) CalculateSharedKey() {
	sharedKey := [32]byte{}
	box.Precompute(&sharedKey, (*[32]byte)(ac.slPublicKey), (*[32]byte)(ac.cliPrivateKey))
	ac.sharedKey = sharedKey[:]
}

func (ac *AuthorizeContext) GetMessageAuthenticator(parts ...[]byte) []byte {
	h := hmac.New(sha256.New, ac.sharedKey)
	h.Write(slices.Concat(parts...))
	return h.Sum(nil)
}
