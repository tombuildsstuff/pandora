package main

import (
	"fmt"
	"log"

	"github.com/tombuildsstuff/pandora/generator/services"
)

func main() {
	endpoint := "http://localhost:5000"
	outputDirectory := "/Users/tharvey/code/src/github.com/tombuildsstuff/pandora"

	log.Printf("[DEBUG] Generating Data Plane..")
	if err := generateDataPlane(endpoint, outputDirectory); err != nil {
		panic(err)
	}
	log.Printf("[DEBUG] Generated Data Plane")

	log.Printf("[DEBUG] Generating Resource Manager..")
	if err := generateResourceManager(endpoint, outputDirectory); err != nil {
		panic(err)
	}

	log.Printf("[DEBUG] Done")
}

func generateDataPlane(endpoint, outputDirectory string) error {
	return nil
	//out := fmt.Sprintf("%s/data-plane", outputDirectory)
	//client := services.NewResourceManagerService(endpoint)
	//generator := services.NewResourceManagerGenerator(client, out)
	//return generator.Generate()
}

func generateResourceManager(endpoint, outputDirectory string) error {
	out := fmt.Sprintf("%s/resource-manager", outputDirectory)
	client := services.NewResourceManagerService(endpoint)
	generator := services.NewResourceManagerGenerator(client, out)
	return generator.Generate()
}
