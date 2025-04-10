package nukible

import (
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func (n *NukiBle) Scan(timeout time.Duration) {
	n.devices = map[string]bluetooth.ScanResult{}
	time.AfterFunc(timeout, func() { n.adapter.StopScan() })
	// Start scanning.
	println("scanning...")
	err := n.adapter.Scan(n.onScan)
	if err != nil {
		panic("Failed to start device scan")
	}
}

func (n *NukiBle) onScan(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	if !strings.HasPrefix(d.LocalName(), "Nuki") {
		return
	}
	if _, exists := n.devices[d.Address.String()]; !exists {
		println("Found new device:", d.Address.String(), d.RSSI, d.LocalName(), d.AdvertisementPayload)
		n.devices[d.Address.String()] = d
	}
}
