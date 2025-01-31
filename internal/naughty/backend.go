package naughty

import (
	"context"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type naughtyBackend struct {
	*framework.Backend

	logger hclog.Logger
}

func BackendFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &naughtyBackend{}
	return b.Setup(ctx, conf)
}

func (b *naughtyBackend) Setup(_ context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	b.logger = conf.Logger
	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help: strings.TrimSpace(`
            The Naughty secrets engine is a toolbox to validate potential ways a malicious plugin can exploit the Vault server.
        `),
		Paths: []*framework.Path{
			pathConfig(b),
			pathEnvVar(b),
			pathSysFile(b),
			pathSysDDoS(b),
		},
	}
	return b.Backend, nil
}
