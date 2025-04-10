package nukible

import (
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

type NukiBle struct {
	adapter *bluetooth.Adapter
	devices map[bluetooth.Address]bluetooth.ScanResult
}

func NewNukiBle() (*NukiBle, error) {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()

	if err != nil {
		return nil, err
	}
	return &NukiBle{
		adapter: adapter,
	}, nil
}

func (n *NukiBle) Scan(timeout time.Duration) {
	n.devices = map[bluetooth.Address]bluetooth.ScanResult{}
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
	if _, exists := n.devices[d.Address]; !exists {
		println("found device:", d.Address.String(), d.RSSI, d.LocalName(), d.AdvertisementPayload)
		n.devices[d.Address] = d
	}
}
