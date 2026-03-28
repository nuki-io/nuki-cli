package blecommands

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strings"
	"time"
)

//go:generate stringer -type=Action
type Action uint8

func (a Action) MarshalText() ([]byte, error) { return []byte(a.String()), nil }

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

//go:generate stringer -type=LockState -trimprefix=LockState
type LockState uint8

func (s LockState) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

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

var _ Request = &LockAction{}

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

//go:generate stringer -type=NukiState -trimprefix=NukiState
type NukiState byte

func (s NukiState) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

const (
	NukiStateUninitialized   NukiState = 0x00
	NukiStatePairingMode     NukiState = 0x01
	NukiStateDoorMode        NukiState = 0x02
	NukiStateMaintenanceMode NukiState = 0x04
)

type Trigger byte

func (t Trigger) MarshalText() ([]byte, error) { return []byte(t.String()), nil }

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
	SSEUplinkAvailable       bool `json:"sseUplinkAvailable"`       // Bit 0: SSE uplink available via BR/WiFi/Thread
	BridgePaired             bool `json:"bridgePaired"`             // Bit 1: Bridge paired
	SSEConnectionViaWiFi     bool `json:"sseConnectionViaWifi"`     // Bit 2: SSE connection via WiFi
	SSEConnectionEstablished bool `json:"sseConnectionEstablished"` // Bit 3: SSE connection established
	SSEConnectionViaThread   bool `json:"sseConnectionViaThread"`   // Bit 4: SSE connection via Thread
	ThreadSSEUplinkEnabled   bool `json:"threadSseUplinkEnabled"`   // Bit 5: Thread SSE uplink enabled (manual setting from user)
	NAT64AvailableViaThread  bool `json:"nat64AvailableViaThread"`  // Bit 6: NAT64 available via Thread (potential SSE uplink)
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
	KeypadSupported           bool `json:"keypadSupported"`           // Bit 0: Feature supported by Keypad
	KeypadBatteryCritical     bool `json:"keypadBatteryCritical"`     // Bit 1: Keypad Battery State Critical
	DoorSensorSupported       bool `json:"doorSensorSupported"`       // Bit 2: Feature supported by Door Sensor
	DoorSensorBatteryCritical bool `json:"doorSensorBatteryCritical"` // Bit 3: Door Sensor Battery State Critical
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
	RSSI   int8                     `json:"rssi"`
	Status ConnectionStrengthStatus `json:"status"`
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
	WifiStatus  WifiStatus `json:"wifiStatus"`  // Bits 0-1: WiFi status
	SseStatus   SseStatus  `json:"sseStatus"`   // Bits 2-3: SSE status
	WifiQuality byte       `json:"wifiQuality"` // Bits 4-7: WiFi quality (0x00 - 0x0F)
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
	MqttStatus MqttStatus `json:"mqttStatus"`
	MqttUplink MqttUplink `json:"mqttUplink"`
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
	ThreadStatus              ThreadStatus `json:"threadStatus"`
	SseStatus                 SseStatus    `json:"sseStatus"`
	MatterCommissioningActive bool         `json:"matterCommissioningActive"`
	WifiSuspended             bool         `json:"wifiSuspended"`
}

func newThreadConnectionStatus(b byte) ThreadConnectionStatus {
	return ThreadConnectionStatus{
		ThreadStatus:              ThreadStatus(b & 0x03),
		SseStatus:                 SseStatus((b >> 2) & 0x03),
		MatterCommissioningActive: b&0x10 != 0,
		WifiSuspended:             b&0x20 != 0,
	}
}

var _ Response = &KeyturnerStates{}

type KeyturnerStates struct {
	NukiState            NukiState `json:"nukiState"`
	LockState            LockState `json:"lockState"`
	Trigger              Trigger   `json:"trigger"`
	CurrentTime          time.Time `json:"currentTime"`
	TimezoneOffset       int16     `json:"timezoneOffset"`
	BatteryStateCritical bool      `json:"batteryStateCritical"`
	Charging             bool      `json:"charging"`
	BatteryPercentage    int       `json:"batteryPercentage"`

	ConfigUpdateCount              byte                   `json:"configUpdateCount"`
	LockNGoTimer                   byte                   `json:"lockNGoTimer"`
	LastLockAction                 LockState              `json:"lastLockAction"`
	LastLockActionTrigger          Trigger                `json:"lastLockActionTrigger"`
	LastLockActionCompletionStatus StatusCode             `json:"lastLockActionCompletionStatus"`
	DoorSensorState                DoorSensorState        `json:"doorSensorState"`
	NightmodeActive                bool                   `json:"nightmodeActive"`
	AccessoryBatteryState          AccessoryStatus        `json:"accessoryBatteryState"`
	RemoteAccessStatus             RemoteAccessStatus     `json:"remoteAccessStatus"`
	BleConnectionStrength          ConnectionStrength     `json:"bleConnectionStrength"`
	WifiConnectionStrength         ConnectionStrength     `json:"wifiConnectionStrength"`
	WifiConnectionStatus           WifiConnectionStatus   `json:"wifiConnectionStatus"`
	MqttConnectionStatus           MqttConnectionStatus   `json:"mqttConnectionStatus"`
	ThreadConnectionStatus         ThreadConnectionStatus `json:"threadConnectionStatus"`
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
	c.CurrentTime = fromNukiTime(b[3:10], time.UTC)
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
