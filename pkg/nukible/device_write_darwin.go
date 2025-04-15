package nukible

import "fmt"

func (n *Device) Write(data []byte) []byte {
	sem := make(chan int, 1)
	sem <- 1

	rxData := make([]byte, 0)
	n.pairingGdioChar.EnableNotifications(func(buf []byte) { rxData = onGdioNotify(buf, sem) })
	n.pairingGdioChar.Write(data)
	n.pairingGdioChar.EnableNotifications(nil)

	fmt.Println("Waiting for response")
	sem <- 1

	return rxData
}
