package sdk

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Poller struct {
	OriginalResponse   *http.Response
	LatestPollResponse *http.Response

	baseClient   *BaseClient
	pollInterval time.Duration
	pollLocation string
}

func NewResourceManagerPoller(response *http.Response, baseClient *BaseClient) (*Poller, error) {
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	if response.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("status code %d (%s) is not a long running operation", response.StatusCode, response.Status)
	}

	locationHeader := response.Header.Get("Azure-AsyncOperation")
	if locationHeader == "" {
		locationHeader = response.Header.Get("Location")
		if locationHeader == "" {
			return nil, fmt.Errorf("the `Azure-AsyncOperation` and `Location` headers were empty")
		}
	}

	retryAfter := 15
	if retryAfterHeader := response.Header.Get("Retry-After"); retryAfterHeader != "" {
		parsed, err := strconv.ParseInt(retryAfterHeader, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as an int: %+v", retryAfterHeader, err)
		}

		retryAfter = int(parsed)
	}

	return &Poller{
		OriginalResponse:   response,
		LatestPollResponse: nil,

		baseClient:   baseClient,
		pollInterval: time.Duration(retryAfter) * time.Second,
		pollLocation: locationHeader,
	}, nil
}

func (p *Poller) PollUntilDone(ctx context.Context) error {
	for {
		// wait for the recommended amount of settings before continuing
		time.Sleep(p.pollInterval)

		input := GetHttpRequestInput{
			Uri: p.pollLocation,
			ExpectedStatusCodes: []int{
				http.StatusAccepted, // in progress
				http.StatusOK,       // finished
			},
		}

		var err error
		p.LatestPollResponse, err = p.baseClient.Get(ctx, input)
		if err != nil {
			return fmt.Errorf("polling: %+v", err)
		}

		// we should be done
		if p.LatestPollResponse.StatusCode == http.StatusOK {
			return nil
		}

		// keep waiting
		if p.LatestPollResponse.StatusCode == http.StatusAccepted {
			// TODO: we could parse the location/retry-after header out, but it appears unnecessary
			continue
		}

		break
	}

	return nil
}
