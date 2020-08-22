package resourcegroups

import "fmt"

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
