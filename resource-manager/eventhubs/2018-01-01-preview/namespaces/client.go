package namespaces

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
	subscriptionId string // TODO: making this Optional?
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

func (client NamespacesClient) Create(ctx context.Context, id NamespacesId, input CreateNamespaceInput) (sdk.Poller, error) {
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

func (client NamespacesClient) Delete(ctx context.Context, id NamespacesId) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // deletion accepted
			http.StatusOK,       // deletion started,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

type GetNamespacesResponse struct {
	HttpResponse *http.Response
	Namespaces   *GetNamespace
}

func (client NamespacesClient) Get(ctx context.Context, id NamespacesId) (*GetNamespacesResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out GetNamespace
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetNamespacesResponse{
		HttpResponse: resp,
		Namespaces:   &out,
	}
	return &result, nil
}

func (client NamespacesClient) Update(ctx context.Context, id NamespacesId, input UpdateNamespaceInput) error {
	req := sdk.PatchHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // TODO: unknown
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	if _, err := client.baseClient.PatchJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

func (client NamespacesClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.EventHub"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
