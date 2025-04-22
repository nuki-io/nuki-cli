package nukible

import (
	"log/slog"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

type NukiBle struct {
	adapter *bluetooth.Adapter
	devices map[string]bluetooth.ScanResult
}

func NewNukiBle() (*NukiBle, error) {
	adapter := bluetooth.DefaultAdapter
	err := adapter.Enable()

	if err != nil {
		return nil, err
	}
	return &NukiBle{
		adapter: adapter,
		devices: map[string]bluetooth.ScanResult{},
	}, nil
}

func (n *NukiBle) GetDevices() map[string]bluetooth.ScanResult {
	return n.devices
}

func (n *NukiBle) GetDeviceAddress(deviceId string) (res *bluetooth.Address, ok bool) {
	d, exists := n.devices[deviceId]

	if !exists {
		return osGetUndiscoveredDeviceAddress(deviceId)
	}
	return &d.Address, true
}

func (n *NukiBle) Connect(addr bluetooth.Address) (*Device, error) {
	device, err := n.adapter.Connect(addr, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(5 * time.Second),
		MinInterval:       bluetooth.NewDuration(15 * time.Millisecond),
		MaxInterval:       bluetooth.NewDuration(15 * time.Millisecond),
		Timeout:           bluetooth.NewDuration(6 * time.Second),
	})
	if err != nil {
		return nil, err
	}
	return &Device{
		btDev: device,
	}, nil
}

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
