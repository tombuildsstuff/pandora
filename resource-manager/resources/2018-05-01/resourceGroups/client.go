package resourceGroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type ResourceGroupsClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

func NewResourceGroupsClient(subscriptionId string, authorizer sdk.Authorizer) ResourceGroupsClient {
	return NewResourceGroupsClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func NewResourceGroupsClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) ResourceGroupsClient {
	return ResourceGroupsClient{
		apiVersion:     "2018-05-01",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}

func (client ResourceGroupsClient) Create(ctx context.Context, id ResourceGroupsId, input CreateInput) error {
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

func (client ResourceGroupsClient) Delete(ctx context.Context, id ResourceGroupsId) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // deletion accepted
			http.StatusOK,       // deletion started,
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

type GetResourceGroupsResponse struct {
	HttpResponse   *http.Response
	ResourceGroups *GetResourceGroup
}

func (client ResourceGroupsClient) Get(ctx context.Context, id ResourceGroupsId) (*GetResourceGroupsResponse, error) {
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // ok
		},
		Uri: sdk.BuildResourceManagerURI(id, client.apiVersion),
	}

	var out GetResourceGroup
	resp, err := client.baseClient.GetJson(ctx, req, &out)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	result := GetResourceGroupsResponse{
		HttpResponse:   resp,
		ResourceGroups: &out,
	}
	return &result, nil
}

func (client ResourceGroupsClient) Update(ctx context.Context, id ResourceGroupsId, input UpdateInput) error {
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

func (client ResourceGroupsClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.Resources"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
