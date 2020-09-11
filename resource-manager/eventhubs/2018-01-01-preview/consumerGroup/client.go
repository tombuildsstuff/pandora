package consumerGroup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type ConsumerGroupClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewConsumerGroupClient(subscriptionId string, authorizer sdk.Authorizer) ConsumerGroupClient {
	return NewConsumerGroupClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewConsumerGroupClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) ConsumerGroupClient {
	return ConsumerGroupClient{
		apiVersion:     "2018-01-01-preview",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client ConsumerGroupClient) Create(ctx context.Context, id ConsumerGroupId, input CreateConsumerGroupInput) error {
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

func (client ConsumerGroupClient) Delete(ctx context.Context, id ConsumerGroupId) (*http.Response, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // deleted
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.Delete(ctx, req)
}

type GetConsumerGroupResponse struct {
	HttpResponse  *http.Response
	ConsumerGroup *GetConsumerGroup
}

func (client ConsumerGroupClient) Get(ctx context.Context, id ConsumerGroupId) (*GetConsumerGroupResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out GetConsumerGroup
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetConsumerGroupResponse{
		HttpResponse:  resp,
		ConsumerGroup: &out,
	}
	return &result, nil
}

func (client ConsumerGroupClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.EventHub"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
