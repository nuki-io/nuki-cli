package blecommands

import (
	"encoding/binary"
	"slices"
)

type CommandCode uint16

const (
	RequestData                 CommandCode = 0x0001
	PublicKey                   CommandCode = 0x0003
	Challenge                   CommandCode = 0x0004
	AuthorizationAuthenticator  CommandCode = 0x0005
	AuthorizationData           CommandCode = 0x0006
	AuthorizationID             CommandCode = 0x0007
	RemoveAuthorizationEntry    CommandCode = 0x0008
	RequestAuthorizationEntries CommandCode = 0x0009
	AuthorizationEntry          CommandCode = 0x000A
	AuthorizationDataInvite     CommandCode = 0x000B
	KeyturnerStates             CommandCode = 0x000C
	LockAction                  CommandCode = 0x000D
	Status                      CommandCode = 0x000E
	MostRecentCommand           CommandCode = 0x000F
	OpeningsClosingsSummary     CommandCode = 0x0010
	BatteryReport               CommandCode = 0x0011
	ErrorReport                 CommandCode = 0x0012
	SetConfig                   CommandCode = 0x0013
	RequestConfig               CommandCode = 0x0014
	Config                      CommandCode = 0x0015
	SetSecurityPIN              CommandCode = 0x0019
	RequestCalibration          CommandCode = 0x001A
	RequestReboot               CommandCode = 0x001D
	AuthorizationIDConfirmation CommandCode = 0x001E
	AuthorizationIDInvite       CommandCode = 0x001F
	VerifySecurityPIN           CommandCode = 0x0020
	UpdateTime                  CommandCode = 0x0021
	UpdateAuthorizationEntry    CommandCode = 0x0025
	AuthorizationEntryCount     CommandCode = 0x0027
	RequestLogEntries           CommandCode = 0x0031
	LogEntry                    CommandCode = 0x0032
	LogEntryCount               CommandCode = 0x0033
	EnableLogging               CommandCode = 0x0034
	SetAdvancedConfig           CommandCode = 0x0035
	RequestAdvancedConfig       CommandCode = 0x0036
	AdvancedConfig              CommandCode = 0x0037
	AddTimeControlEntry         CommandCode = 0x0039
	TimeControlEntryID          CommandCode = 0x003A
	RemoveTimeControlEntry      CommandCode = 0x003B
	RequestTimeControlEntries   CommandCode = 0x003C
	TimeControlEntryCount       CommandCode = 0x003D
	TimeControlEntry            CommandCode = 0x003E
	UpdateTimeControlEntry      CommandCode = 0x003F
	AddKeypadCode               CommandCode = 0x0041
	KeypadCodeID                CommandCode = 0x0042
	RequestKeypadCodes          CommandCode = 0x0043
	KeypadCodeCount             CommandCode = 0x0044
	KeypadCode                  CommandCode = 0x0045
	UpdateKeypadCode            CommandCode = 0x0046
	RemoveKeypadCode            CommandCode = 0x0047
	AuthorizationInfo           CommandCode = 0x004C
	SimpleLockAction            CommandCode = 0x0100
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
	command CommandCode
	payload []byte
}

func NewUnencryptedCommand(cmd CommandCode, payload []byte) UnencryptedCommand {
	return UnencryptedCommand{
		command: cmd,
		payload: payload,
	}
}

func NewUnencryptedRequestData(request CommandCode) UnencryptedCommand {
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
	command CommandCode
	payload []byte
}

func NewEncryptedCommand(crypto Crypto, authId []byte, cmd CommandCode, payload []byte) EncryptedCommand {
	return EncryptedCommand{
		crypto:  crypto,
		authId:  authId,
		command: cmd,
		payload: payload,
	}
}

func NewEncryptedRequestData(crypto Crypto, authId []byte, request CommandCode) EncryptedCommand {
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
	// length = authId + command + payload length + CRC
	pdata := make([]byte, 0, 4+2+len(c.payload)+2)
	pdata = append(pdata, c.authId...)
	pdata = binary.LittleEndian.AppendUint16(pdata, uint16(c.command))
	pdata = append(pdata, c.payload...)
	pdata = binary.LittleEndian.AppendUint16(pdata, CRC(pdata))

	pdataEnc, _ := c.crypto.Encrypt(nonce, pdata)

	// length = nonce + authId + encrypted message length
	adata := make([]byte, 0, 24+4+2)
	adata = append(adata, nonce...)
	adata = append(adata, c.authId...)
	adata = binary.LittleEndian.AppendUint16(adata, uint16(len(pdataEnc)))
	return slices.Concat(adata, pdataEnc)
}
