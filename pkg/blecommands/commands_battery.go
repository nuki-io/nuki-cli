package blecommands

import (
	"encoding/binary"
	"fmt"
)

// BatteryReport (0x0011)

var _ Response = &BatteryReport{}

type BatteryReport struct {
	BatteryDrain      uint16 `json:"batteryDrain"`   // mWs
	BatteryVoltage    uint16 `json:"batteryVoltage"` // mV
	CriticalState     bool   `json:"criticalState"`
	LockAction        Action `json:"lockAction"`
	StartVoltage      uint16 `json:"startVoltage"`     // mV
	LowestVoltage     uint16 `json:"lowestVoltage"`    // mV
	LockDistance      uint16 `json:"lockDistance"`     // degrees
	StartTemperature  int8   `json:"startTemperature"` // °C
	MaxTurnCurrent    uint16 `json:"maxTurnCurrent"`
	BatteryResistance uint16 `json:"batteryResistance"`
}

func (b *BatteryReport) GetCommandCode() CommandCode { return CommandBatteryReport }

func (b *BatteryReport) FromMessage(data []byte) error {
	if len(data) < 17 {
		return fmt.Errorf("battery report too short: %d bytes, need at least 17", len(data))
	}
	b.BatteryDrain = binary.LittleEndian.Uint16(data[0:2])
	b.BatteryVoltage = binary.LittleEndian.Uint16(data[2:4])
	b.CriticalState = data[4] != 0
	b.LockAction = Action(data[5])
	b.StartVoltage = binary.LittleEndian.Uint16(data[6:8])
	b.LowestVoltage = binary.LittleEndian.Uint16(data[8:10])
	b.LockDistance = binary.LittleEndian.Uint16(data[10:12])
	b.StartTemperature = int8(data[12])
	b.MaxTurnCurrent = binary.LittleEndian.Uint16(data[13:15])
	b.BatteryResistance = binary.LittleEndian.Uint16(data[15:17])
	return nil
}
