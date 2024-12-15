package twiliosecrets

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/google/uuid"
	"github.com/hashicorp/cloud-vault-secrets/credentialprovider/builtin"
	"github.com/hashicorp/cloud-vault-secrets/credentialprovider/builtin/twiliointegration"
	"github.com/hashicorp/cloud-vault-secrets/credentialprovider/proto/go/credentials"
	credspb "github.com/hashicorp/cloud-vault-secrets/credentialprovider/proto/go/credentials"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAPIKey(b *twilioBackend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Friendly name of the api key",
				Required:    true,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			// TODO add logical write to configure the credential
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleAPIKeyRead,
			},
		},
	}
}

func (b *twilioBackend) handleAPIKeyRead(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	integration, pbProviderConfig, pbCredentialConfig, err := b.initIntegration(ctx, req)
	if err != nil {
		return nil, err
	}

	operationId := uuid.New().String()
	stageResp, err := integration.Stage(ctx, &credentials.StageRequest{
		OperationId:      operationId, // TODO Cadence leaking out of the interface, we should use a different mechanism
		ProviderConfig:   pbProviderConfig,
		CredentialConfig: pbCredentialConfig,
		CredentialType:   twiliointegration.CredentialTypeTwilioAPIKey,
	})
	if err != nil {
		return nil, err
	}

	_, err = integration.Commit(ctx, &credentials.CommitRequest{
		OperationId:      operationId, // TODO Cadence leaking out of the interface, we should use a different mechanism
		ProviderConfig:   pbProviderConfig,
		CredentialConfig: pbCredentialConfig,
		CredentialSet:    stageResp.CredentialSet,
		CredentialType:   twiliointegration.CredentialTypeTwilioAPIKey,
	})
	if err != nil {
		return nil, err
	}

	marshaledValueData, err := json.Marshal(stageResp.CredentialSet.Values)
	if err != nil {
		return nil, err
	}
	var valueData map[string]any
	err = json.Unmarshal(marshaledValueData, &valueData)
	if err != nil {
		return nil, err
	}

	marshaledInternalData, err := json.Marshal(stageResp.CredentialSet.InternalData)
	if err != nil {
		return nil, err
	}
	var internalData map[string]any
	err = json.Unmarshal(marshaledInternalData, &internalData)
	if err != nil {
		return nil, err
	}

	// TODO we need the SID to revoke the API key when the lease expires
	// The Twilio SDK implementation should persist it as internal data directly
	apiKeySid, ok := valueData[twiliointegration.APIKeySidKey]
	if !ok {
		return nil, fmt.Errorf("couldn't save Twilio API key SId to internal data")
	}
	internalData[twiliointegration.APIKeySidKey] = apiKeySid

	resp := b.Secret(twiliointegration.CredentialTypeTwilioAPIKey).Response(
		valueData,
		internalData,
	)

	// TODO should be configurable or inferred from upstream
	// Low TTL while exploring to test the revoke operation easily
	resp.Secret.TTL = 5 * time.Second
	resp.Secret.MaxTTL = 10 * time.Second

	return resp, nil
}

func (b *twilioBackend) handleAPIKeyRevoke(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	integration, pbProviderConfig, _, err := b.initIntegration(ctx, req)
	if err != nil {
		return nil, err
	}

	apiKeySid, ok := req.Secret.InternalData[twiliointegration.APIKeySidKey]
	if !ok {
		return nil, fmt.Errorf("couldn't revoke Twilio API key, missing internal data for API key SID")
	}

	pbInternalData, err := builtin.MarshalStruct(&twiliointegration.APIKeyCredentialConfig{
		ApiKeySid: apiKeySid.(string),
	})
	if err != nil {
		return nil, err
	}

	operationId := uuid.New().String()
	_, err = integration.Revoke(ctx, &credentials.RevokeRequest{
		OperationId:    operationId, // Cadence leaking out of the interface
		ProviderConfig: pbProviderConfig,
		CredentialSet: &credentials.CredentialSet{
			InternalData: pbInternalData,
		},
		CredentialType: twiliointegration.CredentialTypeTwilioAPIKey,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *twilioBackend) initIntegration(ctx context.Context, req *logical.Request) (
	credspb.CredentialLifecycleServiceServer,
	*structpb.Struct,
	*structpb.Struct,
	error,
) {
	providerConfig, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, nil, nil, err
	}
	pbProviderConfig, err := builtin.MarshalStruct(providerConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	credentialConfig := &twiliointegration.APIKeyCredentialConfig{}
	pbCredentialConfig, err := builtin.MarshalStruct(credentialConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	integration, err := twiliointegration.NewIntegration()
	if err != nil {
		return nil, nil, nil, err
	}

	return integration, pbProviderConfig, pbCredentialConfig, nil
}
