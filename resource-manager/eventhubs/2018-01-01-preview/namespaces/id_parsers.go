package namespaces

import "fmt"

type EventHubNamespaceId struct {
	ResourceGroup string
	Namespace     string
}

func NewEventHubNamespaceId(resourceGroup string, namespace string) EventHubNamespaceId {
	return EventHubNamespaceId{
		ResourceGroup: resourceGroup,
		Namespace:     namespace,
	}
}

func (id EventHubNamespaceId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s", subscriptionId, id.ResourceGroup, id.Namespace)
}
