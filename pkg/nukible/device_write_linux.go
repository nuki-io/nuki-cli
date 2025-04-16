package nukible

import "fmt"

func (n *Device) Write(data []byte) []byte {
	sem := make(chan int, 1)
	sem <- 1

	fmt.Printf("Writing bytes to characteristic %x\n", data)
	rxData := make([]byte, 0)
	n.pairingGdioChar.EnableNotifications(func(buf []byte) { rxData = onGdioNotify(buf, sem) })
	n.pairingGdioChar.WriteWithoutResponse(data)

	fmt.Println("Waiting for response")
	sem <- 1
	// disable notifications again - TODO: sensible, or should we just enable it once?
	n.pairingGdioChar.EnableNotifications(nil)

	return rxData
}
