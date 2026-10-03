package authstore

import (
	"encoding/hex"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/viper"
)

// ViperStore implements bleflows.AuthStore using viper.
// Store only updates v in memory; writing the config file is up to the caller.
type ViperStore struct {
	v *viper.Viper
}

func New(v *viper.Viper) *ViperStore {
	return &ViperStore{v: v}
}

type authorizeContextStorage struct {
	CliPublicKey  string
	CliPrivateKey string
	SlPublicKey   string
	SharedKey     string
	AuthId        string
	AppId         string
	NukiId        string
	Pin           string
	Name          string
}

func contextToStorage(ac *bleflows.AuthorizeContext) *authorizeContextStorage {
	return &authorizeContextStorage{
		CliPublicKey:  fmt.Sprintf("%x", ac.CliPublicKey),
		CliPrivateKey: fmt.Sprintf("%x", ac.CliPrivateKey),
		SlPublicKey:   fmt.Sprintf("%x", ac.SlPublicKey),
		SharedKey:     fmt.Sprintf("%x", ac.SharedKey),
		AuthId:        fmt.Sprintf("%x", ac.AuthId),
		AppId:         fmt.Sprintf("%x", ac.AppId),
		NukiId:        fmt.Sprintf("%X", ac.NukiId),
		Pin:           ac.Pin,
		Name:          ac.Name,
	}
}

func storageToContext(s *authorizeContextStorage) *bleflows.AuthorizeContext {
	ac := &bleflows.AuthorizeContext{}
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
	ac.Pin = s.Pin
	ac.Name = s.Name
	return ac
}

func (s *ViperStore) Load(deviceId string) (*bleflows.AuthorizeContext, error) {
	cfgKey := fmt.Sprintf("authorizations.%s", deviceId)
	if !s.v.IsSet(cfgKey) {
		return nil, fmt.Errorf("no authorization for device with id %s found", deviceId)
	}
	st := &authorizeContextStorage{}
	s.v.UnmarshalKey(cfgKey, st)
	return storageToContext(st), nil
}

func (s *ViperStore) Store(deviceId string, ctx *bleflows.AuthorizeContext) error {
	cfgKey := fmt.Sprintf("authorizations.%s", deviceId)
	s.v.Set(cfgKey, contextToStorage(ctx))
	return nil
}
