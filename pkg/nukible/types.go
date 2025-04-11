package nukible

import (
	"encoding/binary"

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

type Command uint16

const (
	RequestData                 Command = 0x0001
	PublicKey                   Command = 0x0003
	Challenge                   Command = 0x0004
	AuthorizationAuthenticator  Command = 0x0005
	AuthorizationData           Command = 0x0006
	AuthorizationID             Command = 0x0007
	RemoveAuthorizationEntry    Command = 0x0008
	RequestAuthorizationEntries Command = 0x0009
	AuthorizationEntry          Command = 0x000A
	AuthorizationDataInvite     Command = 0x000B
	KeyturnerStates             Command = 0x000C
	LockAction                  Command = 0x000D
	Status                      Command = 0x000E
	MostRecentCommand           Command = 0x000F
	OpeningsClosingsSummary     Command = 0x0010
	BatteryReport               Command = 0x0011
	ErrorReport                 Command = 0x0012
	SetConfig                   Command = 0x0013
	RequestConfig               Command = 0x0014
	Config                      Command = 0x0015
	SetSecurityPIN              Command = 0x0019
	RequestCalibration          Command = 0x001A
	RequestReboot               Command = 0x001D
	AuthorizationIDConfirmation Command = 0x001E
	AuthorizationIDInvite       Command = 0x001F
	VerifySecurityPIN           Command = 0x0020
	UpdateTime                  Command = 0x0021
	UpdateAuthorizationEntry    Command = 0x0025
	AuthorizationEntryCount     Command = 0x0027
	RequestLogEntries           Command = 0x0031
	LogEntry                    Command = 0x0032
	LogEntryCount               Command = 0x0033
	EnableLogging               Command = 0x0034
	SetAdvancedConfig           Command = 0x0035
	RequestAdvancedConfig       Command = 0x0036
	AdvancedConfig              Command = 0x0037
	AddTimeControlEntry         Command = 0x0039
	TimeControlEntryID          Command = 0x003A
	RemoveTimeControlEntry      Command = 0x003B
	RequestTimeControlEntries   Command = 0x003C
	TimeControlEntryCount       Command = 0x003D
	TimeControlEntry            Command = 0x003E
	UpdateTimeControlEntry      Command = 0x003F
	AddKeypadCode               Command = 0x0041
	KeypadCodeID                Command = 0x0042
	RequestKeypadCodes          Command = 0x0043
	KeypadCodeCount             Command = 0x0044
	KeypadCode                  Command = 0x0045
	UpdateKeypadCode            Command = 0x0046
	RemoveKeypadCode            Command = 0x0047
	AuthorizationInfo           Command = 0x004C
	SimpleLockAction            Command = 0x0100
)

type UnencryptedCommand struct {
	command Command
	payload []byte
}

func NewUnencryptedCommand(cmd Command, payload []byte) UnencryptedCommand {
	return UnencryptedCommand{
		command: cmd,
		payload: payload,
	}
}

func (c *UnencryptedCommand) ToMessage() []byte {
	res := make([]byte, 2+len(c.payload))
	binary.BigEndian.PutUint16(res, uint16(c.command))
	for i := 0; i < len(c.payload); i++ {
		res[i+2] = c.payload[i]
	}
	res = binary.BigEndian.AppendUint16(res, CRC(res))
	return res
}
