package nukible

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"

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

	if len(s) != 1 {
		return fmt.Errorf("expected exactly one service, got %d", len(s))
	}
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
			[]bluetooth.UUID{KeyturnerPairingGdioCharacteristicUltra},
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

// stream enables notifications on char, writes data, and returns a channel that
// receives each incoming BLE notification packet and a stop function. The caller
// must call stop() when done to disable notifications. If the receive buffer
// (capacity 32) fills up, excess packets are dropped with a warning.
func (n *Device) stream(char bluetooth.DeviceCharacteristic, data []byte) (<-chan []byte, func()) {
	ch := make(chan []byte, 32)
	var stopped atomic.Bool

	char.EnableNotifications(func(buf []byte) {
		if stopped.Load() {
			return
		}
		copied := make([]byte, len(buf))
		copy(copied, buf)
		select {
		case ch <- copied:
		default:
			slog.Warn("BLE notification dropped: receive buffer full")
		}
	})

	slog.Debug("Writing bytes to characteristic", "data", fmt.Sprintf("%x", data))
	n.osWrite(char, data)

	stop := func() {
		stopped.Store(true)
		char.EnableNotifications(nil)
	}
	return ch, stop
}

// readOne reads the first packet from ch, calls stop, and returns the packet.
// Used for single-response BLE commands.
func readOne(ctx context.Context, ch <-chan []byte, stop func()) ([]byte, error) {
	defer stop()
	slog.Debug("Waiting for response...")
	select {
	case buf := <-ch:
		slog.Debug("Received response", "buf", fmt.Sprintf("%x", buf))
		return buf, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (n *Device) WritePairing(ctx context.Context, data []byte) ([]byte, error) {
	ch, stop := n.stream(n.pairingGdioChar, data)
	return readOne(ctx, ch, stop)
}

func (n *Device) WriteUsdio(ctx context.Context, data []byte) ([]byte, error) {
	ch, stop := n.stream(n.keyturnerUsdioChar, data)
	return readOne(ctx, ch, stop)
}

// WriteUsdioStream sends data and returns a channel of raw BLE notification packets
// and a stop function. The caller consumes packets until the protocol signals
// completion (StatusComplete), then calls stop(). The context is not monitored
// here — callers should select on ctx.Done() alongside the channel.
func (n *Device) WriteUsdioStream(ctx context.Context, data []byte) (<-chan []byte, func()) {
	return n.stream(n.keyturnerUsdioChar, data)
}
