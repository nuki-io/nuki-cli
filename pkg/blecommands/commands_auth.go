package blecommands

import (
	"fmt"
	"slices"
)

var _ Request = &PublicKey{}

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

var _ Request = &AuthorizationAuthenticator{}

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

var _ Request = &AuthorizationData{}

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

var _ Request = &AuthorizationData5G{}

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
	var pinBytes []byte
	if a.SecurityPin != nil {
		pinBytes = a.SecurityPin.GetPinBytes()
	}
	return slices.Concat(
		a.Id,
		appName[:],
		pinBytes,
	)
}

var _ Request = &AuthorizationIDConfirmation{}

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

var _ Response = &AuthorizationID{}

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

var _ Response = &AuthorizationInfo{}

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
