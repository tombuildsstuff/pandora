package groups

import "fmt"

type ResourceGroupId struct {
	SubscriptionId string
	ResourceGroup  string
}

func NewResourceGroupId(subscriptionId string, resourceGroup string) ResourceGroupId {
	return ResourceGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
	}
}

func (id ResourceGroupId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", id.SubscriptionId, id.ResourceGroup)
}
