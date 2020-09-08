package groups

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type CreateResourceGroupInput struct {
	Location string             `json:"location"`
	Tags     *map[string]string `json:"tags,omitempty"`
}

func (m CreateResourceGroupInput) Validate() error {
	var result error

	if m.Location == "" {
		result = multierror.Append(result, fmt.Errorf("Location cannot be empty"))
	}

	return result
}

type GetResourceGroup struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type UpdateResourceGroupInput struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
