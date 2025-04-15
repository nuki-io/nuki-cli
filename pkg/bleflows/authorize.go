package bleflows

import (
	"fmt"

	"go.nuki.io/nuki/nukictl/pkg/blecommands"
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

	req := blecommands.NewUnencryptedRequestData(blecommands.PublicKey)
	res := blecommands.FromDeviceResponse(device.Write(req.ToMessage()))
	fmt.Println(res)
	fmt.Println("Disconnecting...")
	device.Disconnect()
	return nil
}
