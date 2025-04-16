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

	fmt.Println("Requesting SL public key")
	cmd := blecommands.NewUnencryptedRequestData(blecommands.PublicKey)
	res := blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	slPubKey := res.GetPayload()
	fmt.Printf("SL public key: %x\n", slPubKey)

	pubKey, privKey, err := box.GenerateKey(crypto_rand.Reader)

	if err != nil {
		panic(err)
	}

	fmt.Println("Sending public key:", pubKey)
	cmd = blecommands.NewUnencryptedCommand(blecommands.PublicKey, pubKey[:])
	res = blecommands.FromDeviceResponse(device.Write(cmd.ToMessage()))
	challenge := res.GetPayload()
	fmt.Printf("Received challenge: %x\n", challenge)

	sharedKey := [32]byte{}
	box.Precompute(&sharedKey, (*[32]byte)(slPubKey), privKey)
	fmt.Printf("Calculated shared key: %x\n", sharedKey)

	h := hmac.New(sha256.New, sharedKey[:])
	h.Write(slices.Concat(pubKey[:], slPubKey, challenge))
	authenticator := h.Sum(nil)

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
	h = hmac.New(sha256.New, sharedKey[:])
	h.Write(slices.Concat(payload, challenge))
	authenticator = h.Sum(nil)

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

	h = hmac.New(sha256.New, sharedKey[:])
	h.Write(slices.Concat(authId, nonceK))
	authenticator = h.Sum(nil)
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
