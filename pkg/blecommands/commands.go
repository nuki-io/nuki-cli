package blecommands

import (
	"encoding/binary"
	"fmt"
	"slices"
)

//go:generate stringer -type=CommandCode
type CommandCode uint16

const (
	CommandRequestData                 CommandCode = 0x0001
	CommandPublicKey                   CommandCode = 0x0003
	CommandChallenge                   CommandCode = 0x0004
	CommandAuthorizationAuthenticator  CommandCode = 0x0005
	CommandAuthorizationData           CommandCode = 0x0006
	CommandAuthorizationID             CommandCode = 0x0007
	CommandRemoveAuthorizationEntry    CommandCode = 0x0008
	CommandRequestAuthorizationEntries CommandCode = 0x0009
	CommandAuthorizationEntry          CommandCode = 0x000A
	CommandAuthorizationDataInvite     CommandCode = 0x000B
	CommandKeyturnerStates             CommandCode = 0x000C
	CommandLockAction                  CommandCode = 0x000D
	CommandStatus                      CommandCode = 0x000E
	CommandMostRecentCommand           CommandCode = 0x000F
	CommandOpeningsClosingsSummary     CommandCode = 0x0010
	CommandBatteryReport               CommandCode = 0x0011
	CommandErrorReport                 CommandCode = 0x0012
	CommandSetConfig                   CommandCode = 0x0013
	CommandRequestConfig               CommandCode = 0x0014
	CommandConfig                      CommandCode = 0x0015
	CommandSetSecurityPIN              CommandCode = 0x0019
	CommandRequestCalibration          CommandCode = 0x001A
	CommandRequestReboot               CommandCode = 0x001D
	CommandAuthorizationIDConfirmation CommandCode = 0x001E
	CommandAuthorizationIDInvite       CommandCode = 0x001F
	CommandVerifySecurityPIN           CommandCode = 0x0020
	CommandUpdateTime                  CommandCode = 0x0021
	CommandUpdateAuthorizationEntry    CommandCode = 0x0025
	CommandAuthorizationEntryCount     CommandCode = 0x0027
	CommandRequestLogEntries           CommandCode = 0x0031
	CommandLogEntry                    CommandCode = 0x0032
	CommandLogEntryCount               CommandCode = 0x0033
	CommandEnableLogging               CommandCode = 0x0034
	CommandSetAdvancedConfig           CommandCode = 0x0035
	CommandRequestAdvancedConfig       CommandCode = 0x0036
	CommandAdvancedConfig              CommandCode = 0x0037
	CommandAddTimeControlEntry         CommandCode = 0x0039
	CommandTimeControlEntryID          CommandCode = 0x003A
	CommandRemoveTimeControlEntry      CommandCode = 0x003B
	CommandRequestTimeControlEntries   CommandCode = 0x003C
	CommandTimeControlEntryCount       CommandCode = 0x003D
	CommandTimeControlEntry            CommandCode = 0x003E
	CommandUpdateTimeControlEntry      CommandCode = 0x003F
	CommandAddKeypadCode               CommandCode = 0x0041
	CommandKeypadCodeID                CommandCode = 0x0042
	CommandRequestKeypadCodes          CommandCode = 0x0043
	CommandKeypadCodeCount             CommandCode = 0x0044
	CommandKeypadCode                  CommandCode = 0x0045
	CommandUpdateKeypadCode            CommandCode = 0x0046
	CommandRemoveKeypadCode            CommandCode = 0x0047
	CommandAuthorizationInfo           CommandCode = 0x004C
	CommandSimpleLockAction            CommandCode = 0x0100
)

type Action uint8

