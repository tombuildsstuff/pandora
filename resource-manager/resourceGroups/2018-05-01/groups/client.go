package groups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type ResourceGroupClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewResourceGroupClient(subscriptionId string, authorizer sdk.Authorizer) ResourceGroupClient {
	return NewResourceGroupClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewResourceGroupClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) ResourceGroupClient {
	return ResourceGroupClient{
		apiVersion:     "2018-05-01",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client ResourceGroupClient) Create(ctx context.Context, id ResourceGroupId, input CreateResourceGroupInput) error {
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // TODO: unknown
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	if _, err := client.baseClient.PutJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

func (client ResourceGroupClient) Delete(ctx context.Context, id ResourceGroupId) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // deletion accepted,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

type GetResourceGroupResponse struct {
	HttpResponse     *http.Response
	GetResourceGroup *GetResourceGroup
}

func (client ResourceGroupClient) Get(ctx context.Context, id ResourceGroupId) (*GetResourceGroupResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	var out GetResourceGroup
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetResourceGroupResponse{
		HttpResponse:     resp,
		GetResourceGroup: &out,
	}
	return &result, nil
}

func (client ResourceGroupClient) Update(ctx context.Context, id ResourceGroupId, input UpdateResourceGroupInput) error {
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

func (client ResourceGroupClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.Resources"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
