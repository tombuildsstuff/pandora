package eventhub

import "fmt"

type NamespaceID struct {
	Name string
	ResourceGroup string
}

func NewNamespaceID(resourceGroup, name string) NamespaceID {
	return NamespaceID{
		Name: name,
		ResourceGroup: resourceGroup,
	}
}

func (id NamespaceID) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s", subscriptionId, id.ResourceGroup, id.Name)
}
