package eventhub

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type CaptureDescription struct {
	CaptureDestination *CaptureDestination `json:"destination,omitempty"`
	Enabled            bool                `json:"enabled"`
	Encoding           *Encoding           `json:"encoding,omitempty"`
	IntervalInSeconds  *int64              `json:"intervalInSeconds,omitempty"`
	SizeLimitInBytes   *int64              `json:"sizeLimitsInBytes,omitempty"`
	SkipEmptyArchives  *bool               `json:"skipEmptyArchives,omitempty"`
}

type CaptureDestination struct {
	Name       string                       `json:"name"`
	Properties CaptureDestinationProperties `json:"properties"`
}

func (m CaptureDestination) Validate() error {
	var result error

	if m.Name == "" {
		result = multierror.Append(result, fmt.Errorf("Name cannot be empty"))
	}

	return result
}

type CaptureDestinationProperties struct {
	ArchiveNameFormat        string `json:"archiveNameFormat"`
	BlobContainerName        string `json:"blobContainer"`
	StorageAccountResourceID string `json:"storageAccountResourceId"`
}

func (m CaptureDestinationProperties) Validate() error {
	var result error

	if m.BlobContainerName == "" {
		result = multierror.Append(result, fmt.Errorf("BlobContainerName cannot be empty"))
	}

	if m.StorageAccountResourceID == "" {
		result = multierror.Append(result, fmt.Errorf("StorageAccountResourceID cannot be empty"))
	}

	return result
}

type CreateEventHubInput struct {
	Properties CreateEventHubProperties `json:"properties"`
}

type CreateEventHubProperties struct {
	CaptureDescription     *CaptureDescription `json:"captureDescription,omitempty"`
	MessageRetentionInDays *int64              `json:"messageRetentionInDays,omitempty"`
	PartitionCount         *int64              `json:"partitionCount,omitempty"`
}

func (m CreateEventHubProperties) Validate() error {
	var result error

	// TODO: range validation

	// TODO: range validation

	return result
}

type GetEventHub struct {
	Properties GetEventHubProperties `json:"properties"`
}

type GetEventHubProperties struct {
	Status EntityStatus `json:"status"`
}
