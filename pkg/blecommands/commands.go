package blecommands

import (
	"encoding/binary"
	"slices"
)

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

type Action uint8

var (
	Unlock Action = 0x01
	Lock   Action = 0x02

	Unlatch          Action = 0x03
	LockAndGo        Action = 0x04
	LockAndGoUnlatch Action = 0x05
	FullLock         Action = 0x06

	FobAction1 Action = 0x81
	FobAction2 Action = 0x82
	FobAction3 Action = 0x83
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

func NewUnencryptedRequestData(request Command) UnencryptedCommand {
	payload := make([]byte, 2)
	binary.LittleEndian.PutUint16(payload, uint16(request))
	return UnencryptedCommand{
		command: RequestData,
		payload: payload,
	}
}

func (c *UnencryptedCommand) ToMessage() []byte {
	res := make([]byte, 2+len(c.payload))
	binary.LittleEndian.PutUint16(res, uint16(c.command))
	for i, x := range c.payload {
		res[i+2] = x
	}
	res = binary.LittleEndian.AppendUint16(res, CRC(res))
	return res
}

type EncryptedCommand struct {
	crypto  Crypto
	authId  []byte
	command Command
	payload []byte
}

func NewEncryptedCommand(crypto Crypto, authId []byte, cmd Command, payload []byte) EncryptedCommand {
	return EncryptedCommand{
		crypto:  crypto,
		authId:  authId,
		command: cmd,
		payload: payload,
	}
}

func NewEncryptedRequestData(crypto Crypto, authId []byte, request Command) EncryptedCommand {
	payload := make([]byte, 2)
	binary.LittleEndian.PutUint16(payload, uint16(request))
	return EncryptedCommand{
		crypto:  crypto,
		authId:  authId,
		command: RequestData,
		payload: payload,
	}
}

func (c *EncryptedCommand) ToMessage(nonce []byte) []byte {
	pdata := make([]byte, 0, 4+2+len(c.payload)+2)
	pdata = append(pdata, c.authId...)
	pdata = binary.LittleEndian.AppendUint16(pdata, uint16(c.command))
	pdata = append(pdata, c.payload...)
	pdata = binary.LittleEndian.AppendUint16(pdata, CRC(pdata))

	pdataEnc, _ := c.crypto.Encrypt(nonce, pdata)

	adata := make([]byte, 0, 24+4+2)
	adata = append(adata, nonce...)
	adata = append(adata, c.authId...)
	adata = binary.LittleEndian.AppendUint16(adata, uint16(len(pdataEnc)))
	return slices.Concat(adata, pdataEnc)
}
