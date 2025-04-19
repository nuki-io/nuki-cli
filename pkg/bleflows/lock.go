package bleflows

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
)

func (f *Flow) PerformLockOperation(mac string, action blecommands.Action) error {
	addr, ok := f.ble.GetDeviceAddress(mac)
	if !ok {
		return fmt.Errorf("requested device with MAC %s was not discovered", mac)
	}

	device, err := f.ble.Connect(*addr)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to device %s. %s", mac, err.Error()))
	}
	device.DiscoverKeyturnerUsdio()

	crypto := blecommands.NewCrypto(ctx.SharedKey)

	cmd := blecommands.NewEncryptedRequestData(crypto, ac.AuthId, blecommands.Challenge)
	res := blecommands.FromEncryptedDeviceResponse(crypto, device.WriteUsdio(cmd.ToMessage(GetNonce24())))
	nonce := res.GetPayload()

	cmd = blecommands.NewEncryptedCommand(
		crypto,
		ac.AuthId,
		blecommands.LockAction,
		slices.Concat(
			[]byte{byte(action)},
			[]byte{0x27, 0xED, 0x7E, 0x18},
			[]byte{0x00},
			nonce,
		))
	res = blecommands.FromEncryptedDeviceResponse(crypto, device.WriteUsdio(cmd.ToMessage(GetNonce24())))
	fmt.Printf("%x", res.GetPayload())

	// TODO: should read intermediate states as well

	device.Disconnect()
	return nil
}

func loadAuthContext() *AuthorizeContext {
	j, err := os.ReadFile("./ac.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var ac AuthorizeContext
	err = json.Unmarshal(j, &ac)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return &ac
}
