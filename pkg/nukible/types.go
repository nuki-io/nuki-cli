package nukible

import "tinygo.org/x/bluetooth"

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
