package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tombuildsstuff/pandora/generator/templates"
	"github.com/tombuildsstuff/pandora/generator/utils"
)

type ResourceManagerGenerator struct {
	client          ResourceManagerService
	outputDirectory string
}

func NewResourceManagerGenerator(client ResourceManagerService, outputDirectory string) ResourceManagerGenerator {
	return ResourceManagerGenerator{
		outputDirectory: outputDirectory,
		client:          client,
	}
}

func (g ResourceManagerGenerator) Generate() error {
	log.Printf("[DEBUG] Populating API information into Service Definitions..")
	services, err := g.buildServiceDefinitions()
	if err != nil {
		return fmt.Errorf("building service definitions: %+v", err)
	}

	log.Printf("[DEBUG] Removing any existing output directory %q..", g.outputDirectory)
	os.RemoveAll(g.outputDirectory)
	for _, serviceDefinition := range *services {
		log.Printf("[DEBUG] Generating Service %q (API Version %q)..", serviceDefinition.serviceName, serviceDefinition.apiVersion)

		outputDirectory := fmt.Sprintf("%s/%s", g.outputDirectory, serviceDefinition.outputPath())
		log.Printf("[DEBUG] Creating %q..", outputDirectory)
		os.MkdirAll(outputDirectory, os.ModePerm)

		for _, packageDefinition := range serviceDefinition.packages {
			log.Printf("[DEBUG] Generating Package %q..", packageDefinition.packageName)
			packageOutputPath := packageDefinition.outputPath(outputDirectory)
			log.Printf("[DEBUG] Creating %q..", packageOutputPath)
			os.MkdirAll(packageOutputPath, os.ModePerm)

			packageGenerator := packageGenerator{
				serviceDef: serviceDefinition,
				packageDef: packageDefinition,
				outputPath: packageOutputPath,
			}

			log.Printf("[DEBUG] Generating files..")
			if err := packageGenerator.generate(); err != nil {
				return fmt.Errorf("generating package %q (Service %q / API Version %q): %+v", packageDefinition.packageName, serviceDefinition.serviceName, serviceDefinition.apiVersion, err)
			}
		}
	}
	return nil
}

func (g ResourceManagerGenerator) buildServiceDefinitions() (*[]serviceDefinition, error) {
	log.Printf("[DEBUG] Retrieving Supported API's..")
	apis, err := g.client.SupportedApis()
	if err != nil {
		return nil, fmt.Errorf("retrieving supported API's: %+v", err)
	}

	services := make([]serviceDefinition, 0)
	for serviceName, serviceDetails := range apis.Apis {
		log.Printf("[DEBUG] Current Service: %q", serviceName)
		if !serviceDetails.Generate {
			log.Printf("[DEBUG] Generation Disabled - skipping")
			continue
		}

		log.Printf("[DEBUG] Retrieving available API versions..")
		availableApiVersions, err := g.client.SupportedVersionsForApi(serviceDetails)
		if err != nil {
			return nil, fmt.Errorf("retrieving available API versions: %+v", err)
		}

		log.Printf("[DEBUG] Determining API version..")
		apiVersionDetails, err := g.determineApiVersion(availableApiVersions)
		if err != nil {
			return nil, fmt.Errorf("determining API version: %+v", err)
		}

		log.Printf("[DEBUG] Retrieving Operations for API Version %q..", apiVersionDetails.apiVersion)
		supportedOperations, err := g.client.OperationsForApiVersion(apiVersionDetails.details)
		if err != nil {
			return nil, fmt.Errorf("retrieving operations for API version %q: %+v", apiVersionDetails.apiVersion, err)
		}

		log.Printf("[DEBUG] Retrieved %d types.", len(supportedOperations.Types))
		packages := make([]packageDefinition, 0)
		for operationName, operationDetails := range supportedOperations.Types {
			log.Printf("[DEBUG] Current Operation %q..", operationName)

			log.Printf("[DEBUG] Retrieving MetaData for %q..", operationName)
			operationMetaData, err := g.client.MetaDataForOperation(operationDetails)
			if err != nil {
				return nil, fmt.Errorf("retrieving metadata for %q: %+v", operationName, err)
			}

			log.Printf("[DEBUG] Retrieving API Operations for %q..", operationName)
			apiOperations, err := g.client.OperationsForType(*operationMetaData)
			if err != nil {
				return nil, fmt.Errorf("retrieving api operations for %q: %+v", operationName, err)
			}

			log.Printf("[DEBUG] Retrieving Schemas for %q..", operationName)
			schemas, err := g.client.SchemaForType(*operationMetaData)
			if err != nil {
				return nil, fmt.Errorf("retrieving schema for %q: %+v", operationName, err)
			}

			log.Printf("[DEBUG] Appending Package %q..", operationName)
			packages = append(packages, packageDefinition{
				packageName: operationName,
				resourceId:  operationDetails.ResourceId,
				models:      schemas.Models,
				constants:   schemas.Constants,
				operations:  apiOperations.Operations,
				operationId: operationDetails.Uri, // used in models
			})
		}

		log.Printf("[DEBUG] Appending Service..")
		services = append(services, serviceDefinition{
			serviceName:      serviceName,
			apiVersion:       apiVersionDetails.apiVersion,
			resourceProvider: &availableApiVersions.ResourceProvider,
			packages:         packages,
		})
	}

	return &services, nil
}

