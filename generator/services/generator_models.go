package services

import "fmt"

type serviceDefinition struct {
	serviceName      string
	apiVersion       string
	resourceProvider *string
	packages         []packageDefinition
}

func (sd serviceDefinition) outputPath() string {
	return fmt.Sprintf("%s/%s", sd.serviceName, sd.apiVersion)
}

type packageDefinition struct {
	packageName string
	resourceId  ResourceIdDefinition
	models      map[string]ModelDefinition
	constants   map[string]ConstantDefinition
	operations  map[string]OperationDefinition
	operationId string
}

func (pd packageDefinition) outputPath(directory string) string {
	return fmt.Sprintf("%s/%s", directory, pd.packageName)
}
