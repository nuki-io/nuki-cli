package bleflows

import (
	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
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
	ctx.SlPublicKey = res.GetPayload()
	fmt.Printf("SL public key: %x\n", ctx.SlPublicKey)

	pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)
	ctx.CliPublicKey = pubKey[:]
	ctx.CliPrivateKey = privKey[:]
	if err != nil {
		panic(err)
	}

	fmt.Println("Sending public key:", ctx.CliPublicKey)
	cmd = blecommands.NewUnencryptedCommand(blecommands.PublicKey, ctx.CliPublicKey)
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	challenge := res.GetPayload()
	fmt.Printf("Received challenge: %x\n", challenge)

	ctx.CalculateSharedKey()
	fmt.Printf("Calculated shared key: %x\n", ctx.SharedKey)

	authenticator := ctx.GetMessageAuthenticator(ctx.CliPublicKey, ctx.SlPublicKey, challenge)
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
	ctx.AuthId = res.GetPayload()[32 : 32+4]
	fmt.Printf("Received AuthId: %x\n", ctx.AuthId)
	nonceK := res.GetPayload()[52:]
	fmt.Printf("Received nonceK: %x\n", nonceK)

	authenticator = ctx.GetMessageAuthenticator(ctx.AuthId, nonceK)
	cmd = blecommands.NewUnencryptedCommand(
		blecommands.AuthorizationIDConfirmation,
		slices.Concat(
			authenticator,
			ctx.AuthId,
		),
	)
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	complete := res.GetPayload()
	fmt.Printf("Complete: %x\n", complete)
	ctx.DumpJson()

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
	CliPublicKey  []byte
	CliPrivateKey []byte
	SlPublicKey   []byte
	SharedKey     []byte
	AuthId        []byte
}

func (ac *AuthorizeContext) DumpJson() {
	j, err := json.Marshal(ac)
	if err != nil {
		panic(err)
	}
	os.WriteFile("./ac.json", j, 0644)
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
