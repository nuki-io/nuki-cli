package nukible

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

type Device struct {
	btDev           bluetooth.Device
	services        []bluetooth.DeviceService
	characteristics []bluetooth.DeviceCharacteristic

	pairingGdioChar bluetooth.DeviceCharacteristic
}

func (n *Device) DiscoverServicesAndCharacteristics(services []bluetooth.UUID, chars []bluetooth.UUID) error {
	s, err := n.btDev.DiscoverServices(services)
	if err != nil {
		return err
	}
	n.services = s

	c, err := s[0].DiscoverCharacteristics(chars)
	if err != nil {
		return err
	}
	n.characteristics = c
	return nil
}

func (n *Device) DiscoverPairing() error {
	err := n.DiscoverServicesAndCharacteristics(
		[]bluetooth.UUID{KeyturnerPairingService},
		[]bluetooth.UUID{KeyturnerPairingGdioCharacteristic},
	)
	if len(n.services) == 0 && err != nil {
		// expected, maybe it's an Ultra
		err = n.DiscoverServicesAndCharacteristics(
			[]bluetooth.UUID{KeyturnerPairingServiceUltra},
			[]bluetooth.UUID{KeyturnerPairingGdioCharacteristic},
		)
	}
	if err != nil {
		return fmt.Errorf("Could not discover any pairing services or characteristics. %s", err.Error())
	}
	if len(n.services) != 1 {
		return fmt.Errorf("Expected exactly one pairing service, got %d", len(n.services))
	}
	if len(n.characteristics) != 1 {
		return fmt.Errorf("Expected exactly one GDIO characteristic, got %d", len(n.characteristics))
	}
	n.pairingGdioChar = n.characteristics[0]
	fmt.Println("Characteristic", n.pairingGdioChar.String())
	return nil
}

func (n *Device) Disconnect() {
	n.btDev.Disconnect()
	n.services = make([]bluetooth.DeviceService, 0)
	n.characteristics = make([]bluetooth.DeviceCharacteristic, 0)
}

func onGdioNotify(buf []byte, sem chan int) []byte {
	fmt.Printf("Received response: %x\n", buf)
	<-sem
	return buf
}