const (
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

//go:generate stringer -type=StatusCode
type StatusCode uint8

const (
	StatusComplete StatusCode = 0x00
	StatusAccepted StatusCode = 0x01
)

var cmdImplMap = map[CommandCode]func() Command{
	CommandRequestData: func() Command { return Command(&RequestData{}) },
	CommandPublicKey:   func() Command { return Command(&PublicKey{}) },
	CommandChallenge:   func() Command { return Command(&Challenge{}) },
	// CommandAuthorizationAuthenticator:  func() Command { return Command(&AuthorizationAuthenticator{}) },
	// CommandAuthorizationData:           func() Command { return Command(&AuthorizationData{}) },
	CommandAuthorizationID: func() Command { return Command(&AuthorizationID{}) },
	// CommandRemoveAuthorizationEntry:    func() Command { return Command(&RemoveAuthorizationEntry{}) },
	// CommandRequestAuthorizationEntries: func() Command { return Command(&RequestAuthorizationEntries{}) },
	// CommandAuthorizationEntry:          func() Command { return Command(&AuthorizationEntry{}) },
	// CommandAuthorizationDataInvite:     func() Command { return Command(&AuthorizationDataInvite{}) },
	// CommandKeyturnerStates:             func() Command { return Command(&KeyturnerStates{}) },
	// CommandLockAction:                  func() Command { return Command(&LockAction{}) },
	CommandStatus: func() Command { return Command(&Status{}) },
	// CommandMostRecentCommand:           func() Command { return Command(&MostRecentCommand{}) },
	// CommandOpeningsClosingsSummary:     func() Command { return Command(&OpeningsClosingsSummary{}) },
	// CommandBatteryReport:               func() Command { return Command(&BatteryReport{}) },
	// CommandErrorReport:                 func() Command { return Command(&ErrorReport{}) },
	// CommandSetConfig:                   func() Command { return Command(&SetConfig{}) },
	// CommandRequestConfig:               func() Command { return Command(&RequestConfig{}) },
	// CommandConfig:                      func() Command { return Command(&Config{}) },
	// CommandSetSecurityPIN:              func() Command { return Command(&SetSecurityPIN{}) },
	// CommandRequestCalibration:          func() Command { return Command(&RequestCalibration{}) },
	// CommandRequestReboot:               func() Command { return Command(&RequestReboot{}) },
	// CommandAuthorizationIDConfirmation: func() Command { return Command(&AuthorizationIDConfirmation{}) },
	// CommandAuthorizationIDInvite:       func() Command { return Command(&AuthorizationIDInvite{}) },
	// CommandVerifySecurityPIN:           func() Command { return Command(&VerifySecurityPIN{}) },
	// CommandUpdateTime:                  func() Command { return Command(&UpdateTime{}) },
	// CommandUpdateAuthorizationEntry:    func() Command { return Command(&UpdateAuthorizationEntry{}) },
	// CommandAuthorizationEntryCount:     func() Command { return Command(&AuthorizationEntryCount{}) },
	// CommandRequestLogEntries:           func() Command { return Command(&RequestLogEntries{}) },
	// CommandLogEntry:                    func() Command { return Command(&LogEntry{}) },
	// CommandLogEntryCount:               func() Command { return Command(&LogEntryCount{}) },
	// CommandEnableLogging:               func() Command { return Command(&EnableLogging{}) },
	// CommandSetAdvancedConfig:           func() Command { return Command(&SetAdvancedConfig{}) },
	// CommandRequestAdvancedConfig:       func() Command { return Command(&RequestAdvancedConfig{}) },
	// CommandAdvancedConfig:              func() Command { return Command(&AdvancedConfig{}) },
	// CommandAddTimeControlEntry:         func() Command { return Command(&AddTimeControlEntry{}) },
	// CommandTimeControlEntryID:          func() Command { return Command(&TimeControlEntryID{}) },
	// CommandRemoveTimeControlEntry:      func() Command { return Command(&RemoveTimeControlEntry{}) },
	// CommandRequestTimeControlEntries:   func() Command { return Command(&RequestTimeControlEntries{}) },
	// CommandTimeControlEntryCount:       func() Command { return Command(&TimeControlEntryCount{}) },
	// CommandTimeControlEntry:            func() Command { return Command(&TimeControlEntry{}) },
	// CommandUpdateTimeControlEntry:      func() Command { return Command(&UpdateTimeControlEntry{}) },
	// CommandAddKeypadCode:               func() Command { return Command(&AddKeypadCode{}) },
	// CommandKeypadCodeID:                func() Command { return Command(&KeypadCodeID{}) },
	// CommandRequestKeypadCodes:          func() Command { return Command(&RequestKeypadCodes{}) },
	// CommandKeypadCodeCount:             func() Command { return Command(&KeypadCodeCount{}) },
	// CommandKeypadCode:                  func() Command { return Command(&KeypadCode{}) },
	// CommandUpdateKeypadCode:            func() Command { return Command(&UpdateKeypadCode{}) },
	// CommandRemoveKeypadCode:            func() Command { return Command(&RemoveKeypadCode{}) },
	// CommandAuthorizationInfo:           func() Command { return Command(&AuthorizationInfo{}) },
	// CommandSimpleLockAction:            func() Command { return Command(&SimpleLockAction{}) },
}

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
		command: CommandRequestData,
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
		command: CommandRequestData,
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

func ToEncryptedMessage(crypto Crypto, authId []byte, nonce []byte, c Command) []byte {
	payload := c.GetPayload()
	// length = authId + command + payload length + CRC
	pdata := make([]byte, 0, 4+2+len(payload)+2)
	pdata = append(pdata, authId...)
	pdata = binary.LittleEndian.AppendUint16(pdata, uint16(c.GetCommandCode()))
	pdata = append(pdata, payload...)
	pdata = binary.LittleEndian.AppendUint16(pdata, CRC(pdata))

	pdataEnc, _ := crypto.Encrypt(nonce, pdata)

	// length = nonce + authId + encrypted message length
	adata := make([]byte, 0, 24+4+2)
	adata = append(adata, nonce...)
	adata = append(adata, authId...)
	adata = binary.LittleEndian.AppendUint16(adata, uint16(len(pdataEnc)))
	return slices.Concat(adata, pdataEnc)
}

type Command interface {
	FromMessage([]byte) error
	GetCommandCode() CommandCode
	GetPayload() []byte
}

type RequestData struct {
	CommandIdentifier CommandCode
	// AdditionalData    []byte
}

func (c *RequestData) GetCommandCode() CommandCode {
	return CommandRequestData
}

func (c *RequestData) FromMessage([]byte) error {
	return nil
}

func (c *RequestData) GetPayload() []byte {
	payload := make([]byte, 2)
	binary.LittleEndian.PutUint16(payload, uint16(c.CommandIdentifier))
	return payload
}

type PublicKey struct {
	PublicKey []byte
}

func (c *PublicKey) GetCommandCode() CommandCode {
	return CommandPublicKey
}
func (c *PublicKey) FromMessage(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("public key length must be more than 0")
	}
	c.PublicKey = b
	return nil
}
func (c *PublicKey) GetPayload() []byte {
	return c.PublicKey
}

