package templates

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tombuildsstuff/pandora/generator/models"
)

type ClientTemplater struct {
	packageName      string
	typeName         string
	resourceProvider *string
	apiVersion       string
	operations       []models.OperationMetaData
}

func NewClientTemplater(packageName, typeName, apiVersion string, resourceProvider *string, operations []models.OperationMetaData) ClientTemplater {
	return ClientTemplater{
		packageName:      packageName,
		typeName:         typeName,
		apiVersion:       apiVersion,
		resourceProvider: resourceProvider,
		operations:       operations,
	}
}

func (t ClientTemplater) Build() (*string, error) {
	constructors := t.constructors()
	metadata := t.metadata()
	methods, err := t.methods()
	if err != nil {
		return nil, fmt.Errorf("generating methods: %+v", err)
	}

	template := fmt.Sprintf(`package %[1]s

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type %[2]ssClient struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string
}

%[3]s

%[4]s

%[5]s`, t.packageName, t.typeName, constructors, *methods, metadata)
	return &template, nil
}

func (t ClientTemplater) constructors() string {
	template := fmt.Sprintf(`
func New%[1]ssClient(subscriptionId string, authorizer sdk.Authorizer) %[1]ssClient {
	return New%[1]ssClientWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func New%[1]ssClientWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) %[1]ssClient {
	return %[1]ssClient{
		apiVersion:     "%s",
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}
`, t.typeName, t.apiVersion)
	return strings.TrimSpace(template)
}

func (t ClientTemplater) metadata() string {
	// note: this is lazy but it avoids requiring an additional optional import (utils.String)
	template := ""
	if t.resourceProvider != nil {
		template = fmt.Sprintf(`
func (client %[1]ssClient) MetaData() sdk.ClientMetaData {
	resourceProvider := "%[2]s"
	return sdk.ClientMetaData{
		ResourceProvider: &resourceProvider,
	}
}
`, t.typeName, *t.resourceProvider)
	} else {
		template = fmt.Sprintf(`
func (client %[1]ssClient) MetaData() sdk.ClientMetaData {
	return sdk.ClientMetaData{}
}
`, t.typeName)
	}

	return strings.TrimSpace(template)
}

func (t ClientTemplater) methods() (*string, error) {
	output := make([]string, 0)

	sortedMethods := sortMethodsAlphabetically(t.operations)
	for _, method := range sortedMethods {
		templater := methodTemplater{
			typeName:             t.typeName,
			method:               method.Method,
			name:                 method.Name,
			longRunningOperation: method.LongRunningOperation,
			expectedStatusCodes:  method.ExpectedStatusCodes,
		}
		formatted, err := templater.Build()
		if err != nil {
			return nil, err
		}

		output = append(output, *formatted)
	}

	result := strings.Join(output, "\n\n")
	return &result, nil
}

func sortMethodsAlphabetically(input []models.OperationMetaData) []models.OperationMetaData {
	names := make([]string, 0)
	indexes := make(map[string]int, len(input))
	for i, v := range input {
		names = append(names, v.Name)
		indexes[v.Name] = i
	}

	sort.Strings(names)
	out := make([]models.OperationMetaData, 0)
	for _, v := range names {
		index := indexes[v]
		out = append(out, input[index])
	}

	return out
}
