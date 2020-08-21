package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ApiClient interface {
	MetaData() ClientMetaData
}

type ArmResourceId interface {
	ID(subscriptionId string) string
}

type ClientMetaData struct {
	ResourceProvider *string
}

type BaseClient struct {
	authorizer Authorizer
	endpoint   string
	httpClient *http.Client
}

func DefaultBaseClient(endpoint string, authorizer Authorizer) BaseClient {
	return BaseClient{
		authorizer: authorizer,
		endpoint:   endpoint,
		httpClient: &http.Client{
			Transport: http.DefaultTransport,
		},
	}
}

type DeleteHttpRequestInput struct {
	ExpectedStatusCodes []int
	Uri                 string
}

func (c BaseClient) Delete(ctx context.Context, input DeleteHttpRequestInput) (*http.Response, error) {
	url := fmt.Sprintf("https://%s%s", c.endpoint, input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := c.performAuthenticatedHttpRequest(ctx, req, input.ExpectedStatusCodes)
	if err != nil {
		return nil, fmt.Errorf("making request: %+v", err)
	}

	return resp, nil
}

type GetHttpRequestInput struct {
	ExpectedStatusCodes []int
	Uri                 string
}

func (c BaseClient) Get(ctx context.Context, input GetHttpRequestInput) (*http.Response, error) {
	url, err := c.buildUri(input.Uri)
	if err != nil {
		return nil, fmt.Errorf("building uri: %+v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := c.performAuthenticatedHttpRequest(ctx, req, input.ExpectedStatusCodes)
	if err != nil {
		return resp, fmt.Errorf("making request: %+v", err)
	}

	return resp, nil
}

func (c BaseClient) GetJson(ctx context.Context, input GetHttpRequestInput, out interface{}) (*http.Response, error) {
	url, err := c.buildUri(input.Uri)
	if err != nil {
		return nil, fmt.Errorf("building uri: %+v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, *url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := c.performAuthenticatedHttpRequest(ctx, req, input.ExpectedStatusCodes)
	if err != nil {
		return resp, fmt.Errorf("making request: %+v", err)
	}

	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		return resp, fmt.Errorf("expected the 'Content-Type' to be 'application/json' but got %q", contentType)
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return resp, fmt.Errorf("unmarshalling response: %+v", err)
	}

	return resp, nil
}

type PatchHttpRequestInput struct {
	Body                interface{}
	ExpectedStatusCodes []int
	Uri                 string
}

func (c BaseClient) PatchJson(ctx context.Context, input PatchHttpRequestInput) (*http.Response, error) {
	marshalledBytes, err := json.Marshal(input.Body)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %+v", err)
	}

	url, err := c.buildUri(input.Uri)
	if err != nil {
		return nil, fmt.Errorf("building uri: %+v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, *url, bytes.NewReader(marshalledBytes))
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := c.performAuthenticatedHttpRequest(ctx, req, input.ExpectedStatusCodes)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type PutHttpRequestInput struct {
	Body                interface{}
	ExpectedStatusCodes []int
	Uri                 string
}

func (c BaseClient) PutJson(ctx context.Context, input PutHttpRequestInput) (*http.Response, error) {
	marshalledBytes, err := json.Marshal(input.Body)
	if err != nil {
		return nil, fmt.Errorf("marshalling body: %+v", err)
	}

	url, err := c.buildUri(input.Uri)
	if err != nil {
		return nil, fmt.Errorf("building uri: %+v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, *url, bytes.NewReader(marshalledBytes))
	if err != nil {
		return nil, fmt.Errorf("building request: %+v", err)
	}

	resp, err := c.performAuthenticatedHttpRequest(ctx, req, input.ExpectedStatusCodes)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (c BaseClient) performAuthenticatedHttpRequest(ctx context.Context, req *http.Request, expectedStatusCodes []int) (*http.Response, error) {
	token, err := c.authorizer.Token(ctx, "https://management.azure.com")
	if err != nil {
		return nil, fmt.Errorf("retrieving auth token: %+v", err)
	}
	req.Header.Add("Authorization", token.AuthorizationHeader())
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	// TODO: handle retries, 429's etc
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("sending request: %+v", err)
	}

	if exists := containsStatusCode(expectedStatusCodes, resp.StatusCode); !exists {
		return resp, fmt.Errorf("unexpected status %d (%s)", resp.StatusCode, resp.Status)
	}

	return resp, nil
}

func (c BaseClient) buildUri(input string) (*string, error) {
	uri, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	// it's a full URI so let's use it
	if uri.IsAbs() {
		return &input, nil
	}

	output := fmt.Sprintf("https://%s%s", c.endpoint, input)
	return &output, nil
}

func containsStatusCode(expected []int, actual int) bool {
	for _, v := range expected {
		if actual == v {
			return true
		}
	}

	return false
}
