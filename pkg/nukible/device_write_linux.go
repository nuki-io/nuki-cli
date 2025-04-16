package nukible

import (
	"tinygo.org/x/bluetooth"
)

func (n *Device) osWrite(char bluetooth.DeviceCharacteristic, data []byte) {
	char.WriteWithoutResponse(data)
}
