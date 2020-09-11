package keys

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type KeysClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewKeysClient(subscriptionId string, authorizer sdk.Authorizer) KeysClient {
	return NewKeysClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewKeysClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) KeysClient {
	return KeysClient{
		apiVersion:     "1.0",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

type GetKeysResponse struct {
	HttpResponse *http.Response
	Keys         *GetKeysResponse
}

func (client KeysClient) Get(ctx context.Context, id KeysId) (*GetKeysResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out GetKeysResponse
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetKeysResponse{
		HttpResponse: resp,
		Keys:         &out,
	}
	return &result, nil
}

func (client KeysClient) MetaData() sdk.ClientMetaData {
	return sdk.ClientMetaData{}
}
