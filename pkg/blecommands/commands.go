package blecommands

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strings"
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
	CommandAuthorizationData5G         CommandCode = 0x0006
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

//go:generate stringer -type=Action
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

//go:generate stringer -type=LockState
type LockState uint8

const (
	LockStateUncalibrated    LockState = 0x00
	LockStateLocked          LockState = 0x01
	LockStateUnlocking       LockState = 0x02
	LockStateUnlocked        LockState = 0x03
	LockStateLocking         LockState = 0x04
	LockStateUnlatched       LockState = 0x05
	LockStateUnlockedLockNGo LockState = 0x06
	LockStateUnlatching      LockState = 0x07
	LockStateCalibration     LockState = 0xFC
	LockStateBootRun         LockState = 0xFD
	LockStateMotorBlocked    LockState = 0xFE
	LockStateUndefined       LockState = 0xFF
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
	CommandConfig: func() Command { return Command(&Config{}) },
	// CommandSetSecurityPIN:              func() Command { return Command(&SetSecurityPIN{}) },
	// CommandRequestCalibration:          func() Command { return Command(&RequestCalibration{}) },
	// CommandRequestReboot:               func() Command { return Command(&RequestReboot{}) },
	// CommandAuthorizationIDConfirmation: func() Command { return Command(&AuthorizationIDConfirmation{}) },
	// CommandAuthorizationIDInvite:       func() Command { return Command(&AuthorizationIDInvite{}) },
	// CommandVerifySecurityPIN:           func() Command { return Command(&VerifySecurityPIN{}) },
	// CommandUpdateTime:                  func() Command { return Command(&UpdateTime{}) },
	// CommandUpdateAuthorizationEntry:    func() Command { return Command(&UpdateAuthorizationEntry{}) },
	// CommandAuthorizationEntryCount:     func() Command { return Command(&AuthorizationEntryCount{}) },
	CommandRequestLogEntries: func() Command { return Command(&RequestLogEntries{}) },
	CommandLogEntry:          func() Command { return Command(&LogEntry{}) },
	// CommandLogEntryCount:     func() Command { return Command(&LogEntryCount{}) },
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
	CommandAuthorizationInfo: func() Command { return Command(&AuthorizationInfo{}) },
	// CommandSimpleLockAction:            func() Command { return Command(&SimpleLockAction{}) },
}
var timezoneMap = map[uint16]string{
	0:     "Africa/Cairo",
	1:     "Africa/Lagos",
	2:     "Africa/Maputo",
	3:     "Africa/Nairobi",
	4:     "America/Anchorage",
	5:     "America/Argentina/Buenos_Aires",
	6:     "America/Chicago",
	7:     "America/Denver",
	8:     "America/Halifax",
	9:     "America/Los_Angeles",
	10:    "America/Manaus",
	11:    "America/Mexico_City",
	12:    "America/New_York",
	13:    "America/Phoenix",
	14:    "America/Regina",
	15:    "America/Santiago",
	16:    "America/Sao_Paulo",
	17:    "America/St_Johns",
	18:    "Asia/Bangkok",
	19:    "Asia/Dubai",
	20:    "Asia/Hong_Kong",
	21:    "Asia/Jerusalem",
	22:    "Asia/Karachi",
	23:    "Asia/Kathmandu",
	24:    "Asia/Kolkata",
	25:    "Asia/Riyadh",
	26:    "Asia/Seoul",
	27:    "Asia/Shanghai",
	28:    "Asia/Tehran",
	29:    "Asia/Tokyo",
	30:    "Asia/Yangon",
	31:    "Australia/Adelaide",
	32:    "Australia/Brisbane",
	33:    "Australia/Darwin",
	34:    "Australia/Hobart",
	35:    "Australia/Perth",
	36:    "Australia/Sydney",
	37:    "Europe/Berlin",
	38:    "Europe/Helsinki",
	39:    "Europe/Istanbul",
	40:    "Europe/London",
	41:    "Europe/Moscow",
	42:    "Pacific/Auckland",
	43:    "Pacific/Guam",
	44:    "Pacific/Honolulu",
	45:    "Pacific/Pago_Pago",
	65535: "", // Special case:  No timezone
}

func byteToBool(b byte) bool {
	return b != 0
}

