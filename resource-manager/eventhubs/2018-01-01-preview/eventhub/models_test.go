package eventhub

import "github.com/tombuildsstuff/pandora/sdk"

var _ sdk.ModelWithValidation = CaptureDestination{}

var _ sdk.ModelWithValidation = CaptureDestinationProperties{}

var _ sdk.ModelWithValidation = CreateEventHubProperties{}

// TODO: unit tests for the API methods based on sample responses
