package namespaces

import "fmt"

type EventHubNamespaceId struct {
	SubscriptionId string
	ResourceGroup  string
	Namespace      string
}

func NewEventHubNamespaceId(subscriptionId string, resourceGroup string, namespace string) EventHubNamespaceId {
	return EventHubNamespaceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Namespace:      namespace,
	}
}

func (id EventHubNamespaceId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s", id.SubscriptionId, id.ResourceGroup, id.Namespace)
}
