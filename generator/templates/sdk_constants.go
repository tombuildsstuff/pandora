package templates

import (
	"fmt"
	"sort"
	"strings"
)

type ConstantMetaData struct {
	// the type *has* to be a string, since these are treated as Enum's on the Azure side

	Values          map[string]string
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

	// if there are no constants
	if constants == nil {
		return nil, nil
	}

	template := fmt.Sprintf(`package %[1]s

%[2]s
`, t.packageName, *constants)
	return &template, nil
}

func (t ConstantsTemplater) constantsModels() (*string, error) {
	output := make([]string, 0)

	// sort the constant types
	constantNames := make([]string, 0)
	for k := range t.constants {
		constantNames = append(constantNames, k)
	}
	sort.Strings(constantNames)

	for _, constantName := range constantNames {
		metadata := t.constants[constantName]
		lines := make([]string, 0)

		// then sort the values for this constant
		keys := make([]string, 0)
		for k := range metadata.Values {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := metadata.Values[k]
			lines = append(lines, fmt.Sprintf("\t%s %s = %q", k, constantName, v))
		}
		code := fmt.Sprintf(`type %s string

var (
%s
)
`, constantName, strings.Join(lines, "\n"))
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
