package twiliosecrets

import (
	"context"
	"strings"

	"github.com/hashicorp/cloud-vault-secrets/credentialprovider/builtin/twiliointegration"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type twilioBackend struct {
	*framework.Backend

	logger hclog.Logger
}

func TwilioBackendFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &twilioBackend{}
	return b.Setup(ctx, conf)
}

func (b *twilioBackend) Setup(_ context.Context, conf *logical.BackendConfig) (logical.Backend, error) {

	b.logger = conf.Logger
	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help: strings.TrimSpace(`
            The Twilio secrets engine dynamically manages Twilio API keys.
        `),
		Paths: []*framework.Path{
			pathConfig(b),
			pathAPIKey(b),
		},
		Secrets: []*framework.Secret{
			{
				Type: twiliointegration.CredentialTypeTwilioAPIKey,
				Fields: map[string]*framework.FieldSchema{
					twiliointegration.APIKeySidKey: {
						Type:        framework.TypeString,
						Description: "API Key Public ID",
					},
					twiliointegration.APIKeySecretKey: {
						Type:        framework.TypeString,
						Description: "API Key Secret",
					},
				},
				Renew:  nil, // TODO the Credentials Lifecycle Interface doesn't support the renewal operation
				Revoke: b.handleAPIKeyRevoke,
			},
		},
	}
	return b.Backend, nil
}