type apiVersionDetails struct {
	apiVersion string
	details    VersionDetails
}

func (g ResourceManagerGenerator) determineApiVersion(versions *SupportedVersionsResponse) (*apiVersionDetails, error) {
	if versions == nil {
		return nil, fmt.Errorf("no versions were available!")
	}

	for versionNumber, versionDetails := range versions.Versions {
		if versionDetails.Generate {
			return &apiVersionDetails{
				apiVersion: versionNumber,
				details:    versionDetails,
			}, nil
		}
	}

	return nil, fmt.Errorf("no version was marked as to generate")
}

type packageGenerator struct {
	serviceDef serviceDefinition
	packageDef packageDefinition
	outputPath string
}

func (p packageGenerator) generate() error {
	log.Printf("[DEBUG] Generating Client for Operations..")
	if err := p.generateOperations(fmt.Sprintf("%s/client.go", p.outputPath)); err != nil {
		return fmt.Errorf("generating client for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating Constants for Operations..")
	if err := p.generateConstants(fmt.Sprintf("%s/constants.go", p.outputPath)); err != nil {
		return fmt.Errorf("generating constants for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating ID Parsers for Operations..")
	if err := p.generateIDParsers(fmt.Sprintf("%s/id_parsers.go", p.outputPath)); err != nil {
		return fmt.Errorf("generating ID parsers for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating Models for Operations..")
	if err := p.generateModels(fmt.Sprintf("%s/models.go", p.outputPath)); err != nil {
		return fmt.Errorf("generating models for operations: %+v", err)
	}

	log.Printf("[DEBUG] Generating Model Tests for Operations..")
	if err := p.generateModelTests(fmt.Sprintf("%s/models_test.go", p.outputPath)); err != nil {
		return fmt.Errorf("generating model tests for operations: %+v", err)
	}

	return nil
}

func (p packageGenerator) generateConstants(filePath string) error {
	packageName := p.packageDef.packageName
	constants := make(map[string]templates.ConstantMetaData, 0)
	for k, v := range p.packageDef.constants {
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

func (p packageGenerator) generateIDParsers(filePath string) error {
	packageName := p.packageDef.packageName
	resourceId := p.packageDef.resourceId

	templater := templates.NewResourceIDTemplate(packageName, packageName, resourceId.Format, resourceId.Segments)
	out, err := templater.Build()
	if err != nil {
		return fmt.Errorf("building template: %+v", err)
	}

	return goFmtAndWriteToFile(filePath, *out)
}

func (p packageGenerator) generateModels(filePath string) error {
	models, err := p.packageDef.buildModelDefinitions()
	if err != nil {
		return fmt.Errorf("building model definitions: %+v", err)
	}

	templater := templates.NewModelsTemplater(p.packageDef.packageName, *models)
	out, err := templater.Build()
	if err != nil {
		return fmt.Errorf("generating models: %+v", err)
	}

	return goFmtAndWriteToFile(filePath, *out)
}

func (p packageGenerator) generateModelTests(filePath string) error {
	models, err := p.packageDef.buildModelDefinitions()
	if err != nil {
		return fmt.Errorf("building model definitions: %+v", err)
	}

	templater := templates.NewModelTestsTemplater(p.packageDef.packageName, *models)
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

func (p packageGenerator) generateOperations(filePath string) error {
	operations := make([]templates.ClientOperation, 0)
	for name, operation := range p.packageDef.operations {
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

	apiVersion := p.serviceDef.apiVersion
	packageName := p.packageDef.packageName
	resourceProvider := p.serviceDef.resourceProvider

	// this is resource manager so resourceProvider is guaranteed
	if resourceProvider == nil {
		return fmt.Errorf("resourceProvider was nil for a Resource Manager Operation")
	}

	templater := templates.NewResourceManagerClientTemplater(packageName, packageName, apiVersion, *resourceProvider, operations)
	output, err := templater.Build()
	if err != nil {
		return fmt.Errorf("templating: %+v", err)
	}

	return goFmtAndWriteToFile(filePath, *output)
}

func goFmtAndWriteToFile(filePath, fileContents string) error {
	fmt, err := utils.GolangCodeFormatter{}.Format(fileContents)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(*fmt), 0644)
}
