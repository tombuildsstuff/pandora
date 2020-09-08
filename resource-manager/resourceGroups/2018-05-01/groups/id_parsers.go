package groups

import "fmt"

type GroupsId struct {
	ResourceGroup string
}

func NewGroupsId(resourceGroup string) GroupsId {
	return GroupsId{
		ResourceGroup: resourceGroup,
	}
}

func (id GroupsId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, id.ResourceGroup)
}
