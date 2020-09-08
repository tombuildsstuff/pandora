package services

import (
	"fmt"
	"net/http"

	"github.com/rust-lang/rust/src/llvm-project/llgo/third_party/gofrontend/libgo/go/encoding/json"
)

type ResourceManagerService struct {
	endpoint string
}

func NewResourceManagerService(endpoint string) ResourceManagerService {
	return ResourceManagerService{
		endpoint: endpoint,
	}
}

// TODO: rename this functions and structs

// SupportedApis returns the supported Resource Manager API's available in Pandora
func (rm ResourceManagerService) SupportedApis() (*ResourceManagerApiResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s/apis/v1/resource-manager", rm.endpoint))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out ResourceManagerApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

// SupportedVersionsForApi returns the supported API versions for this Api definition
func (rm ResourceManagerService) SupportedVersionsForApi(api ApiDetails) (*SupportedVersionsResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, api.Uri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out SupportedVersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

// OperationsForApiVersion returns the supported operation types for this Api version
func (rm ResourceManagerService) OperationsForApiVersion(version VersionDetails) (*SupportedTypesResponse, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, version.Uri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out SupportedTypesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

// MetaDataForOperation returns the metadata for this API version
func (rm ResourceManagerService) MetaDataForOperation(definition TypeDefinition) (*OperationMetaData, error) {
	resp, err := rm.getJson(fmt.Sprintf("%s%s", rm.endpoint, definition.Uri))
	if err != nil {
		return nil, fmt.Errorf("retrieving JSON: %+v", err)
	}

	var out OperationMetaData
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decoding JSON: %+v", err)
	}

	return &out, nil
}

// OperationsForType returns the Operations supported by this Type for this Api version
func (rm ResourceManagerService) OperationsForType(operation OperationMetaData) (*OperationsResponse, error) {
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

// SchemaForType returns the Schema supported by this Type for this API version
func (rm ResourceManagerService) SchemaForType(operation OperationMetaData) (*SchemaResponse, error) {
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

func (rm ResourceManagerService) getJson(endpoint string) (*http.Response, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d %s", resp.StatusCode, resp.Status)
	}

	return resp, nil
}
