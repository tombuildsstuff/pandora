package configurationStore

import "fmt"

type ConfigurationStoreId struct {
	SubscriptionId  string
	ResourceGroup   string
	ConfigStoreName string
}

func NewConfigurationStoreId(subscriptionId string, resourceGroup string, configStoreName string) ConfigurationStoreId {
	return ConfigurationStoreId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ConfigStoreName: configStoreName,
	}
}

func (id ConfigurationStoreId) ID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s", id.SubscriptionId, id.ResourceGroup, id.ConfigStoreName)
}
