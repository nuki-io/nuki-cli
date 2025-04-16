package nukible

import (
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func (n *NukiBle) Scan(timeout time.Duration) {
	n.ScanForDevice("", timeout)
}

func (n *NukiBle) ScanForDevice(deviceId string, timeout time.Duration) {
	n.devices = map[string]bluetooth.ScanResult{}
	t := time.AfterFunc(timeout, func() { n.adapter.StopScan() })
	// Start scanning.
	println("scanning...")
	err := n.adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) { n.onScan(a, sr, deviceId) })
	t.Stop()
	if err != nil {
		panic("Failed to start device scan")
	}
}

func (n *NukiBle) onScan(a *bluetooth.Adapter, d bluetooth.ScanResult, stopOnDeviceId string) {
	if !strings.HasPrefix(d.LocalName(), "Nuki") {
		return
	}
	if _, exists := n.devices[d.Address.String()]; !exists {
		println("Found new device:", d.Address.String(), d.RSSI, d.LocalName(), d.AdvertisementPayload)
		n.devices[d.Address.String()] = d
		if d.Address.String() == stopOnDeviceId {
			a.StopScan()
			return
		}
	}
}
