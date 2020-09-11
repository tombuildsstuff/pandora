package resourceGroups

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type CreateInput struct {
	Location string             `json:"location"`
	Tags     *map[string]string `json:"tags,omitempty"`
}

func (m CreateInput) Validate() error {
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

type UpdateInput struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
