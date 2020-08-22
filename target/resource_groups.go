package target

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
	subscriptionId string
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

type ResourceGroupID struct {
	Name string
}

func NewResourceGroupID(name string) ResourceGroupID {
	return ResourceGroupID{
		Name: name,
	}
}

func (id ResourceGroupID) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, id.Name)
}

type CreateResourceGroupInput struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

func (input CreateResourceGroupInput) Validate() error {
	errors := make([]error, 0)

	if input.Location == "" {
		errors = append(errors, fmt.Errorf("`location` cannot be empty"))
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("errors: %+v", errors)
}

func (client ResourceGroupsClient) Create(ctx context.Context, id ResourceGroupID, input CreateResourceGroupInput) error {
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

func (client ResourceGroupsClient) Delete(ctx context.Context, id ResourceGroupID) (sdk.Poller, error) {
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // delete accepted
		},
		Uri: sdk.BuildResourceManagerURI(id, client.subscriptionId, client.apiVersion),
	}

	return client.baseClient.DeleteThenPoll(ctx, req)
}

type GetResourceGroup struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type GetResourceGroupResponse struct {
	HttpResponse  *http.Response
	ResourceGroup *GetResourceGroup
}

func (client ResourceGroupsClient) Get(ctx context.Context, id ResourceGroupID) (*GetResourceGroupResponse, error) {
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

type UpdateResourceGroupInput struct {
	Tags *map[string]string `json:"tags,omitempty"`
}

func (client ResourceGroupsClient) Update(ctx context.Context, id ResourceGroupID, input UpdateResourceGroupInput) error {
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

func (client ResourceGroupsClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "Microsoft.Resources"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
