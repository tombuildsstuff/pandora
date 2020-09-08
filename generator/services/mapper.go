package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tombuildsstuff/pandora/generator/templates"
)

func (pd packageDefinition) buildModelDefinitions() (*[]templates.ModelDefinition, error) {
	models := make([]templates.ModelDefinition, 0)

	for modelName, properties := range pd.models {
		fields := make([]templates.PropertyDefinition, 0)
		for propertyName, property := range properties {
			propType, err := pd.determineTypeForProperty(property)
			if err != nil {
				return nil, fmt.Errorf("determining type for property %q: %+v", propertyName, err)
			}

			fieldName := strings.Title(propertyName)
			result := templates.PropertyDefinition{
				Name:     fieldName,
				JsonName: propertyName,
				Type:     *propType,
				Required: property.Required,
				Optional: property.Optional,
			}

			if property.Optional {
				// add a pointer
				result.Type = fmt.Sprintf("*%s", result.Type)
			}

			if property.Validation != nil {
				result.Validation = &templates.PropertyValidationDefinition{
					Type:   templates.PropertyValidationType(property.Validation.Type),
					Values: property.Validation.Values,
				}
			}

			fields = append(fields, result)
		}

		fields = sortFields(fields)

		models = append(models, templates.ModelDefinition{
			Fields:   fields,
			Name:     strings.Title(modelName),
			JsonName: modelName,
		})
	}

	models = sortModels(models)

	return &models, nil
}

func (pd packageDefinition) determineTypeForProperty(def PropertyDefinition) (*string, error) {
	if def.Type == Constant {
		if def.ConstantReference == nil {
			return nil, fmt.Errorf("constant without a reference")
		}

		ref, err := parseReference(*def.ConstantReference)
		if err != nil {
			return nil, fmt.Errorf("parsing reference %q", *def.ConstantReference)
		}

		return &ref.name, nil
	}

	// TODO: support for Lists and Sets of Objects
	if def.Type == Object {
		if def.ModelReference == nil {
			return nil, fmt.Errorf("model without a reference")
		}

		ref, err := parseReference(*def.ModelReference)
		if err != nil {
			return nil, fmt.Errorf("parsing reference %q", *def.ModelReference)
		}

		return &ref.name, nil
	}

	out := string(def.Type)
	switch def.Type {
	case Boolean:
		out = "bool"
		break

	case Location:
		out = "string"
		break

	case Integer:
		out = "int64"
		break

	case Tags:
		out = "map[string]string"
		break
	}

	return &out, nil
}

func sortFields(input []templates.PropertyDefinition) []templates.PropertyDefinition {
	keys := make([]string, 0)
	vals := make(map[string]templates.PropertyDefinition)
	for _, v := range input {
		keys = append(keys, v.Name)
		vals[v.Name] = v
	}
	sort.Strings(keys)

	out := make([]templates.PropertyDefinition, 0)
	for _, key := range keys {
		out = append(out, vals[key])
	}

	return out
}

func sortModels(input []templates.ModelDefinition) []templates.ModelDefinition {
	keys := make([]string, 0)
	vals := make(map[string]templates.ModelDefinition)
	for _, v := range input {
		keys = append(keys, v.Name)
		vals[v.Name] = v
	}
	sort.Strings(keys)

	out := make([]templates.ModelDefinition, 0)
	for _, key := range keys {
		out = append(out, vals[key])
	}

	return out
}
