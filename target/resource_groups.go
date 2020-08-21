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

func (client ResourceGroupsClient) Create(ctx context.Context, name string, input CreateResourceGroupInput) error {
	uri := NewResourceGroupID(name).ID(client.subscriptionId)
	req := sdk.PutHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK,      // already exists
			http.StatusCreated, // new
		},
		Uri: fmt.Sprintf("%s?api-version=%s", uri, client.apiVersion), // TODO: some kind of helper but whatever
	}
	if _, err := client.baseClient.PutJson(ctx, req); err != nil {
		return fmt.Errorf("sending Request: %+v", err)
	}
	return nil
}

func (client ResourceGroupsClient) Delete(ctx context.Context, id ResourceGroupID) (*sdk.Poller, error) {
	uri := id.ID(client.subscriptionId)
	req := sdk.DeleteHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusAccepted, // delete accepted
		},
		Uri: fmt.Sprintf("%s?api-version=%s", uri, client.apiVersion), // TODO: some kind of helper but whatever
	}
	originalResp, err := client.baseClient.Delete(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sending Request: %+v", err)
	}

	poller, err := sdk.NewResourceManagerPoller(originalResp, &client.baseClient)
	if err != nil {
		return nil, fmt.Errorf("building poller: %+v", err)
	}

	return poller, nil
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
	uri := id.ID(client.subscriptionId)
	req := sdk.GetHttpRequestInput{
		ExpectedStatusCodes: []int{
			http.StatusOK, // Exists
		},
		Uri: fmt.Sprintf("%s?api-version=%s", uri, client.apiVersion), // TODO: some kind of helper but whatever
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
	uri := id.ID(client.subscriptionId)
	req := sdk.PatchHttpRequestInput{
		Body: input,
		ExpectedStatusCodes: []int{
			http.StatusOK, // already exists
		},
		Uri: fmt.Sprintf("%s?api-version=%s", uri, client.apiVersion), // TODO: some kind of helper but whatever
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
