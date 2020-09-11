package configurationStore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type ConfigurationStoreClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewConfigurationStoreClient(subscriptionId string, authorizer sdk.Authorizer) ConfigurationStoreClient {
	return NewConfigurationStoreClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewConfigurationStoreClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) ConfigurationStoreClient {
	return ConfigurationStoreClient{
		apiVersion:     "2019-10-01",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client ConfigurationStoreClient) Create(ctx context.Context, id ConfigurationStoreId, input CreateStoreInput) (sdk.Poller, error) {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // TODO: unknown
			http.StatusCreated,  // TODO: unknown,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.PutJsonThenPoll(ctx, req)
}

func (client ConfigurationStoreClient) Delete(ctx context.Context, id ConfigurationStoreId) (*http.Response, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // deleted
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.Delete(ctx, req)
}

type GetConfigurationStoreResponse struct {
	HttpResponse       *http.Response
	ConfigurationStore *GetStore
}

func (client ConfigurationStoreClient) Get(ctx context.Context, id ConfigurationStoreId) (*GetConfigurationStoreResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out GetStore
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetConfigurationStoreResponse{
		HttpResponse:       resp,
		ConfigurationStore: &out,
	}
	return &result, nil
}

func (client ConfigurationStoreClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.AppConfiguration"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
