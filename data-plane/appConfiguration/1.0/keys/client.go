package keys

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
)

type KeysClient struct {
	apiVersion string
	baseClient sdk.BaseClient
}

func NewKeysClientWithBaseURI(endpoint string, authorizer sdk.Authorizer) KeysClient {
	return KeysClient{
		apiVersion: "1.0",
		baseClient: sdk.DefaultBaseClient(endpoint, authorizer),
	}
}

type GetKeysResponse struct {
	HttpResponse *http.Response
	Keys         *GetKeys
}

func (client KeysClient) Get(ctx context.Context, id KeysId) (*GetKeysResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildDataPlaneURI(id, client.baseClient.Endpoint, client.apiVersion),
	}

	var out GetKeys
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
