package nukible

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
)

func (n *NukiBle) Authorize(mac string) error {
	devs := n.GetDevices()
	var addr bluetooth.Address
	d, exists := devs[mac]

	if !exists {
		return fmt.Errorf("requested device with MAC %s was not discovered", mac)
	}
	addr = d.Address

	device, err := n.adapter.Connect(addr, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(5 * time.Second),
		MinInterval:       bluetooth.NewDuration(15 * time.Millisecond),
		MaxInterval:       bluetooth.NewDuration(15 * time.Millisecond),
		Timeout:           bluetooth.NewDuration(6 * time.Second),
	})
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to device %s. %s", mac, err.Error()))
	}

	services, err := device.DiscoverServices([]bluetooth.UUID{KeyturnerPairingService})
	if len(services) == 0 && err != nil {
		// expected, maybe it's an Ultra
		services, err = device.DiscoverServices([]bluetooth.UUID{KeyturnerPairingServiceUltra})
	}
	if len(services) == 0 && err != nil {
		panic(fmt.Sprintf("Could not discover any pairing services. %s", err.Error()))
	}
	if len(services) != 1 {
		panic(fmt.Sprintf("Expected exactly one pairing service, got %d", len(services)))
	}

	fmt.Println("Service", services[0].String())
	characteristics, err := services[0].DiscoverCharacteristics([]bluetooth.UUID{KeyturnerPairingGdioCharacteristic})
	if err != nil {
		device.Disconnect()
		panic(fmt.Sprintf("Could not discover GDIO characteristic. %s", err.Error()))
	}
	if len(characteristics) != 1 {
		panic(fmt.Sprintf("Expected exactly one GDIO characteristic, got %d", len(characteristics)))
	}

	gdio := characteristics[0]
	sem := make(chan int, 1)
	sem <- 1
	gdio.EnableNotifications(func(buf []byte) { onGdioNotify(buf, sem) })
	fmt.Println("Writing request for key exchange")
	gdio.WriteWithoutResponse([]byte{0x01, 0x00, 0x03, 0x00, 0x27, 0xA7})
	fmt.Println("Waiting for response")
	sem <- 1

	fmt.Println("Disconnecting...")
	device.Disconnect()
	return nil
}

func onGdioNotify(buf []byte, sem chan int) {
	fmt.Println("Received response")

	crc := CRC(buf[:len(buf)-2])
	fmt.Printf("Full response: %x\n", buf)
	fmt.Printf("Without CRC: %x\n", buf[:len(buf)-2])
	fmt.Printf("Own CRC calc: %x\n", crc)
	<-sem
}
