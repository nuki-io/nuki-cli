package bleflows

import (
	"crypto/hmac"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/spf13/viper"
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
	Name          string
}

func (ac *AuthorizeContext) toStorage() *authorizeContextStorage {
	return &authorizeContextStorage{
		CliPublicKey:  fmt.Sprintf("%x", ac.CliPublicKey),
		CliPrivateKey: fmt.Sprintf("%x", ac.CliPrivateKey),
		SlPublicKey:   fmt.Sprintf("%x", ac.SlPublicKey),
		SharedKey:     fmt.Sprintf("%x", ac.SharedKey),
		AuthId:        fmt.Sprintf("%x", ac.AuthId),
		AppId:         fmt.Sprintf("%x", ac.AppId),
		NukiId:        fmt.Sprintf("%X", ac.NukiId),
		Name:          ac.Name,
	}
}
func (ac *AuthorizeContext) fromStorage(s *authorizeContextStorage) {
	if v, err := hex.DecodeString(s.CliPublicKey); err == nil {
		ac.CliPublicKey = v
	}
	if v, err := hex.DecodeString(s.CliPrivateKey); err == nil {
		ac.CliPrivateKey = v
	}
	if v, err := hex.DecodeString(s.SlPublicKey); err == nil {
		ac.SlPublicKey = v
	}
	if v, err := hex.DecodeString(s.SharedKey); err == nil {
		ac.SharedKey = v
	}
	if v, err := hex.DecodeString(s.AuthId); err == nil {
		ac.AuthId = v
	}
	if v, err := hex.DecodeString(s.AppId); err == nil {
		ac.AppId = v
	}
	ac.Name = s.Name
}

type authorizeContextStorage struct {
	CliPublicKey  string
	CliPrivateKey string
	SlPublicKey   string
	SharedKey     string
	AuthId        string
	AppId         string
	NukiId        string
	Name          string
}

func NewAuthorizeContext() *AuthorizeContext {
	ctx := &AuthorizeContext{}
	ctx.AppId = binary.LittleEndian.AppendUint32(ctx.AppId, rand.Uint32())
	return ctx
}

func (ac *AuthorizeContext) Load(id string) error {
	cfgKey := fmt.Sprintf("authorizations.%s", id)
	if !viper.IsSet(cfgKey) {
		return fmt.Errorf("no authorization for device with id %s found", id)
	}
	s := &authorizeContextStorage{}
	viper.UnmarshalKey(cfgKey, s)
	ac.fromStorage(s)
	return nil
}

func (ac *AuthorizeContext) Store(id string) {
	cfgKey := fmt.Sprintf("authorizations.%s", id)
	viper.Set(cfgKey, ac.toStorage())
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
	pub, priv, err := box.GenerateKey(crypto_rand.Reader)
	if err != nil {
		panic(err)
	}
	ac.CliPublicKey = pub[:]
	ac.CliPrivateKey = priv[:]
}
