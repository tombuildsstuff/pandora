package templates

import (
	"fmt"
	"strings"
)

type ConstantMetaData struct {
	// the type *has* to be a string, since these are treated as Enum's on the Azure side

	Values          []string
	CaseInsensitive bool
}

type ConstantsTemplater struct {
	packageName string
	constants   map[string]ConstantMetaData
}

func NewConstantsTemplater(packageName string, constants map[string]ConstantMetaData) ConstantsTemplater {
	return ConstantsTemplater{
		packageName: packageName,
		constants:   constants,
	}
}

func (t ConstantsTemplater) Build() (*string, error) {
	constants, err := t.constantsModels()
	if err != nil {
		return nil, fmt.Errorf("building constants: %+v", err)
	}

	template := fmt.Sprintf(`package %[1]s

%[2]s
`, t.packageName, *constants)
	return &template, nil
}

func (t ConstantsTemplater) constantsModels() (*string, error) {
	output := make([]string, 0)

	for name, metadata := range t.constants {
		values := make([]string, 0)
		for _, v := range metadata.Values {
			values = append(values, fmt.Sprintf("\t%s %s = %q", v, name, v))
		}
		code := fmt.Sprintf(`type %s string

var (
%s
)
`, name, strings.Join(values, "\n"))
		//if metadata.CaseInsensitive {
		//	// TODO: add a parse function that's case sensitive
		//	// NOTE: should we add a parse function anyway?
		//}

		output = append(output, code)
	}

	// not everything has constants so these are conditionally output
	if len(output) == 0 {
		return nil, nil
	}

	result := strings.Join(output, "\n\n")
	return &result, nil
}