type Command interface {
	GetCommandCode() CommandCode
}

type Request interface {
	Command
	GetPayload() []byte
}

type Response interface {
	Command
	FromMessage([]byte) error
}

type RequestData struct {
	CommandIdentifier CommandCode
	// AdditionalData    []byte
}

func (c *RequestData) GetCommandCode() CommandCode {
	return CommandRequestData
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

type AuthorizationType uint8

const (
	AuthorizationTypeApp    AuthorizationType = 0x00 // App
	AuthorizationTypeBridge AuthorizationType = 0x01 // Bridge
	AuthorizationTypeFob    AuthorizationType = 0x02 // Fob
	AuthorizationTypeKeypad AuthorizationType = 0x03 // Keypad
)

type AuthorizationData struct {
	Authenticator []byte
	IdType        AuthorizationType
	Id            []byte
	Name          string
	Nonce         []byte
}

func (c *AuthorizationData) GetCommandCode() CommandCode {
	return CommandAuthorizationData
}
func (c *AuthorizationData) GetPayload() []byte {
	appName := [32]byte{}
	copy(appName[:], c.Name)
	return slices.Concat(
		c.Authenticator,
		[]byte{byte(c.IdType)},
		c.Id,
		appName[:],
		c.Nonce,
	)
}

var _ Command = &AuthorizationData5G{}

type AuthorizationData5G struct {
	Id          []byte
	Name        string
	SecurityPin Pin
}

func (a *AuthorizationData5G) GetCommandCode() CommandCode {
	return CommandAuthorizationData5G
}

func (a *AuthorizationData5G) GetPayload() []byte {
	appName := [32]byte{}
	copy(appName[:], a.Name)
	return slices.Concat(
		a.Id,
		appName[:],
		a.SecurityPin.GetPinBytes(),
	)
}

type AuthorizationIDConfirmation struct {
	Authenticator []byte
	AuthId        []byte
}

func (c *AuthorizationIDConfirmation) GetCommandCode() CommandCode {
	return CommandAuthorizationIDConfirmation
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

// TODO: probably better to split it into 5G and pre-5G versions
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
	if len(b) != 84 && len(b) != 20 { // 20 bytes from 5G onwards
		return fmt.Errorf("authorization ID length must be exactly 84 bytes (until 5G) or 20 bytes (from 5G onwards), got: %d", len(b))
	}
	if len(b) == 20 {
		// 5G onwards
		c.AuthId = b[:4]
		c.Uuid = b[4:20]
		return nil
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
	ErrorCode         byte
	Error             string
	CommandIdentifier CommandCode
}

func (c *ErrorReport) GetCommandCode() CommandCode {
	return CommandErrorReport
}
func (c *ErrorReport) FromMessage(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("error report length must be exactly 3 bytes, got: %d", len(b))
	}
	c.ErrorCode = b[0]
	c.Error = errorCodeNames[c.ErrorCode]
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
func (c *LockAction) GetPayload() []byte {
	return slices.Concat(
		[]byte{byte(c.Action)},
		c.AppId,
		[]byte{c.Flags},
		c.Nonce,
	)
}

//go:generate stringer -type=NukiState
type NukiState byte

const (
	NukiStateUninitialized   NukiState = 0x00
	NukiStatePairingMode     NukiState = 0x01
	NukiStateDoorMode        NukiState = 0x02
	NukiStateMaintenanceMode NukiState = 0x04
)

type Trigger byte

const (
	TriggerSystem    Trigger = 0x00 // via bluetooth command
	TriggerManual    Trigger = 0x01 // by using a key from outside the door or rotating the wheel on the inside
	TriggerButton    Trigger = 0x02 // by pressing the Smart Lock's button
	TriggerAutomatic Trigger = 0x03 // executed automatically (e.g. at a specific time) by the Smart Lock
	TriggerAutoLock  Trigger = 0x06 // auto lock of the Smart Lock
)

func (t Trigger) String() string {
	switch t {
	case TriggerSystem:
		return "System"
	case TriggerManual:
		return "Manual"
	case TriggerButton:
		return "Button"
	case TriggerAutomatic:
		return "Automatic"
	case TriggerAutoLock:
		return "Auto Lock"
	}
	return ""
}

// KeyturnerStates holds the state information for the Smart Lock.
type DoorSensorState byte

const (
	DoorSensorUnavailable  DoorSensorState = 0x00 // Not paired
	DoorSensorClosed       DoorSensorState = 0x02
	DoorSensorOpened       DoorSensorState = 0x03
	DoorSensorUncalibrated DoorSensorState = 0x10
	DoorSensorTampered     DoorSensorState = 0xF0
	DoorSensorUnknown      DoorSensorState = 0xFF
)

type RemoteAccessStatus struct {
	SSEUplinkAvailable       bool // Bit 0: SSE uplink available via BR/WiFi/Thread
	BridgePaired             bool // Bit 1: Bridge paired
	SSEConnectionViaWiFi     bool // Bit 2: SSE connection via WiFi
	SSEConnectionEstablished bool // Bit 3: SSE connection established
	SSEConnectionViaThread   bool // Bit 4: SSE connection via Thread
	ThreadSSEUplinkEnabled   bool // Bit 5: Thread SSE uplink enabled (manual setting from user)
	NAT64AvailableViaThread  bool // Bit 6: NAT64 available via Thread (potential SSE uplink)
}

func (r RemoteAccessStatus) String() string {
	var vals []string
	if r.SSEUplinkAvailable {
		vals = append(vals, "SSE Uplink Available")
	}
	if r.BridgePaired {
		vals = append(vals, "Bridge Paired")
	}
	if r.SSEConnectionViaWiFi {
		vals = append(vals, "SSE Connection Via WiFi")
	}
	if r.SSEConnectionEstablished {
		vals = append(vals, "SSE Connection Established")
	}
	if r.SSEConnectionViaThread {
		vals = append(vals, "SSE Connection Via Thread")
	}
	if r.ThreadSSEUplinkEnabled {
		vals = append(vals, "Thread SSE Uplink Enabled")
	}
	if r.NAT64AvailableViaThread {
		vals = append(vals, "NAT64 Available Via Thread")
	}
	return strings.Join(vals, ", ")
}

func newRemoteAccessStatus(b byte) RemoteAccessStatus {
	return RemoteAccessStatus{
		SSEUplinkAvailable:       b&0x01 != 0,
		BridgePaired:             b&0x02 != 0,
		SSEConnectionViaWiFi:     b&0x04 != 0,
		SSEConnectionEstablished: b&0x08 != 0,
		SSEConnectionViaThread:   b&0x10 != 0,
		ThreadSSEUplinkEnabled:   b&0x20 != 0,
		NAT64AvailableViaThread:  b&0x40 != 0,
	}
}

type AccessoryStatus struct {
	KeypadSupported           bool // Bit 0: Feature supported by Keypad
	KeypadBatteryCritical     bool // Bit 1: Keypad Battery State Critical
	DoorSensorSupported       bool // Bit 2: Feature supported by Door Sensor
	DoorSensorBatteryCritical bool // Bit 3: Door Sensor Battery State Critical
}

func newAccessoryStatus(b byte) AccessoryStatus {
	return AccessoryStatus{
		KeypadSupported:           b&0x01 != 0,
		KeypadBatteryCritical:     b&0x02 != 0,
		DoorSensorSupported:       b&0x04 != 0,
		DoorSensorBatteryCritical: b&0x08 != 0,
	}
}

type ConnectionStrengthStatus byte

const (
	ConnectionStrengthInvalid      ConnectionStrengthStatus = 0x00
	ConnectionStrengthNotSupported ConnectionStrengthStatus = 0x01
	ConnectionStrengthOK           ConnectionStrengthStatus = 0x02
)

type ConnectionStrength struct {
	RSSI   int8
	Status ConnectionStrengthStatus
}

func newConnectionStrength(b byte) ConnectionStrength {
	switch b {
	case 0x00:
		return ConnectionStrength{RSSI: 0, Status: ConnectionStrengthInvalid}
	case 0x01:
		return ConnectionStrength{RSSI: 0, Status: ConnectionStrengthNotSupported}
	default:
		return ConnectionStrength{RSSI: int8(b), Status: ConnectionStrengthOK}
	}
}

type WifiStatus byte

const (
	WifiDisabled     WifiStatus = 0x00
	WifiDisconnected WifiStatus = 0x01
	WifiConnecting   WifiStatus = 0x02
	WifiConnected    WifiStatus = 0x03
)

type SseStatus byte

const (
	SseSuspended    SseStatus = 0x00
	SseNotReachable SseStatus = 0x01
	SseConnecting   SseStatus = 0x02
	SseConnected    SseStatus = 0x03
)

type WifiConnectionStatus struct {
	WifiStatus  WifiStatus // Bits 0-1: WiFi status
	SseStatus   SseStatus  // Bits 2-3: SSE status
	WifiQuality byte       // Bits 4-7: WiFi quality (0x00 - 0x0F)
}

func newWifiConnectionStatus(b byte) WifiConnectionStatus {
	return WifiConnectionStatus{
		WifiStatus:  WifiStatus(b & 0x03),
		SseStatus:   SseStatus((b >> 2) & 0x03),
		WifiQuality: (b >> 4) & 0x0F,
	}
}

type MqttStatus byte

const (
	MqttDisabled     MqttStatus = 0x00
	MqttDisconnected MqttStatus = 0x01
	MqttConnecting   MqttStatus = 0x02
	MqttConnected    MqttStatus = 0x03
)

type MqttUplink byte

const (
	MqttUplinkWiFi   MqttUplink = 0x00
	MqttUplinkThread MqttUplink = 0x01
)

type MqttConnectionStatus struct {
	MqttStatus MqttStatus
	MqttUplink MqttUplink
}

func newMqttConnectionStatus(b byte) MqttConnectionStatus {
	return MqttConnectionStatus{
		MqttStatus: MqttStatus(b & 0x03),
		MqttUplink: MqttUplink((b >> 2) & 0x01),
	}
}

type ThreadStatus byte

const (
	ThreadStatusMatterDisabled ThreadStatus = 0x00
	ThreadStatusDisconnected   ThreadStatus = 0x01
	ThreadStatusConnecting     ThreadStatus = 0x02
	ThreadStatusConnected      ThreadStatus = 0x03
)

type ThreadConnectionStatus struct {
	ThreadStatus              ThreadStatus
	SseStatus                 SseStatus
	MatterCommissioningActive bool
	WifiSuspended             bool
}

func newThreadConnectionStatus(b byte) ThreadConnectionStatus {
	return ThreadConnectionStatus{
		ThreadStatus:              ThreadStatus(b & 0x03),
		SseStatus:                 SseStatus((b >> 2) & 0x03),
		MatterCommissioningActive: b&0x10 != 0,
		WifiSuspended:             b&0x20 != 0,
	}
}

type KeyturnerStates struct {
	NukiState            NukiState
	LockState            LockState
	Trigger              Trigger
	CurrentTime          time.Time
	TimezoneOffset       int16
	BatteryStateCritical bool
	Charging             bool
	BatteryPercentage    int

	ConfigUpdateCount              byte
	LockNGoTimer                   byte
	LastLockAction                 LockState
	LastLockActionTrigger          Trigger
	LastLockActionCompletionStatus StatusCode
	DoorSensorState                DoorSensorState
	NightmodeActive                bool
	AccessoryBatteryState          AccessoryStatus
	RemoteAccessStatus             RemoteAccessStatus
	BleConnectionStrength          ConnectionStrength
	WifiConnectionStrength         ConnectionStrength
	WifiConnectionStatus           WifiConnectionStatus
	MqttConnectionStatus           MqttConnectionStatus
	ThreadConnectionStatus         ThreadConnectionStatus
}

func (c *KeyturnerStates) GetCommandCode() CommandCode {
	return CommandKeyturnerStates
}
func (c *KeyturnerStates) FromMessage(b []byte) error {
	if len(b) < 26 && len(b) > 27 { // TODO: which one is the correct length?
		return fmt.Errorf("keyturner states length must be between 26 and 27 bytes, got: %d", len(b))
	}
	c.NukiState = NukiState(b[0])
	c.LockState = LockState(b[1])
	c.Trigger = Trigger(b[2])
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
	c.TimezoneOffset = int16(binary.LittleEndian.Uint16(b[10:12]))
	c.BatteryStateCritical = byteToBool(b[12] & 0x01)
	c.Charging = byteToBool(b[12] & 0x02)
	c.BatteryPercentage = int(b[12]>>2) * 2
	c.ConfigUpdateCount = b[13]
	c.LockNGoTimer = b[14]
	c.LastLockAction = LockState(b[15])
	c.LastLockActionTrigger = Trigger(b[16])
	c.LastLockActionCompletionStatus = StatusCode(b[17])
	c.DoorSensorState = DoorSensorState(b[18])
	c.NightmodeActive = byteToBool(b[19])
	c.AccessoryBatteryState = newAccessoryStatus(b[20])
	c.RemoteAccessStatus = newRemoteAccessStatus(b[21])
	c.BleConnectionStrength = newConnectionStrength(b[22])
	c.WifiConnectionStrength = newConnectionStrength(b[23])
	c.WifiConnectionStatus = newWifiConnectionStatus(b[24])
	c.MqttConnectionStatus = newMqttConnectionStatus(b[25])
	c.ThreadConnectionStatus = newThreadConnectionStatus(b[26])
	return nil
}

// Config Command 0x0015
type Config struct {
	NukiID           uint32
	Name             string
	Latitude         float32
	Longitude        float32
	AutoUnlatch      bool
	PairingEnabled   bool
	ButtonEnabled    bool
	LedEnabled       bool
	LedBrightness    uint8
	CurrentTime      time.Time
	TimezoneOffset   int16
	DstMode          uint8
	HasFob           bool
	FobAction1       uint8
	FobAction2       uint8
	FobAction3       uint8
	SingleLock       bool
	AdvertisingMode  uint8
	HasKeypad        bool
	FirmwareVersion  string
	HardwareRevision string
	HomeKitStatus    uint8
	TimezoneID       uint16
	DeviceType       uint8
	Capabilities     uint8
	HasKeypad2       bool
	MatterStatus     uint8
}

func (c *Config) GetTimezoneLocation() *time.Location {
	tz, err := time.LoadLocation(timezoneMap[c.TimezoneID])
	if err != nil {
		return time.UTC
	}
	return tz
}
func (c *Config) FromMessage(b []byte) error {
	if len(b) < 76 {
		return fmt.Errorf("invalid Config message length")
	}

	c.NukiID = binary.LittleEndian.Uint32(b[0:4])
	c.Name = string(bytes.Trim(b[4:36], "\x00"))

	c.Latitude = math.Float32frombits(binary.LittleEndian.Uint32(b[36:40]))
	c.Longitude = math.Float32frombits(binary.LittleEndian.Uint32(b[40:44]))
	c.AutoUnlatch = byteToBool(b[44])
	c.PairingEnabled = byteToBool(b[45])
	c.ButtonEnabled = byteToBool(b[46])
	c.LedEnabled = byteToBool(b[47])
	c.LedBrightness = b[48]

	year := int(binary.LittleEndian.Uint16(b[49:51]))
	month := int(b[51])
	day := int(b[52])
	hour := int(b[53])
	minute := int(b[54])
	second := int(b[55])

	c.TimezoneID = binary.LittleEndian.Uint16(b[72:74])
	tz := c.GetTimezoneLocation()
	c.CurrentTime = time.Date(year, time.Month(month), day, hour, minute, second, 0, tz)

	c.TimezoneOffset = int16(binary.LittleEndian.Uint16(b[56:58]))
	c.DstMode = b[58]
	c.HasFob = byteToBool(b[59])
	c.FobAction1 = b[60]
	c.FobAction2 = b[61]
	c.FobAction3 = b[62]
	c.SingleLock = byteToBool(b[63])
	c.AdvertisingMode = b[64]
	c.HasKeypad = byteToBool(b[65])

	c.FirmwareVersion = fmt.Sprintf("%d.%d.%d", b[66], b[67], b[68])
	c.HardwareRevision = fmt.Sprintf("%d.%d", b[69], b[70])

	c.HomeKitStatus = b[71]
	// timezoneID is set above, next to the current time
	c.DeviceType = b[74]
	c.Capabilities = b[75]
	if len(b) > 76 { // MatterStatus is optional? TODO: verify
		c.HasKeypad2 = byteToBool(b[76])
	}
	if len(b) > 77 { // MatterStatus is optional? TODO: verify
		c.MatterStatus = b[77]
	}
	return nil
}

func (c *Config) GetCommandCode() CommandCode {
	return CommandConfig
}

type RequestConfig struct {
	Nonce []byte
}

func (c *RequestConfig) FromMessage(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("invalid RequestConfig message length")
	}
	c.Nonce = b
	return nil
}

func (c *RequestConfig) GetCommandCode() CommandCode {
	return CommandRequestConfig
}

func (c *RequestConfig) GetPayload() []byte {
	return c.Nonce
}

//go:generate stringer -type=LogSortOrder
type LogSortOrder byte

const (
	LogSortOrderAscending  LogSortOrder = 0x00
	LogSortOrderDescending LogSortOrder = 0x01
)

type RequestLogEntries struct {
	StartIndex  uint32
	Count       uint16
	SortOrder   LogSortOrder
	TotalCount  byte
	Nonce       []byte
	SecurityPin Pin
}

func (c *RequestLogEntries) GetCommandCode() CommandCode {
	return CommandRequestLogEntries
}

func (c *RequestLogEntries) GetPayload() []byte {
	return slices.Concat(
		binary.LittleEndian.AppendUint32(nil, c.StartIndex),
		binary.LittleEndian.AppendUint16(nil, c.Count),
		[]byte{byte(c.SortOrder), c.TotalCount},
		c.Nonce,
		c.SecurityPin.GetPinBytes(),
	)
}

//go:generate stringer -type=LogEntryType
type LogEntryType uint8

const (
	LoggingEnabledDisabled           LogEntryType = 0x01
	LogLockAction                    LogEntryType = 0x02
	LogCalibration                   LogEntryType = 0x03
	LogInitializationRun             LogEntryType = 0x04
	LogKeypadAction                  LogEntryType = 0x05
	LogDoorSensor                    LogEntryType = 0x06
	DoorSensorLoggingEnabledDisabled LogEntryType = 0x07
	LogFirmwareUpdate                LogEntryType = 0x0A
)

type LogEntry struct {
	Index    uint32
	Time     time.Time
	AuthId   uint32
	AuthName string
	Type     LogEntryType
	Data     []byte
}

func (c *LogEntry) FromMessage(b []byte) error {
	if len(b) < 48 {
		return fmt.Errorf("log entry length must be at least 48 bytes, got: %d", len(b))
	}
	c.Index = binary.LittleEndian.Uint32(b[0:4])

	year := int(binary.LittleEndian.Uint16(b[4:6]))
	month := int(b[6])
	day := int(b[7])
	hour := int(b[8])
	minute := int(b[9])
	second := int(b[10])
	// TODO: this should respect the SLs timezone config --> read during authorization
	c.Time = time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)

	c.AuthId = binary.LittleEndian.Uint32(b[11:15])
	c.AuthName = string(bytes.Trim(b[15:47], "\x00"))
	c.Type = LogEntryType(b[47])
	c.Data = b[48:]
	return nil
}

