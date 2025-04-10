package nukible

import "tinygo.org/x/bluetooth"

var (
	// https://developer.nuki.io/t/bluetooth-api/27
	baseUuid = bluetooth.NewUUID([16]byte{
		0xa9, 0x2e, 0x00, 0x00,
		0x55, 0x01,
		0x11, 0xe4,
		0x91, 0x6c,
		0x08, 0x00, 0x20, 0x0c, 0x9a, 0x66})

	KeyturnerInitializationService = baseUuid.Replace16BitComponent(0xe000)

	KeyturnerPairingService            = baseUuid.Replace16BitComponent(0xe100)
	KeyturnerPairingServiceUltra       = baseUuid.Replace16BitComponent(0xe300)
	KeyturnerPairingGdioCharacteristic = baseUuid.Replace16BitComponent(0xe101)

	KeyturnerService             = baseUuid.Replace16BitComponent(0xe200)
	KeyturnerGdioCharacteristic  = baseUuid.Replace16BitComponent(0xe201)
	KeyturnerUsdioCharacteristic = baseUuid.Replace16BitComponent(0xe202)
)
