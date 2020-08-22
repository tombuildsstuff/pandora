package resourcegroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type Client struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string
}

func NewClient(subscriptionId string, authorizer sdk.Authorizer) Client {
	return NewClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) Client {
	return Client{
		apiVersion:     "2018-05-01",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client Client) Create(ctx context.Context, id ResourceGroupID, input CreateResourceGroupInput) error {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK,      // already exists
			http.StatusCreated, // new
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	if _, err := client.baseClient.PutJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

func (client Client) Delete(ctx context.Context, id ResourceGroupID) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // delete accepted
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

func (client Client) Get(ctx context.Context, id ResourceGroupID) (*GetResourceGroupResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // Exists
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	var out GetResourceGroup
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetResourceGroupResponse{
		HttpResponse:  resp,
		ResourceGroup: &out,
	}
	return &result, nil
}

func (client Client) Update(ctx context.Context, id ResourceGroupID, input UpdateResourceGroupInput) error {
	req := sdk.PatchHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // already exists
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}
	if _, err := client.baseClient.PatchJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

func (client Client) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.Resources"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
