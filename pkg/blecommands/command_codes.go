package blecommands

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

var responseImplMap = map[CommandCode]func() Response{
	CommandPublicKey:                  func() Response { return &PublicKey{} },
	CommandChallenge:                  func() Response { return &Challenge{} },
	CommandAuthorizationAuthenticator: func() Response { return &AuthorizationAuthenticator{} },
	CommandAuthorizationID:            func() Response { return &AuthorizationID{} },
	// CommandRemoveAuthorizationEntry:    func() Response { return &RemoveAuthorizationEntry{} },
	// CommandRequestAuthorizationEntries: func() Response { return &RequestAuthorizationEntries{} },
	// CommandAuthorizationEntry:          func() Response { return &AuthorizationEntry{} },
	// CommandAuthorizationDataInvite:     func() Response { return &AuthorizationDataInvite{} },
	CommandKeyturnerStates: func() Response { return &KeyturnerStates{} },
	CommandStatus:          func() Response { return &Status{} },
	// CommandMostRecentCommand:           func() Response { return &MostRecentCommand{} },
	// CommandOpeningsClosingsSummary:     func() Response { return &OpeningsClosingsSummary{} },
	// CommandBatteryReport:               func() Response { return &BatteryReport{} },
	CommandErrorReport: func() Response { return &ErrorReport{} },
	// CommandSetConfig:                   func() Response { return &SetConfig{} },
	CommandRequestConfig: func() Response { return &RequestConfig{} },
	CommandConfig:        func() Response { return &Config{} },
	// CommandSetSecurityPIN:              func() Response { return &SetSecurityPIN{} },
	// CommandRequestCalibration:          func() Response { return &RequestCalibration{} },
	// CommandRequestReboot:               func() Response { return &RequestReboot{} },
	// CommandAuthorizationIDInvite:       func() Response { return &AuthorizationIDInvite{} },
	// CommandVerifySecurityPIN:           func() Response { return &VerifySecurityPIN{} },
	// CommandUpdateTime:                  func() Response { return &UpdateTime{} },
	// CommandUpdateAuthorizationEntry:    func() Response { return &UpdateAuthorizationEntry{} },
	// CommandAuthorizationEntryCount:     func() Response { return &AuthorizationEntryCount{} },
	CommandLogEntry: func() Response { return &LogEntry{} },
	// CommandLogEntryCount:             func() Response { return &LogEntryCount{} },
	// CommandEnableLogging:             func() Response { return &EnableLogging{} },
	// CommandSetAdvancedConfig:         func() Response { return &SetAdvancedConfig{} },
	// CommandRequestAdvancedConfig:     func() Response { return &RequestAdvancedConfig{} },
	// CommandAdvancedConfig:            func() Response { return &AdvancedConfig{} },
	// CommandAddTimeControlEntry:       func() Response { return &AddTimeControlEntry{} },
	// CommandTimeControlEntryID:        func() Response { return &TimeControlEntryID{} },
	// CommandRemoveTimeControlEntry:    func() Response { return &RemoveTimeControlEntry{} },
	// CommandRequestTimeControlEntries: func() Response { return &RequestTimeControlEntries{} },
	// CommandTimeControlEntryCount:     func() Response { return &TimeControlEntryCount{} },
	// CommandTimeControlEntry:          func() Response { return &TimeControlEntry{} },
	// CommandUpdateTimeControlEntry:    func() Response { return &UpdateTimeControlEntry{} },
	// CommandAddKeypadCode:             func() Response { return &AddKeypadCode{} },
	// CommandKeypadCodeID:              func() Response { return &KeypadCodeID{} },
	// CommandRequestKeypadCodes:        func() Response { return &RequestKeypadCodes{} },
	// CommandKeypadCodeCount:           func() Response { return &KeypadCodeCount{} },
	// CommandKeypadCode:                func() Response { return &KeypadCode{} },
	// CommandUpdateKeypadCode:          func() Response { return &UpdateKeypadCode{} },
	// CommandRemoveKeypadCode:          func() Response { return &RemoveKeypadCode{} },
	CommandAuthorizationInfo: func() Response { return &AuthorizationInfo{} },
	// CommandSimpleLockAction:          func() Response { return &SimpleLockAction{} },
}
