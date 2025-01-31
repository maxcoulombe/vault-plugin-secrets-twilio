package naughty

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathSysDDoS(b *naughtyBackend) *framework.Path {
	return &framework.Path{
		Pattern: "ddos",
		Fields:  map[string]*framework.FieldSchema{},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleFileRead,
			},
		},
	}
}

func (b *naughtyBackend) handleSysDDoS(_ context.Context, _ *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
