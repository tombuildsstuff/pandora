package namespaces

import "fmt"

type NamespacesId struct {
	SubscriptionId string
	ResourceGroup  string
	Namespace      string
}

func NewNamespacesId(subscriptionId string, resourceGroup string, namespace string) NamespacesId {
	return NamespacesId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Namespace:      namespace,
	}
}

func (id NamespacesId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s", id.SubscriptionId, id.ResourceGroup, id.Namespace)
}