type AuthorizationAuthenticator struct {
	Authenticator []byte
}

func (c *AuthorizationAuthenticator) GetCommandCode() CommandCode {
	return CommandAuthorizationAuthenticator
}
func (c *AuthorizationAuthenticator) FromMessage(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("authenticator length must be more than 0")
	}
	c.Authenticator = b
	return nil
}
func (c *AuthorizationAuthenticator) GetPayload() []byte {
	return c.Authenticator
}

type AuthorizationData struct {
	// TODO: should use more concrete types
	Authenticator []byte
	IdType        uint8 // 0x00 = App, 0x01 = Bridge, 0x02 = Fob, 0x03 = Keypad
	Id            []byte
	Name          string
	Nonce         []byte
}

func (c *AuthorizationData) GetCommandCode() CommandCode {
	return CommandAuthorizationData
}
func (c *AuthorizationData) FromMessage(b []byte) error {
	return fmt.Errorf("not implemented")
}
func (c *AuthorizationData) GetPayload() []byte {
	appName := [32]byte{}
	copy(appName[:], c.Name)
	return slices.Concat(
		c.Authenticator,
		[]byte{c.IdType},
		c.Id,
		appName[:],
		c.Nonce,
	)
}

type AuthorizationIDConfirmation struct {
	Authenticator []byte
	AuthId        []byte
}

func (c *AuthorizationIDConfirmation) GetCommandCode() CommandCode {
	return CommandAuthorizationIDConfirmation
}
func (c *AuthorizationIDConfirmation) FromMessage(b []byte) error {
	return fmt.Errorf("not implemented")
}
func (c *AuthorizationIDConfirmation) GetPayload() []byte {
	return slices.Concat(c.Authenticator, c.AuthId)
}

type Challenge struct {
	Nonce []byte
}

func (c *Challenge) GetCommandCode() CommandCode {
	return CommandChallenge
}
func (c *Challenge) FromMessage(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("challenge length must be exactly 32 bytes, got: %d", len(b))
	}
	c.Nonce = b
	return nil
}
func (c *Challenge) GetPayload() []byte {
	return c.Nonce
}

type AuthorizationID struct {
	Authenticator []byte
	AuthId        []byte
	Uuid          []byte
	Nonce         []byte
}

func (c *AuthorizationID) GetCommandCode() CommandCode {
	return CommandAuthorizationID
}
func (c *AuthorizationID) FromMessage(b []byte) error {
	if len(b) != 84 {
		return fmt.Errorf("authorization ID length must be exactly 84 bytes, got: %d", len(b))
	}
	c.Authenticator = b[:32]
	c.AuthId = b[32:36]
	c.Uuid = b[36:52]
	c.Nonce = b[52:]
	return nil
}
func (c *AuthorizationID) GetPayload() []byte {
	return slices.Concat(c.Authenticator, c.AuthId, c.Uuid, c.Nonce)
}

type Status struct {
	Status StatusCode
}

func (c *Status) GetCommandCode() CommandCode {
	return CommandStatus
}
func (c *Status) FromMessage(b []byte) error {
	if len(b) != 1 {
		return fmt.Errorf("status length must be exactly 1 byte, got: %d", len(b))
	}
	c.Status = StatusCode(b[0])
	return nil
}
func (c *Status) GetPayload() []byte {
	return []byte{byte(c.Status)}
}
