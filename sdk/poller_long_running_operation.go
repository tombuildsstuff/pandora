package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LongRunningOperationPoller struct {
	baseClient         *BaseClient
	latestPollResponse *http.Response
	originalResponse   *http.Response
	pollInterval       time.Duration
	pollLocation       string
}

func (p LongRunningOperationPoller) GetLatestPollResponse() *http.Response {
	return p.latestPollResponse
}

func (p LongRunningOperationPoller) GetOriginalResponse() *http.Response {
	return p.originalResponse
}

func newLongRunningOperationPoller(response *http.Response, baseClient *BaseClient) (Poller, error) {
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusAccepted {
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

	return &LongRunningOperationPoller{
		baseClient:       baseClient,
		originalResponse: response,
		pollInterval:     time.Duration(retryAfter) * time.Second,
		pollLocation:     locationHeader,
	}, nil
}

func (p *LongRunningOperationPoller) PollUntilDone(ctx context.Context) error {
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
		p.latestPollResponse, err = p.baseClient.Get(ctx, input)
		if err != nil {
			return fmt.Errorf("polling: %+v", err)
		}

		// we should be done
		if p.latestPollResponse.StatusCode == http.StatusOK {
			var details longRunningOperationResponse
			if json.NewDecoder(p.latestPollResponse.Body).Decode(&details); err != nil {
				return fmt.Errorf("decoding latest polling response")
			}

			if details.Status == "" || strings.EqualFold(details.Status, "Succeeded") {
				return nil
			}

			if strings.EqualFold(details.Status, "Failed") {
				if details.Error != nil {
					return fmt.Errorf("polling for status (Code %q / Message %q)", details.Error.Code, details.Error.Message)
				}

				return fmt.Errorf("polling for status: %q", details.Status)
			}

			return fmt.Errorf("unexpected status %q", details.Status)
		}

		continue
	}

	return nil
}

type longRunningOperationResponse struct {
	Id     string                     `json:"id"`
	Name   string                     `json:"name"`
	Status string                     `json:"status"`
	Error  *longRunningOperationError `json:"error,omitempty"'`
}

type longRunningOperationError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
