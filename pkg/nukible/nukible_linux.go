package nukible

import "tinygo.org/x/bluetooth"

// osGetUndiscoveredDeviceAddress on Linux will never return a device address
// because the device must be discovered beforehand with a scan in order to connect.
func osGetUndiscoveredDeviceAddress(id string) (res *bluetooth.Address, ok bool) {
	return nil, false
}
