package consumerGroup

import "fmt"

type ConsumerGroupId struct {
	SubscriptionId string
	ResourceGroup  string
	Namespaces     string
	EventHub       string
	ConsumerGroup  string
}

func NewConsumerGroupId(subscriptionId string, resourceGroup string, namespaces string, eventHub string, consumerGroup string) ConsumerGroupId {
	return ConsumerGroupId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Namespaces:     namespaces,
		EventHub:       eventHub,
		ConsumerGroup:  consumerGroup,
	}
}

func (id ConsumerGroupId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventHub/namespaces/%s/eventhubs/%s/consumerGroups/%s", id.SubscriptionId, id.ResourceGroup, id.Namespaces, id.EventHub, id.ConsumerGroup)
}
