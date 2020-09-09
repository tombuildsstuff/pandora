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
				JsonName: property.JsonName,
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
	if def.Type == List {
		if def.ListElementType == nil {
			return nil, fmt.Errorf("`listElementType` cannot be nil when `propertyType` is a `List`")
		}

		innerType, err := pd.determineTypeForPropertyInternal(*def.ListElementType, def.ConstantReference, def.ModelReference)
		if err != nil {
			return nil, fmt.Errorf("retrieving `listElementType`: %+v", err)
		}

		out := fmt.Sprintf("[]%s", *innerType)
		return &out, nil
	}

	return pd.determineTypeForPropertyInternal(def.Type, def.ConstantReference, def.ModelReference)
}

func (pd packageDefinition) determineTypeForPropertyInternal(propertyType PropertyType, constantReference, modelReference *string) (*string, error) {
	if propertyType == Constant {
		return pd.determineTypeForConstant(constantReference)
	}

	if propertyType == Object {
		return pd.determineTypeForObject(modelReference)
	}

	out := string(propertyType)
	switch propertyType {
	case Boolean:
		out = "bool"
		break

	case Location:
		out = "string"
		break

	case Integer:
		out = "int64"
		break

	case String:
		out = "string"
		break

	case Tags:
		out = "map[string]string"
		break

		// TODO: other types
	}

	return &out, nil
}

func (pd packageDefinition) determineTypeForConstant(constantReference *string) (*string, error) {
	if constantReference == nil {
		return nil, fmt.Errorf("constant without a reference")
	}

	ref, err := parseReference(*constantReference)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %q", *constantReference)
	}

	return &ref.name, nil
}

func (pd packageDefinition) determineTypeForObject(modelReference *string) (*string, error) {
	if modelReference == nil {
		return nil, fmt.Errorf("model without a reference")
	}

	ref, err := parseReference(*modelReference)
	if err != nil {
		return nil, fmt.Errorf("parsing reference %q", *modelReference)
	}

	return &ref.name, nil
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
