package sdk

import (
	"context"
	"fmt"
	"net/http"
)

type Poller interface {
	GetOriginalResponse() *http.Response
	GetLatestPollResponse() *http.Response
	PollUntilDone(ctx context.Context) error
}

func DeterminePoller(response *http.Response, baseClient *BaseClient, uri string) (Poller, error) {
	// we could clearly make this smarter, but this is fine for now
	poller, err := newLongRunningOperationPoller(response, baseClient)
	if err == nil {
		return poller, nil
	}

	poller, err = newProvisioningStatePoller(response, baseClient, uri)
	if err == nil {
		return poller, nil
	}

	return nil, fmt.Errorf("unable to determine poller type")
}
