package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PandoraApi interface {
	Apis() (*ApisResponse, error)
	VersionsForApi(api ApiReference) (*ApiVersionsResponse, error)
	MetaDataForOperation(definition Operation) (*ApiOperationMetaData, error)
	OperationsForApiVersion(version VersionDetails) (*ApiVersionOperationsResponse, error)
	APIOperationsForApiVersion(operation ApiOperationMetaData) (*OperationsResponse, error)
	SchemasForApiVersion(operation ApiOperationMetaData) (*SchemaResponse, error)
}

type PandoraApiService struct {
	endpoint        string
	resourceManager bool
}

func NewDataPlaneService(endpoint string) PandoraApiService {
	return PandoraApiService{
		endpoint:        endpoint,
		resourceManager: false,
	}
}

func NewResourceManagerService(endpoint string) PandoraApiService {
	return PandoraApiService{
		endpoint:        endpoint,
		resourceManager: true,
	}
}

func (rm PandoraApiService) segment() string {
	if rm.resourceManager {
		return "resource-manager"
	}

	return "data-plane"
}

func (rm PandoraApiService) Apis() (*ApisResponse, error) {
	segment := rm.segment()
	resp, err := rm.getJson(fmt.Sprintf("%s/apis/v1/%s", rm.endpoint, segment))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out ApisResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

func (rm PandoraApiService) VersionsForApi(api ApiReference) (*ApiVersionsResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, api.Uri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out ApiVersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

func (rm PandoraApiService) OperationsForApiVersion(version VersionDetails) (*ApiVersionOperationsResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, version.Uri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out ApiVersionOperationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

func (rm PandoraApiService) MetaDataForOperation(definition Operation) (*ApiOperationMetaData, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, definition.Uri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out ApiOperationMetaData
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

// APIOperationsForApiVersion returns the Operations supported by this Type for this Api version
func (rm PandoraApiService) APIOperationsForApiVersion(operation ApiOperationMetaData) (*OperationsResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, operation.OperationsUri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out OperationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

// SchemasForApiVersion returns the Schema supported by this Type for this API version
func (rm PandoraApiService) SchemasForApiVersion(operation ApiOperationMetaData) (*SchemaResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, operation.SchemaUri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out SchemaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

func (rm PandoraApiService) getJson(endpoint string) (*http.Response, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}
