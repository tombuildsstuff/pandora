package templates

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tombuildsstuff/pandora/generator/models"
)

type ModelsTemplater struct {
	packageName string
	typeName    string
	operations  []models.OperationMetaData
}

func NewModelsTemplater(packageName, typeName string, operations []models.OperationMetaData) ModelsTemplater {
	return ModelsTemplater{
		packageName: packageName,
		typeName:    typeName,
		operations:  operations,
	}
}

func (t ModelsTemplater) Build() (*string, error) {
	models, err := t.models()
	if err != nil {
		return nil, fmt.Errorf("building models: %+v", err)
	}
	template := fmt.Sprintf(`package %[1]s

import (
	"fmt"
	"net/http"
)

%[2]s`, t.packageName, *models)
	return &template, nil
}

func (t ModelsTemplater) models() (*string, error) {
	// first collate all of the types from all operations
	// then sort them and output them
	types := make(map[string]string, 0)

	for _, operation := range t.operations {
		newTypes, err := t.typesForOperation(operation, t.typeName)
		if err != nil {
			return nil, fmt.Errorf("building types for %q (method %q)", operation.Name, operation.Method)
		}
		for k, v := range *newTypes {
			// ensure no duplicates, which should be impossible but defensive against bugs
			if _, existing := types[k]; existing {
				return nil, fmt.Errorf("invalid duplicate type for %q", k)
			}

			types[k] = v
		}
	}

	sortedKeys := make([]string, 0)
	for k, _ := range types {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	sortedStructs := make([]string, 0)
	for _, key := range sortedKeys {
		value := types[key]
		sortedStructs = append(sortedStructs, value)
	}

	result := strings.Join(sortedStructs, "\n\n")
	return &result, nil
}

func (t ModelsTemplater) typesForOperation(input models.OperationMetaData, typeName string) (*map[string]string, error) {
	method := strings.ToUpper(input.Method)
	if method == "DELETE" {
		// TODO: some operations do have some querystrings however
		return &map[string]string{}, nil
	}

	// TODO: validation methods

	if method == "GET" {
		result := t.getOperationTypes(input, typeName)
		return &result, nil
	}

	if method == "PATCH" {
		result := t.patchOperationTypes(input, typeName)
		return &result, nil
	}

	if method == "PUT" {
		result := t.putOperationTypes(input, typeName)
		return &result, nil
	}

	// TODO: temp
	return nil, fmt.Errorf("unsupported method %q", input.Method)
}

func (t ModelsTemplater) getOperationTypes(input models.OperationMetaData, typeName string) map[string]string {
	structName := fmt.Sprintf("%s%s", input.Name, typeName)
	wrapperStructName := fmt.Sprintf("%sResponse", structName)
	return map[string]string{
		structName: fmt.Sprintf(`type %s struct {
	// TODO: implementation
}`, structName),
		wrapperStructName: fmt.Sprintf(`type %[1]s struct {
	HttpResponse  *http.Response
	ResourceGroup *%[2]s
}`, wrapperStructName, structName),
	}
}

func (t ModelsTemplater) patchOperationTypes(input models.OperationMetaData, typeName string) map[string]string {
	structName := fmt.Sprintf("%s%sInput", input.Name, typeName)
	return map[string]string{
		structName: fmt.Sprintf(`type %s struct {
	// TODO: implementation
	// TODO: notably here all fields need to be 'omitempty'
}`, structName),
	}
}

func (t ModelsTemplater) putOperationTypes(input models.OperationMetaData, typeName string) map[string]string {
	structName := fmt.Sprintf("%s%sInput", input.Name, typeName)
	return map[string]string{
		structName: fmt.Sprintf(`type %s struct {
	// TODO: implementation
}`, structName),
	}
}
