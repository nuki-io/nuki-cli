package nukible

import (
	"fmt"
	"log/slog"

	"tinygo.org/x/bluetooth"
)

type Device struct {
	btDev           bluetooth.Device
	services        []bluetooth.DeviceService
	characteristics []bluetooth.DeviceCharacteristic

	pairingGdioChar    bluetooth.DeviceCharacteristic
	keyturnerUsdioChar bluetooth.DeviceCharacteristic
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
		return fmt.Errorf("could not discover any pairing services or characteristics. %s", err.Error())
	}
	if len(n.services) != 1 {
		return fmt.Errorf("expected exactly one pairing service, got %d", len(n.services))
	}
	if len(n.characteristics) != 1 {
		return fmt.Errorf("expected exactly one GDIO characteristic, got %d", len(n.characteristics))
	}
	n.pairingGdioChar = n.characteristics[0]
	slog.Debug("Discovered pairing characteristic", "uuid", n.pairingGdioChar.String())
	return nil
}

func (n *Device) DiscoverKeyturnerUsdio() error {
	err := n.DiscoverServicesAndCharacteristics(
		[]bluetooth.UUID{KeyturnerService},
		[]bluetooth.UUID{KeyturnerUsdioCharacteristic},
	)
	if err != nil {
		return fmt.Errorf("could not discover any Keyturner services or characteristics. %s", err.Error())
	}
	if len(n.services) != 1 {
		return fmt.Errorf("expected exactly one Keyturner service, got %d", len(n.services))
	}
	if len(n.characteristics) != 1 {
		return fmt.Errorf("expected exactly one USDIO characteristic, got %d", len(n.characteristics))
	}
	n.keyturnerUsdioChar = n.characteristics[0]
	slog.Debug("Discovered Keyturner USDIO characteristic", "uuid", n.keyturnerUsdioChar.String())
	return nil
}

func (n *Device) Disconnect() {
	err := n.btDev.Disconnect()
	if err != nil {
		slog.Error("Error disconnecting from device", "error", err)
	}
	n.services = make([]bluetooth.DeviceService, 0)
	n.characteristics = make([]bluetooth.DeviceCharacteristic, 0)
}

func (n *Device) WritePairing(data []byte) []byte {
	return n.write(n.pairingGdioChar, data, onGdioNotify)
}

func (n *Device) WriteUsdio(data []byte) []byte {
	return n.write(n.keyturnerUsdioChar, data, onGdioNotify)
}

func (n *Device) WriteUsdioWithCallback(data []byte, cb func([]byte, chan int) []byte) []byte {
	return n.write(n.keyturnerUsdioChar, data, cb)
}

func (n *Device) write(char bluetooth.DeviceCharacteristic, data []byte, cb func([]byte, chan int) []byte) []byte {
	sem := make(chan int, 1)
	sem <- 1

	slog.Debug("Writing bytes to characteristic", "data", fmt.Sprintf("%x", data))
	rxData := make([]byte, 0)
	char.EnableNotifications(func(buf []byte) { rxData = cb(buf, sem) })
	n.osWrite(char, data)

	slog.Debug("Waiting for response...")
	sem <- 1
	// disable notifications again - TODO: sensible, or should we just enable it once?
	char.EnableNotifications(nil)

	return rxData
}

func onGdioNotify(buf []byte, sem chan int) []byte {
	slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
	<-sem
	return buf
}
