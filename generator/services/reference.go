package services

import (
	"fmt"
	"strings"
)

type reference struct {
	operationId string
	name        string
}

func parseReference(input string) (*reference, error) {
	// format:
	//	"/apis/v1/resource-manager/eventhubs/2018-01-01-preview/namespaces/schema#Sku"
	segments := strings.Split(input, "#")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected 2 segments but got %d (input %q)", len(segments), input)
	}

	// this exists to allow cross-model references in future (e.g. resource groups)
	// this likely won't be used anytime soon, but is handy to have available from the start
	// to be able to look these up as required
	ref := reference{
		operationId: segments[0],
		name:        segments[1],
	}
	return &ref, nil
}
