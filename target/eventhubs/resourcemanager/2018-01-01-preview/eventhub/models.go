package eventhub

import "net/http"

type SkuName string

const (
	Basic    SkuName = "Basic"
	Standard SkuName = "Standard"
)

type Sku struct {
	Name     SkuName `json:"name,omitempty"`
	Tier     SkuName `json:"tier,omitempty"`
	Capacity *int    `json:"capacity,omitempty"`
}

type CreateNamespaceInput struct {
	Location   string                    `json:"location"`
	Properties CreateNamespaceProperties `json:"properties"`
	Sku        Sku                       `json:"sku"`
	Tags       map[string]string         `json:"tags"`
}

type CreateNamespaceProperties struct {
	IsAutoInflateEnabled bool `json:"isAutoInflateEnabled"`
	ZoneRedundant        bool `json:"zoneRedundant"`
}

type GetNamespaceResponse struct {
	HttpResponse *http.Response
	Namespace    *GetNamespace
}

type GetNamespace struct {
	Location   string                 `json:"location"`
	Properties GetNamespaceProperties `json:"properties"`
	Sku        Sku                    `json:"sku"`
	Tags       map[string]string      `json:"tags"`
}

type GetNamespaceProperties struct {
	IsAutoInflateEnabled bool   `json:"isAutoInflateEnabled"`
	ServiceBusEndpoint   string `json:"serviceBusEndpoint"`
	ZoneRedundant        bool   `json:"zoneRedundant"`
}
