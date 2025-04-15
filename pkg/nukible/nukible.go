package nukible

import (
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
	}, nil
}

func (n *NukiBle) GetDevices() map[string]bluetooth.ScanResult {
	return n.devices
}

func (n *NukiBle) GetDeviceAddress(deviceId string) (res *bluetooth.Address, ok bool) {
	d, exists := n.devices[deviceId]

	if !exists {
		return nil, false
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
