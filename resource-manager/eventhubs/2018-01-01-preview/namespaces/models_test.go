package namespaces

import "github.com/tombuildsstuff/pandora/sdk"

var _ sdk.ModelWithValidation = CreateNamespaceInput{}

var _ sdk.ModelWithValidation = Sku{}

var _ sdk.ModelWithValidation = UpdateNamespaceInput{}

// TODO: unit tests for the API methods based on sample responses
