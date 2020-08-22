package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ProvisioningStatePoller struct {
	baseClient         *BaseClient
	latestPollResponse *http.Response
	originalResponse   *http.Response
	pollInterval       time.Duration
	pollLocation       string
}

func (p ProvisioningStatePoller) GetLatestPollResponse() *http.Response {
	return p.latestPollResponse
}

func (p ProvisioningStatePoller) GetOriginalResponse() *http.Response {
	return p.originalResponse
}

func newProvisioningStatePoller(response *http.Response, baseClient *BaseClient, uri string) (Poller, error) {
	return &ProvisioningStatePoller{
		baseClient:       baseClient,
		originalResponse: response,
		pollInterval:     15 * time.Second,
		pollLocation:     uri,
	}, nil
}

func (p *ProvisioningStatePoller) PollUntilDone(ctx context.Context) error {
	for {
		// wait for the recommended amount of settings before continuing
		time.Sleep(p.pollInterval)

		input := GetHttpRequestInput{
			Uri: p.pollLocation,
			ExpectedStatusCodes: []int{
				http.StatusOK,
			},
		}

		var err error
		p.latestPollResponse, err = p.baseClient.Get(ctx, input)
		if err != nil {
			return fmt.Errorf("polling: %+v", err)
		}

		var out ProvisioningStateResponse
		if err := json.NewDecoder(p.latestPollResponse.Body).Decode(&out); err != nil {
			return fmt.Errorf("decoding response: %+v", err)
		}

		if strings.EqualFold(out.Properties.ProvisioningState, "Succeeded") {
			return nil
		}

		// TODO: handle failed

		continue
	}

	return nil
}

type ProvisioningStateResponse struct {
	Properties ProvisioningStateResponseProperties `json:"properties"`
}

type ProvisioningStateResponseProperties struct {
	ProvisioningState string `json:"provisioningState`
}
