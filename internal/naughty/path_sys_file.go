package naughty

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathSysFile(b *naughtyBackend) *framework.Path {
	return &framework.Path{
		Pattern: "file/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the file to steal on the host",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleFileRead,
			},
		},
	}
}

func (b *naughtyBackend) handleFileRead(_ context.Context, _ *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	return nil, nil
}
