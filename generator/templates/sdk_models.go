package templates

import (
	"fmt"
	"log"
	"strings"
)

type ModelsTemplater struct {
	packageName string
	models      []ModelDefinition
}

func NewModelsTemplater(packageName string, models []ModelDefinition) ModelsTemplater {
	return ModelsTemplater{
		packageName: packageName,
		models:      models,
	}
}

type ModelDefinition struct {
	Name     string
	JsonName string
	Fields   []PropertyDefinition
}

type PropertyDefinition struct {
	Name       string
	JsonName   string
	Type       string
	Required   bool
	Optional   bool
	Validation *PropertyValidationDefinition
}

type PropertyValidationDefinition struct {
	Type   PropertyValidationType
	Values *[]interface{}
}

type PropertyValidationType string

var (
	Range PropertyValidationType = "range"
)

func (g ModelsTemplater) Build() (*string, error) {
	// TODO: sort the fields then parse them in below

	models := make([]string, 0)
	for _, model := range g.models {
		log.Printf("[DEBUG] Generating %q..", model.Name)
		code, err := g.codeForModel(model)
		if err != nil {
			return nil, fmt.Errorf("generating model for %q: %+v", model.Name, err)
		}

		models = append(models, *code)
	}

	out := fmt.Sprintf(`package %s

import "github.com/hashicorp/go-multierror"

%s
`, g.packageName, strings.Join(models, "\n\n"))
	return &out, nil
}

func (g ModelsTemplater) codeForModel(definition ModelDefinition) (*string, error) {
	output := g.structForModel(definition)

	validation, err := definition.validationCode()
	if err != nil {
		return nil, err
	}

	if validation != nil {
		output += fmt.Sprintf("\n\n%s", *validation)
	}

	return &output, nil
}

func (g ModelsTemplater) structForModel(definition ModelDefinition) string {
	fields := make([]string, 0)

	for _, v := range definition.Fields {
		jsonTag := v.JsonName
		if v.Optional {
			jsonTag = fmt.Sprintf("%s,omitempty", v.JsonName)
		}

		format := "\t%s %s `json:\"%s\"`" // e.g. Foo string `json:"foo"`
		fields = append(fields, fmt.Sprintf(format, v.Name, v.Type, jsonTag))
	}

	return fmt.Sprintf(`type %s struct {
%s
}
`, definition.Name, strings.Join(fields, "\n"))
}

func (definition ModelDefinition) validationCode() (*string, error) {
	fields := make([]string, 0)
	for _, v := range definition.Fields {
		validationCode, err := v.validationCode()
		if err != nil {
			return nil, fmt.Errorf("generating validation for %q: %+v", v.Name, err)
		}
		if validationCode != nil {
			fields = append(fields, *validationCode)
		}
	}
	if len(fields) == 0 {
		return nil, nil
	}

	formattedFields := make([]string, 0)
	for _, field := range fields {
		formattedFields = append(formattedFields, fmt.Sprintf("\t%s", field))
	}

	out := fmt.Sprintf(`func (m %s) Validate() error {
  var result error

%s

  return result
}`, definition.Name, strings.Join(formattedFields, "\n\n"))
	return &out, nil
}

func (property PropertyDefinition) validationCode() (*string, error) {
	if !property.Required && !property.Optional && property.Validation == nil {
		return nil, nil
	}

	output := make([]string, 0)

	if property.Required && strings.EqualFold(property.Type, "String") {
		output = append(output, fmt.Sprintf(`if m.%[1]s == "" {
	result = multierror.Append(result, fmt.Errorf("%[1]s cannot be empty"))
}
`, property.Name))
	}

	// TODO: if there's nested objects, do they have validation functions?

	if property.Validation != nil {
		switch property.Validation.Type {
		case Range:
			{
				// TODO: implement me
				output = append(output, fmt.Sprintf(`// TODO: range validation`))
			}

		default:
			return nil, fmt.Errorf("unimplemented validation type %q!", property.Validation.Type)
		}
	}

	if len(output) == 0 {
		return nil, nil
	}

	result := strings.Join(output, "\n\n")
	return &result, nil
}
