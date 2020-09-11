package services

import (
	"fmt"
	"log"

	"github.com/tombuildsstuff/pandora/generator/templates"
)

type sdkGenerator struct{}

func (p sdkGenerator) directory(workingDirectory, packageName string) string {
	return fmt.Sprintf("%s/%s", workingDirectory, packageName)
}

func (p sdkGenerator) generate(serviceDef serviceDefinition, packageDef packageDefinition, outputPath string) error {
	log.Printf("[DEBUG] Generating Client for Operations..")
	clientsFile := fmt.Sprintf("%s/client.go", outputPath)
	if err := p.generateOperations(serviceDef, packageDef, clientsFile); err != nil {
		return fmt.Errorf("generating client for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating Constants for Operations..")
	constantsFile := fmt.Sprintf("%s/constants.go", outputPath)
	if err := p.generateConstants(packageDef, constantsFile); err != nil {
		return fmt.Errorf("generating constants for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating ID Parsers for Operations..")
	idParsersFile := fmt.Sprintf("%s/id_parsers.go", outputPath)
	if err := p.generateIDParsers(packageDef, idParsersFile); err != nil {
		return fmt.Errorf("generating ID parsers for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating Models for Operations..")
	modelsFile := fmt.Sprintf("%s/models.go", outputPath)
	if err := p.generateModels(packageDef, modelsFile); err != nil {
		return fmt.Errorf("generating models for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating Model Tests for Operations..")
	modelTestsFile := fmt.Sprintf("%s/models_test.go", outputPath)
	if err := p.generateModelTests(packageDef, modelTestsFile); err != nil {
		return fmt.Errorf("generating model tests for operations: %+v", err)
	}

	return nil
}

func (p sdkGenerator) name() string {
	return "SDK"
}

func (p sdkGenerator) generateConstants(packageDef packageDefinition, filePath string) error {
	packageName := packageDef.packageName
	constants := make(map[string]templates.ConstantMetaData, 0)
	for k, v := range packageDef.constants {
		constants[k] = templates.ConstantMetaData{
			Values:          v.Values,
			CaseInsensitive: v.CaseInsensitive,
		}
	}

	templater := templates.NewConstantsTemplater(packageName, constants)
	out, err := templater.Build()
	if err != nil {
		return fmt.Errorf("building template: %+v", err)
	}

	// not everything has constants so this file is conditionally output
	if out == nil {
		return nil
	}

	return goFmtAndWriteToFile(filePath, *out)
}

func (p sdkGenerator) generateIDParsers(packageDef packageDefinition, filePath string) error {
	packageName := packageDef.packageName
	typeName := packageDef.typeName
	resourceId := packageDef.resourceId

	templater := templates.NewResourceIDTemplate(packageName, typeName, resourceId.Format, resourceId.Segments)
	out, err := templater.Build()
	if err != nil {
		return fmt.Errorf("building template: %+v", err)
	}

	return goFmtAndWriteToFile(filePath, *out)
}

func (p sdkGenerator) generateModels(packageDef packageDefinition, filePath string) error {
	models, err := packageDef.buildModelDefinitions()
	if err != nil {
		return fmt.Errorf("building model definitions: %+v", err)
	}

	templater := templates.NewModelsTemplater(packageDef.packageName, *models)
	out, err := templater.Build()
	if err != nil {
		return fmt.Errorf("generating models: %+v", err)
	}

	return goFmtAndWriteToFile(filePath, *out)
}

func (p sdkGenerator) generateModelTests(packageDef packageDefinition, filePath string) error {
	models, err := packageDef.buildModelDefinitions()
	if err != nil {
		return fmt.Errorf("building model definitions: %+v", err)
	}

	templater := templates.NewModelTestsTemplater(packageDef.packageName, *models)
	out, err := templater.Build()
	if err != nil {
		return fmt.Errorf("generating models tests: %+v", err)
	}

	// model tests are only generated when there's a validation function
	// to ensure the model complies with the interface, so check before
	// writing an empty file out
	if out == nil {
		return nil
	}

	return goFmtAndWriteToFile(filePath, *out)
}

func (p sdkGenerator) generateOperations(serviceDef serviceDefinition, packageDef packageDefinition, filePath string) error {
	operations := make([]templates.ClientOperation, 0)
	for name, operation := range packageDef.operations {
		clientOperation := templates.ClientOperation{
			Name:                 name,
			Method:               operation.Method,
			LongRunningOperation: operation.LongRunning,
			ExpectedStatusCodes:  operation.ExpectedStatusCodes,
		}

		if operation.RequestObject != nil {
			ref, err := parseReference(*operation.RequestObject)
			if err != nil {
				return fmt.Errorf("parsing reference %q: %+v", *operation.RequestObject, err)
			}

			clientOperation.RequestObjectName = &ref.name
		}

		if operation.ResponseObject != nil {
			ref, err := parseReference(*operation.ResponseObject)
			if err != nil {
				return fmt.Errorf("parsing reference %q: %+v", *operation.ResponseObject, err)
			}

			clientOperation.ResponseObjectName = &ref.name
		}

		operations = append(operations, clientOperation)
	}

	apiVersion := serviceDef.apiVersion
	packageName := packageDef.packageName
	typeName := packageDef.typeName
	resourceProvider := serviceDef.resourceProvider

	var output *string
	var err error
	if resourceProvider != nil {
		log.Printf("[DEBUG] Generating Resource Manager Client..")
		templater := templates.NewResourceManagerClientTemplater(packageName, typeName, apiVersion, *resourceProvider, operations)
		output, err = templater.Build()
		if err != nil {
			return fmt.Errorf("templating: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Generating Data Plane Client..")
		templater := templates.NewDataPlaneClientTemplater(packageName, typeName, apiVersion, operations)
		output, err = templater.Build()
		if err != nil {
			return fmt.Errorf("templating: %+v", err)
		}
	}

	return goFmtAndWriteToFile(filePath, *output)
}
