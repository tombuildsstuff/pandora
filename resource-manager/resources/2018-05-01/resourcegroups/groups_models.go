package resourcegroups

import (
	"fmt"
	"net/http"
)

type CreateResourceGroupInput struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

func (input CreateResourceGroupInput) Validate() error {
	errors := make([]error, 0)

	if input.Location == "" {
		errors = append(errors, fmt.Errorf("`location` cannot be empty"))
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("errors: %+v", errors)
}

type GetResourceGroup struct {
	Location string            `json:"location"`
	Tags     map[string]string `json:"tags"`
}

type GetResourceGroupResponse struct {
	HttpResponse  *http.Response
	ResourceGroup *GetResourceGroup
}

type UpdateResourceGroupInput struct {
	Tags *map[string]string `json:"tags,omitempty"`
}
