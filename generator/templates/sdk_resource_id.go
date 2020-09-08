package templates

import (
	"fmt"
	"strings"

	"github.com/tombuildsstuff/pandora/generator/utils"
)

type ResourceIDTemplate struct {
	packageName            string
	typeName               string
	resourceIdFormatString string
	resourceIdSegments     []string
}

func NewResourceIDTemplate(packageName, typeName, resourceIdFormat string, resourceIdSegments []string) ResourceIDTemplate {
	return ResourceIDTemplate{
		packageName:            packageName,
		typeName:               strings.Title(typeName),
		resourceIdFormatString: resourceIdFormat,
		resourceIdSegments:     resourceIdSegments,
	}
}

func (t ResourceIDTemplate) Build() (*string, error) {
	structFields := t.fields("\t")
	arguments := t.arguments()
	constructorFields := t.constructorFieldAssignment("\t\t")
	fieldArguments := t.fieldArguments("id.")
	template := fmt.Sprintf(`package %[1]s

import "fmt"

type %[2]sId struct {
%[3]s
}

func New%[2]sId(%[4]s) %[2]sId {
	return %[2]sId{
%[5]s,
	}
}

func (id %[2]sId) ID() string {
	return fmt.Sprintf("%[6]s", %[7]s)
}`, t.packageName, t.typeName, structFields, arguments, constructorFields, t.resourceIdFormatString, fieldArguments)
	return &template, nil
}

// properties returns the segments formatted for used as the struct properties for an ID Parser
func (t ResourceIDTemplate) fields(indent string) string {
	output := make([]string, 0)
	for _, k := range t.resourceIdSegments {
		normalized := utils.NormalizePropertyName(k)
		formatted := fmt.Sprintf("%s%s\tstring", indent, normalized)
		output = append(output, formatted)
	}

	return strings.Join(output, "\n")
}

// arguments returns the segments formatted for used as the arguments for an ID Parser
func (t ResourceIDTemplate) arguments() string {
	args := make([]string, 0)
	for _, k := range t.resourceIdSegments {
		args = append(args, fmt.Sprintf("%s string", k))
	}

	return strings.Join(args, ", ")
}

func (t ResourceIDTemplate) constructorFieldAssignment(indent string) string {
	output := make([]string, 0)

	for _, k := range t.resourceIdSegments {
		normalized := utils.NormalizePropertyName(k)
		formatted := fmt.Sprintf("%s%s: %s", indent, normalized, k)
		output = append(output, formatted)
	}

	return strings.Join(output, ",\n")
}

func (t ResourceIDTemplate) fieldArguments(prefix string) string {
	output := make([]string, 0)

	for _, k := range t.resourceIdSegments {
		// e.g. "id.Name"
		normalized := utils.NormalizePropertyName(k)
		formatted := fmt.Sprintf("%s%s", prefix, normalized)
		output = append(output, formatted)
	}

	return strings.Join(output, ", ")
}
