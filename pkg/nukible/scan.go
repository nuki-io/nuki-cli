package nukible

import (
	"log/slog"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func (n *NukiBle) Scan(timeout time.Duration) error {
	return n.ScanForDevice("", timeout)
}

func (n *NukiBle) ScanForDevice(deviceId string, timeout time.Duration) error {
	n.devices = map[string]bluetooth.ScanResult{}
	t := time.AfterFunc(timeout, func() { n.adapter.StopScan() })

	slog.Info("Scanning for devices...")
	err := n.adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) { n.onScan(a, sr, deviceId) })
	t.Stop()
	return err
}

func (n *NukiBle) onScan(a *bluetooth.Adapter, d bluetooth.ScanResult, stopOnDeviceId string) {
	if !strings.HasPrefix(d.LocalName(), "Nuki") {
		return
	}
	if _, exists := n.devices[d.Address.String()]; !exists {
		slog.Info("Found new device", "address", d.Address.String(), "rssi", d.RSSI, "name", d.LocalName())
		n.devices[d.Address.String()] = d
		if d.Address.String() == stopOnDeviceId {
			a.StopScan()
			return
		}
	}
}