func (c *LogEntry) GetCommandCode() CommandCode {
	return CommandLogEntry
}

func (c *LogEntry) String() string {
	switch c.Type {
	case LoggingEnabledDisabled:
		if c.Data[0] == 0 {
			return "Logging Disabled"
		} else {
			return "Logging Enabled"
		}
	case LogLockAction, LogCalibration, LogInitializationRun:
		var suffix string
		if t := Trigger(c.Data[1]).String(); t != "" {
			suffix = fmt.Sprintf(" (%s)", t)
		}
		return fmt.Sprintf("%v%s", Action(c.Data[0]), suffix)
	}
	if c.Type == LogFirmwareUpdate {
		return fmt.Sprintf("Firmware Update (%d.%d.%d)", c.Data[0], c.Data[1], c.Data[2])
	}
	return fmt.Sprintf("%s by %s (ID: %d) at %s", c.Type.String(), c.AuthName, c.AuthId, c.Time.Format(time.RFC3339))
}

var _ Command = &AuthorizationInfo{}

type AuthorizationInfo struct {
	SecurityPinSet bool
}

func (a *AuthorizationInfo) FromMessage(b []byte) error {
	a.SecurityPinSet = byteToBool(b[0])
	return nil
}

func (a *AuthorizationInfo) GetCommandCode() CommandCode {
	return CommandAuthorizationInfo
}
