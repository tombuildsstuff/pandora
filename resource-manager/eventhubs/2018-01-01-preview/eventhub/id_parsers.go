package eventhub

import "fmt"

type EventhubId struct {
	SubscriptionId string
	ResourceGroup  string
	Namespace      string
	EventHub       string
}

func NewEventhubId(subscriptionId string, resourceGroup string, namespace string, eventHub string) EventhubId {
	return EventhubId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Namespace:      namespace,
		EventHub:       eventHub,
	}
}

func (id EventhubId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s", id.SubscriptionId, id.ResourceGroup, id.Namespace, id.EventHub)
}
