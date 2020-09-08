package namespaces

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type EventHubNamespaceClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewEventHubNamespaceClient(subscriptionId string, authorizer sdk.Authorizer) EventHubNamespaceClient {
	return NewEventHubNamespaceClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewEventHubNamespaceClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) EventHubNamespaceClient {
	return EventHubNamespaceClient{
		apiVersion:     "2018-01-01-preview",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client EventHubNamespaceClient) Create(ctx context.Context, id EventHubNamespaceId, input CreateNamespaceInput) (sdk.Poller, error) {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // TODO: unknown,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.PutJsonThenPoll(ctx, req)
}

func (client EventHubNamespaceClient) Delete(ctx context.Context, id EventHubNamespaceId) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // deletion accepted,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

type GetEventHubNamespaceResponse struct {
	HttpResponse      *http.Response
	EventHubNamespace *GetNamespace
}

func (client EventHubNamespaceClient) Get(ctx context.Context, id EventHubNamespaceId) (*GetEventHubNamespaceResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	var out GetNamespace
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetEventHubNamespaceResponse{
		HttpResponse:      resp,
		EventHubNamespace: &out,
	}
	return &result, nil
}

func (client EventHubNamespaceClient) Update(ctx context.Context, id EventHubNamespaceId, input PatchNamespaceInput) error {
	req := sdk.PatchHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // TODO: unknown
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	if _, err := client.baseClient.PatchJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

func (client EventHubNamespaceClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.EventHub"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
