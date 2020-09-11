package eventhub

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type EventhubClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewEventhubClient(subscriptionId string, authorizer sdk.Authorizer) EventhubClient {
	return NewEventhubClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewEventhubClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) EventhubClient {
	return EventhubClient{
		apiVersion:     "2018-01-01-preview",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client EventhubClient) Create(ctx context.Context, id EventhubId, input CreateEventHubInput) error {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusCreated, // TODO: unknown
			http.StatusOK,      // TODO: unknown
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	if _, err := client.baseClient.PutJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

type GetEventhubResponse struct {
	HttpResponse *http.Response
	Eventhub     *GetEventHub
}

func (client EventhubClient) Get(ctx context.Context, id EventhubId) (*GetEventhubResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out GetEventHub
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetEventhubResponse{
		HttpResponse: resp,
		Eventhub:     &out,
	}
	return &result, nil
}

func (client EventhubClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.EventHub"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
