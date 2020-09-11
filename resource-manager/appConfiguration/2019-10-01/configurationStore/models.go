package configurationStore

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type CreateStoreInput struct {
	Location string             `json:"location"`
	Sku      Sku                `json:"sku"`
	Tags     *map[string]string `json:"tags,omitempty"`
}

func (m CreateStoreInput) Validate() error {
	var result error

	if m.Location == "" {
		result = multierror.Append(result, fmt.Errorf("Location cannot be empty"))
	}

	return result
}

type GetStore struct {
	Location   string             `json:"location"`
	Properties GetStoreProperties `json:"properties"`
	Sku        Sku                `json:"sku"`
	Tags       map[string]string  `json:"tags"`
}

type GetStoreProperties struct {
	ConfigurationStoreEndpoint string `json:"endpoint"`
}

type Sku struct {
	Name SkuName `json:"name"`
}
