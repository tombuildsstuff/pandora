package resourceGroups

import "fmt"

type ResourceGroupsId struct {
	SubscriptionId string
	ResourceGroup  string
}

func NewResourceGroupsId(subscriptionId string, resourceGroup string) ResourceGroupsId {
	return ResourceGroupsId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
	}
}

func (id ResourceGroupsId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", id.SubscriptionId, id.ResourceGroup)
}
