package blecommands

import (
	"encoding/binary"
	"fmt"
	"slices"
	"time"
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
	CommandKeyturnerStates: func() Command { return Command(&KeyturnerStates{}) },
	// CommandLockAction:                  func() Command { return Command(&LockAction{}) },
	CommandStatus: func() Command { return Command(&Status{}) },
	// CommandMostRecentCommand:           func() Command { return Command(&MostRecentCommand{}) },
	// CommandOpeningsClosingsSummary:     func() Command { return Command(&OpeningsClosingsSummary{}) },
	// CommandBatteryReport:               func() Command { return Command(&BatteryReport{}) },
	CommandErrorReport: func() Command { return Command(&ErrorReport{}) },
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
	return fmt.Errorf("not implemented")
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

type ErrorReport struct {
	ErrorCode         int8
	CommandIdentifier CommandCode
}

func (c *ErrorReport) GetCommandCode() CommandCode {
	return CommandErrorReport
}
func (c *ErrorReport) FromMessage(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("error report length must be exactly 3 bytes, got: %d", len(b))
	}
	c.ErrorCode = int8(b[0])
	c.CommandIdentifier = CommandCode(binary.LittleEndian.Uint16(b[1:]))
	return nil
}
func (c *ErrorReport) GetPayload() []byte {
	return []byte{byte(c.ErrorCode), byte(c.CommandIdentifier), byte(c.CommandIdentifier >> 8)}
}

type LockAction struct {
	Action Action
	AppId  []byte
	Flags  byte
	Nonce  []byte
}

func (c *LockAction) GetCommandCode() CommandCode {
	return CommandLockAction
}
func (c *LockAction) FromMessage(b []byte) error {
	return fmt.Errorf("not implemented")
}
func (c *LockAction) GetPayload() []byte {
	return slices.Concat(
		[]byte{byte(c.Action)},
		c.AppId,
		[]byte{c.Flags},
		c.Nonce,
	)
}

type KeyturnerStates struct {
	NukiState                      byte
	LockState                      byte
	Trigger                        byte
	CurrentTime                    time.Time
	TimezoneOffset                 byte
	CriticalBatteryState           byte
	ConfigUpdateCount              byte
	LockNGoTimer                   byte
	LastLockAction                 byte
	LastLockActionTrigger          byte
	LastLockActionCompletionStatus byte
	DoorSensorState                byte
	NightmodeActive                byte
	AccessoryBatteryState          byte
	RemoteAccessStatus             byte
	BleConnectionStrength          int8
	WifiConnectionStrength         int8
	WifiConnectionStatus           byte
	MqttConnectionStatus           byte
	ThreadConnectionStatus         byte
}

func (c *KeyturnerStates) GetCommandCode() CommandCode {
	return CommandKeyturnerStates
}
func (c *KeyturnerStates) FromMessage(b []byte) error {
	if len(b) != 27 { // TODO: which one is the correct length?
		return fmt.Errorf("keyturner states length must be exactly 27 bytes, got: %d", len(b))
	}
	c.NukiState = b[0]
	c.LockState = b[1]
	c.Trigger = b[2]
	c.CurrentTime = time.Date(
		int(binary.LittleEndian.Uint16(b[3:5])),
		time.Month(b[5]),
		int(b[6]),
		int(b[7]),
		int(b[8]),
		int(b[9]),
		0,
		time.UTC,
	)
	c.TimezoneOffset = b[10]
	c.CriticalBatteryState = b[11]
	c.ConfigUpdateCount = b[12]
	c.LockNGoTimer = b[13]
	c.LastLockAction = b[14]
	c.LastLockActionTrigger = b[15]
	c.LastLockActionCompletionStatus = b[16]
	c.DoorSensorState = b[17]
	c.NightmodeActive = b[18]
	c.AccessoryBatteryState = b[19]
	c.RemoteAccessStatus = b[20]
	c.BleConnectionStrength = int8(b[21])
	c.WifiConnectionStrength = int8(b[22])
	c.WifiConnectionStatus = b[23]
	c.MqttConnectionStatus = b[24]
	c.ThreadConnectionStatus = b[25]
	return nil
}
func (c *KeyturnerStates) GetPayload() []byte {
	panic("not implemented")
}
