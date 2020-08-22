package eventhub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type NamespacesClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string
}

func NewNamespacesClient(subscriptionId string, authorizer sdk.Authorizer) NamespacesClient {
	return NewNamespacesClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewNamespacesClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) NamespacesClient {
	return NamespacesClient{
		apiVersion:     "2018-01-01-preview",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client NamespacesClient) Create(ctx context.Context, id NamespaceID, input CreateNamespaceInput) (sdk.Poller, error) {
	uri := sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion)
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // new
		},
		Uri: uri,
	}

	return client.baseClient.PutJsonThenPoll(ctx, req)
}

func (client NamespacesClient) Delete(ctx context.Context, id NamespaceID) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // delete accepted
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

func (client NamespacesClient) Get(ctx context.Context, id NamespaceID) (*GetNamespaceResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // Exists
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	var out GetNamespace
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetNamespaceResponse{
		HttpResponse: resp,
		Namespace:    &out,
	}
	return &result, nil
}

func (client NamespacesClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.EventHub"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
