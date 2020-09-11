package namespaces

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type CreateNamespaceInput struct {
	Location   string                    `json:"location"`
	Properties CreateNamespaceProperties `json:"properties"`
	Sku        Sku                       `json:"sku"`
	Tags       *map[string]string        `json:"tags,omitempty"`
}

func (m CreateNamespaceInput) Validate() error {
	var result error

	if m.Location == "" {
		result = multierror.Append(result, fmt.Errorf("Location cannot be empty"))
	}

	return result
}

type CreateNamespaceProperties struct {
	IsAutoInflateEnabled *bool `json:"isAutoInflateEnabled,omitempty"`
	ZoneRedundant        *bool `json:"ZoneRedundant,omitempty"`
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

type Sku struct {
	Capacity *int64  `json:"capacity,omitempty"`
	Name     SkuTier `json:"name"`
	Tier     SkuTier `json:"tier"`
}

type UpdateNamespaceInput struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
