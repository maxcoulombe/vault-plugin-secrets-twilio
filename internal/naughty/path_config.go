package naughty

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const configStoragePath = "config"

type Config struct{}

func pathConfig(b *naughtyBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields:  map[string]*framework.FieldSchema{},
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

func (b *naughtyBackend) handleConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &Config{}
	}

	err = b.saveConfig(ctx, config, req.Storage)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (b *naughtyBackend) handleConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &Config{}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}

	return resp, nil
}

func (b *naughtyBackend) config(ctx context.Context, s logical.Storage) (*Config, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	config := &Config{}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}

	return config, nil
}

func (b *naughtyBackend) saveConfig(ctx context.Context, config *Config, s logical.Storage) error {
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
