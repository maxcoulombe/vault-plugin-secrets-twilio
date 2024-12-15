package twiliosecrets

import (
	"context"

	"github.com/hashicorp/cloud-vault-secrets/credentialprovider/builtin/twiliointegration"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const twiliointegrationDotAccountID = "account_sid" // TOOO export this constant from twiliointegration
const configStoragePath = "config"

func pathConfig(b *twilioBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config/root",
		Fields: map[string]*framework.FieldSchema{
			twiliointegrationDotAccountID: {
				Type:        framework.TypeString,
				Description: "Twilio Account SID",
				Required:    true,
			},
			twiliointegration.APIKeySidKey: {
				Type:        framework.TypeString,
				Description: "Twilio API Key SID",
				Required:    true,
			},
			twiliointegration.APIKeySecretKey: {
				Type:        framework.TypeString,
				Description: "Twilio API Key Secret",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.handleConfigWrite,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.handleConfigWrite,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleConfigRead,
			},
		},
	}
}

func (b *twilioBackend) handleConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &twiliointegration.ProviderConfig{}
	}

	if accountSID, ok := data.GetOk(twiliointegrationDotAccountID); ok {
		config.AccountSid = accountSID.(string)
	}
	if apiKeySID, ok := data.GetOk(twiliointegration.APIKeySidKey); ok {
		config.ApiKeySid = apiKeySID.(string)
	}
	if apiKeySecret, ok := data.GetOk(twiliointegration.APIKeySecretKey); ok {
		config.ApiKeySecret = apiKeySecret.(string)
	}

	// TODO Call SDK to validate the root config

	err = b.saveConfig(ctx, config, req.Storage)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (b *twilioBackend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &twiliointegration.ProviderConfig{}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			twiliointegrationDotAccountID:     config.AccountSid,
			twiliointegration.APIKeySidKey:    config.ApiKeySid,
			twiliointegration.APIKeySecretKey: "***", // Obfuscate sensitive root config data
		},
	}

	return resp, nil
}

func (b *twilioBackend) config(ctx context.Context, s logical.Storage) (*twiliointegration.ProviderConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	config := &twiliointegration.ProviderConfig{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (b *twilioBackend) saveConfig(ctx context.Context, config *twiliointegration.ProviderConfig, s logical.Storage) error {
	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return err
	}

	err = s.Put(ctx, entry)
	if err != nil {
		return err
	}

	return nil
}
