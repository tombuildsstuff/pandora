package main

import (
	"github.com/tombuildsstuff/pandora/generator/services"
)

func main() {
	if err := generateResourceManager("http://localhost:2020"); err != nil {
		panic(err)
	}
}

func generateResourceManager(endpoint string) error {
	outputDirectory := "/Users/tharvey/code/src/github.com/tombuildsstuff/pandora/resource-manager/out"
	client := services.NewResourceManagerService(endpoint)
	generator := services.NewResourceManagerGenerator(client, outputDirectory)
	return generator.Generate()
}
