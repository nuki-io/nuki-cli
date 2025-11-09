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
