package bleflows

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"slices"

	"golang.org/x/crypto/nacl/box"
)

type AuthorizeContext struct {
	CliPublicKey  []byte
	CliPrivateKey []byte
	SlPublicKey   []byte
	SharedKey     []byte
	AuthId        []byte
	AppId         []byte
	NukiId        uint32
	Pin           string
	Name          string
}

func NewAuthorizeContext() *AuthorizeContext {
	ctx := &AuthorizeContext{}
	ctx.AppId = make([]byte, 4)
	if _, err := rand.Read(ctx.AppId); err != nil {
		panic(err)
	}
	return ctx
}

func (ac *AuthorizeContext) CalculateSharedKey() {
	sharedKey := [32]byte{}
	box.Precompute(&sharedKey, (*[32]byte)(ac.SlPublicKey), (*[32]byte)(ac.CliPrivateKey))
	ac.SharedKey = sharedKey[:]
}

func (ac *AuthorizeContext) GetMessageAuthenticator(parts ...[]byte) []byte {
	h := hmac.New(sha256.New, ac.SharedKey)
	h.Write(slices.Concat(parts...))
	return h.Sum(nil)
}

func (ac *AuthorizeContext) GenerateKeyPair() {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	ac.CliPublicKey = pub[:]
	ac.CliPrivateKey = priv[:]
}
