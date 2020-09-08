package groups

import "fmt"

type ResourceGroupId struct {
	ResourceGroup string
}

func NewResourceGroupId(resourceGroup string) ResourceGroupId {
	return ResourceGroupId{
		ResourceGroup: resourceGroup,
	}
}

func (id ResourceGroupId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, id.ResourceGroup)
}
