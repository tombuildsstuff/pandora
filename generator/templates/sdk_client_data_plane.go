package templates

import (
	"fmt"
	"sort"
	"strings"
)

type DataPlaneClientTemplater struct {
	packageName string
	typeName    string
	apiVersion  string
	operations  []ClientOperation
}

func NewDataPlaneClientTemplater(packageName, typeName, apiVersion string, operations []ClientOperation) DataPlaneClientTemplater {
	return DataPlaneClientTemplater{
		packageName: packageName,
		typeName:    typeName,
		apiVersion:  apiVersion,
		operations:  operations,
	}
}

func (t DataPlaneClientTemplater) Build() (*string, error) {
	clientName := fmt.Sprintf("%sClient", strings.Title(t.typeName))
	constructors, err := t.constructors(clientName)
	if err != nil {
		return nil, fmt.Errorf("generating constructors: %+v", err)
	}

	methods, err := t.methods(clientName)
	if err != nil {
		return nil, fmt.Errorf("generating methods: %+v", err)
	}

	metadata := t.metadata(clientName)

	out := fmt.Sprintf(`package %s

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/sdk/endpoints"
)

type %s struct {
	apiVersion     string
	baseClient     sdk.BaseClient
	subscriptionId string // TODO: making this Optional?
}

%s

%s

%s
`, t.packageName, clientName, *constructors, *methods, metadata)
	return &out, nil
}

func (t DataPlaneClientTemplater) constructors(clientName string) (*string, error) {
	format := fmt.Sprintf(`
func New%[1]s(subscriptionId string, authorizer sdk.Authorizer) %[1]s {
	return New%[1]sWithBaseURI(endpoints.DefaultManagementEndpoint, subscriptionId, authorizer)
}

func New%[1]sWithBaseURI(endpoint string, subscriptionId string, authorizer sdk.Authorizer) %[1]s {
	return %[1]s{
		apiVersion:     %[2]q,
		baseClient:     sdk.DefaultBaseClient(endpoint, authorizer),
		subscriptionId: subscriptionId,
	}
}
`, clientName, t.apiVersion)
	return &format, nil
}

func (t DataPlaneClientTemplater) methods(clientName string) (*string, error) {
	output := make([]string, 0)

	sortedMethods := t.sortMethods(t.operations)
	for _, method := range sortedMethods {
		templater := methodTemplater{
			clientName: clientName,
			typeName:   strings.Title(t.typeName),
			operation:  method,
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

func (t DataPlaneClientTemplater) metadata(clientName string) string {
	return fmt.Sprintf(`
func (client %[1]s) MetaData() sdk.ClientMetaData {
	return sdk.ClientMetaData{}
}
`, clientName)
}

func (t DataPlaneClientTemplater) sortMethods(input []ClientOperation) []ClientOperation {
	names := make([]string, 0)
	indexes := make(map[string]int, len(input))
	for i, v := range input {
		names = append(names, v.Name)
		indexes[v.Name] = i
	}

	sort.Strings(names)
	out := make([]ClientOperation, 0)
	for _, v := range names {
		index := indexes[v]
		out = append(out, input[index])
	}

	return out
}
