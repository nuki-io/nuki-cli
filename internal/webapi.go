package internal

import (
	"context"
	"fmt"

	client "github.com/nuki-io/go-nuki"
)

type webApiClient struct {
	cl *client.APIClient
}

type WebApiClient interface {
	GetMyAccount() (*client.MyAccount, error)
	GetDevices() ([]client.Smartlock, error)
}

func NewWebApiClient(apiKey string) WebApiClient {
	cfg := client.NewConfiguration()
	cfg.Host = "api.nuki.io"
	cfg.Scheme = "https"
	cfg.UserAgent = "nukictl"
	cfg.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	return &webApiClient{
		cl: client.NewAPIClient(cfg),
	}
}

func (w *webApiClient) GetMyAccount() (*client.MyAccount, error) {
	accountGet := w.cl.AccountAPI.GetAccounts(context.Background())
	res, _, err := accountGet.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get account details: %w", err)
	}
	return res, nil
}

func (w *webApiClient) GetDevices() ([]client.Smartlock, error) {
	req := w.cl.SmartlockAPI.GetSmartlocks(context.Background())
	res, _, err := req.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get smartlocks: %w", err)
	}
	return res, nil
}
