package nukible

import "tinygo.org/x/bluetooth"

// osGetUndiscoveredDeviceAddress on Darwin will construct a bluetooth.Address from
// the given id.
func osGetUndiscoveredDeviceAddress(id string) (res *bluetooth.Address, ok bool) {
	uuid, _ := bluetooth.ParseUUID(id)
	addr := &bluetooth.Address{
		UUID: uuid,
	}
	return addr, true
}
