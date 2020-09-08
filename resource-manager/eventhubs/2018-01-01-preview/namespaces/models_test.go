package namespaces

import "github.com/tombuildsstuff/pandora/sdk"

var _ sdk.ModelWithValidation = CreateNamespaceInput{}

var _ sdk.ModelWithValidation = PatchNamespaceInput{}

var _ sdk.ModelWithValidation = Sku{}

// TODO: unit tests for the API methods based on sample responses
