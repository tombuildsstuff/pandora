package sdk

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type Authorizer interface {
	Token(ctx context.Context, endpoint string) (*Token, error)
}

type Token struct {
	// todo: expires etc
	accessToken string
	kind        string
}

func (t Token) AuthorizationHeader() string {
	return fmt.Sprintf("%s %s", t.kind, t.accessToken)
}

type ClientSecretAuthorizer struct {
	activeDirectoryEndpoint string
	clientId                string
	clientSecret            string
	tenantId                string
}

func NewClientSecretAuthorizer(clientId, clientSecret, tenantId string) Authorizer {
	return NewClientSecretAuthorizerForEndpoint(clientId, clientSecret, tenantId, endpoints.DefaultActiveDirectoryEndpoint)
}

func NewClientSecretAuthorizerForEndpoint(clientId, clientSecret, tenantId, activeDirectoryEndpoint string) Authorizer {
	return &ClientSecretAuthorizer{
		activeDirectoryEndpoint: activeDirectoryEndpoint,
		clientId:     clientId,
		clientSecret: clientSecret,
		tenantId:     tenantId,
	}
}

func (a ClientSecretAuthorizer) Token(ctx context.Context, endpoint string) (*Token, error) {
	// TODO: obviously make something ourselves here
	oauth, err := adal.NewOAuthConfig(a.activeDirectoryEndpoint, a.tenantId)
	if err != nil {
		return nil, err
	}

	spt, err := adal.NewServicePrincipalToken(*oauth, a.clientId, a.clientSecret, endpoint)
	if err != nil {
		return nil, err
	}

	if err = spt.RefreshWithContext(ctx); err != nil {
		return nil, err
	}

	token := spt.Token()
	result := Token{
		accessToken: token.AccessToken,
		kind:        "Bearer",
	}
	return &result, nil
}
