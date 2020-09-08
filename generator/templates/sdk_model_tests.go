package templates

import (
	"fmt"
	"strings"
)

type ModelsTestsTemplater struct {
	packageName string
	models      []ModelDefinition
}

func NewModelTestsTemplater(packageName string, models []ModelDefinition) ModelsTestsTemplater {
	return ModelsTestsTemplater{
		packageName: packageName,
		models:      models,
	}
}

func (t ModelsTestsTemplater) Build() (*string, error) {
	modelsRequiringValidation := make([]string, 0)
	for _, v := range t.models {
		code, err := v.validationCode()
		if err != nil {
			return nil, fmt.Errorf("generating validation code for %s: %+v", v.Name, err)
		}

		if code != nil {
			modelsRequiringValidation = append(modelsRequiringValidation, v.Name)
		}
	}

	// don't bother writing out an empty file
	if len(modelsRequiringValidation) == 0 {
		return nil, nil
	}

	lines := make([]string, 0)
	for _, modelName := range modelsRequiringValidation {
		lines = append(lines, fmt.Sprintf("var _ sdk.ModelWithValidation = %s{}", modelName))
	}

	format := fmt.Sprintf(`package %s

import "github.com/tombuildsstuff/pandora/sdk"

%s

// TODO: unit tests for the API methods based on sample responses
`, t.packageName, strings.Join(lines, "\n\n"))
	return &format, nil
}
