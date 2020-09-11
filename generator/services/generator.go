package services

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tombuildsstuff/pandora/generator/utils"
)

type SdkGenerator struct {
	client          PandoraApi
	outputDirectory string
}

func NewSdkGenerator(client PandoraApi, outputDirectory string) SdkGenerator {
	return SdkGenerator{
		outputDirectory: outputDirectory,
		client:          client,
	}
}

// TODO: use the Operation Name in the file name so that we can have
// multiple operations within a single package
// e.g. `namespace_client.go`, `namespace_models.go` and `namespace`

func (g SdkGenerator) Generate() error {
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
				serviceDef:      serviceDefinition,
				packageDef:      packageDefinition,
				outputDirectory: outputDirectory,
			}

			log.Printf("[DEBUG] Generating files..")
			if err := packageGenerator.generate(); err != nil {
				return fmt.Errorf("generating package %q (Service %q / API Version %q): %+v", packageDefinition.packageName, serviceDefinition.serviceName, serviceDefinition.apiVersion, err)
			}
		}
	}
	return nil
}

func (g SdkGenerator) buildServiceDefinitions() (*[]serviceDefinition, error) {
	log.Printf("[DEBUG] Retrieving Supported API's..")
	apis, err := g.client.Apis()
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
		availableApiVersions, err := g.client.VersionsForApi(serviceDetails)
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
			apiOperations, err := g.client.APIOperationsForApiVersion(*operationMetaData)
			if err != nil {
				return nil, fmt.Errorf("retrieving api operations for %q: %+v", operationName, err)
			}

			log.Printf("[DEBUG] Retrieving Schemas for %q..", operationName)
			schemas, err := g.client.SchemasForApiVersion(*operationMetaData)
			if err != nil {
				return nil, fmt.Errorf("retrieving schema for %q: %+v", operationName, err)
			}

			log.Printf("[DEBUG] Appending Package %q..", operationName)
			packages = append(packages, packageDefinition{
				packageName: operationName,
				typeName:    operationMetaData.Name,
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
			resourceProvider: availableApiVersions.ResourceProvider,
			packages:         packages,
		})
	}

	return &services, nil
}

type apiVersionDetails struct {
	apiVersion string
	details    VersionDetails
}

func (g SdkGenerator) determineApiVersion(versions *ApiVersionsResponse) (*apiVersionDetails, error) {
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
	serviceDef      serviceDefinition
	packageDef      packageDefinition
	outputDirectory string
	tfOutputPath    string
}

type generator interface {
	directory(workingDirectory, packageName string) string
	generate(serviceDef serviceDefinition, packageDef packageDefinition, outputPath string) error
	name() string
}

func (p packageGenerator) generate() error {
	generators := []generator{
		sdkGenerator{},
	}

	for _, generator := range generators {
		log.Printf("[DEBUG] Starting generating %s..", generator.name())
		directory := generator.directory(p.outputDirectory, p.packageDef.packageName)
		os.MkdirAll(directory, os.ModePerm)

		if err := generator.generate(p.serviceDef, p.packageDef, directory); err != nil {
			return err
		}
		log.Printf("[DEBUG] Finished generating for %s..", generator.name())
	}

	return nil
}

func goFmtAndWriteToFile(filePath, fileContents string) error {
	fmt, err := utils.GolangCodeFormatter{}.Format(fileContents)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(*fmt), 0644)
}
