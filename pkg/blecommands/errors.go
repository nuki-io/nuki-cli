package blecommands

var errorCodeNames = map[byte]string{
	0xFD: "ERROR_BAD_CRC",                  // CRC of received command is invalid
	0xFE: "ERROR_BAD_LENGTH",               // Length of retrieved command payload does not match expected length
	0xFF: "ERROR_UNKNOWN",                  // Used if no other error code matches
	0x10: "P_ERROR_NOT_PAIRING",            // Returned if public key is being requested via request data command, but the Smart Lock is not in pairing mode
	0x11: "P_ERROR_BAD_AUTHENTICATOR",      // Returned if the received authenticator does not match the own calculated authenticator
	0x12: "P_ERROR_BAD_PARAMETER",          // Returned if a provided parameter is outside of its valid range
	0x13: "P_ERROR_MAX_USER",               // Returned if the maximum number of users has been reached
	0x20: "K_ERROR_NOT_AUTHORIZED",         // Returned if the provided authorization id is invalid or the payload could not be decrypted using the shared key for this authorization id
	0x21: "K_ERROR_BAD_PIN",                // Returned if the provided pin does not match the stored one.
	0x22: "K_ERROR_BAD_NONCE",              // Returned if the provided nonce does not match the last stored one of this authorization id or has already been used before.
	0x23: "K_ERROR_BAD_PARAMETER",          // Returned if a provided parameter is outside of its valid range.
	0x24: "K_ERROR_INVALID_AUTH_ID",        // Returned if the desired authorization id could not be deleted because it does not exist.
	0x25: "K_ERROR_DISABLED",               // Returned if the provided authorization id is currently disabled.
	0x26: "K_ERROR_REMOTE_NOT_ALLOWED",     // Returned if the request has been forwarded by the Nuki Bridge and the provided authorization id has not been granted remote access.
	0x27: "K_ERROR_TIME_NOT_ALLOWED",       // Returned if the provided authorization id has not been granted access at the current time.
	0x28: "K_ERROR_TOO_MANY_PIN_ATTEMPTS",  // Returned if an invalid pin has been provided too often
	0x29: "K_ERROR_TOO_MANY_ENTRIES",       // Returned if no more entries can be stored
	0x2A: "K_ERROR_CODE_ALREADY_EXISTS",    // Returned if a Keypad Code should be added but the given code already exists.
	0x2B: "K_ERROR_CODE_INVALID",           // Returned if a Keypad Code that has been entered is invalid.
	0x2C: "K_ERROR_CODE_INVALID_TIMEOUT_1", // Returned if an invalid pin has been provided multiple times.
	0x2D: "K_ERROR_CODE_INVALID_TIMEOUT_2", // Returned if an invalid pin has been provided multiple times.
	0x2E: "K_ERROR_CODE_INVALID_TIMEOUT_3", // Returned if an invalid pin has been provided multiple times.
	0x40: "K_ERROR_AUTO_UNLOCK_TOO_RECENT", // Returned on an incoming auto unlock request and if an lock action has already been executed within short time.
	0x41: "K_ERROR_POSITION_UNKNOWN",       // Returned on an incoming unlock request if the request has been forwarded by the Nuki Bridge and the Smart Lock is unsure about its actual lock position.
	0x42: "K_ERROR_MOTOR_BLOCKED",          // Returned if the motor blocks.
	0x43: "K_ERROR_CLUTCH_FAILURE",         // Returned if there is a problem with the clutch during motor movement.
	0x44: "K_ERROR_MOTOR_TIMEOUT",          // Returned if the motor moves for a given period of time but did not block.
	0x45: "K_ERROR_BUSY",                   // Returned on any lock action via bluetooth if there is already a lock action processing.
	0x46: "K_ERROR_CANCELED",               // Returned on any lock action or during calibration if the user canceled the motor movement by pressing the button
	0x47: "K_ERROR_NOT_CALIBRATED",         // Returned on any lock action if the Smart Lock has not yet been calibrated
	0x48: "K_ERROR_MOTOR_POSITION_LIMIT",   // Returned during calibration if the internal position database is not able to store any more values
	0x49: "K_ERROR_MOTOR_LOW_VOLTAGE",      // Returned if the motor blocks because of low voltage.
	0x4A: "K_ERROR_MOTOR_POWER_FAILURE",    // Returned if the power drain during motor movement is zero
	0x4B: "K_ERROR_CLUTCH_POWER_FAILURE",   // Returned if the power drain during clutch movement is zero
	0x4C: "K_ERROR_VOLTAGE_TOO_LOW",        // Returned on a calibration request if the battery voltage is too low and a calibration will therefore not be started
	0x4D: "K_ERROR_FIRMWARE_UPDATE_NEEDED", // Returned during any motor action if a firmware update is mandatory
}
