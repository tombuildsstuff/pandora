package namespaces

import "fmt"

type NamespacesId struct {
	ResourceGroup string
	Namespace     string
}

func NewNamespacesId(resourceGroup string, namespace string) NamespacesId {
	return NamespacesId{
		ResourceGroup: resourceGroup,
		Namespace:     namespace,
	}
}

func (id NamespacesId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s", subscriptionId, id.ResourceGroup, id.Namespace)
}
