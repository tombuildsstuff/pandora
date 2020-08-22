package main

import (
	"log"
	"net/http"

	"github.com/tombuildsstuff/pandora/generator/models"
	"github.com/tombuildsstuff/pandora/generator/templates"
	"github.com/tombuildsstuff/pandora/generator/utils"
)

func main() {
	if err := clientMain(); err != nil {
		panic(err)
	}
}

func clientMain() error {
	methods := []models.OperationMetaData{
		{
			Name:                 "Delete",
			Method:               http.MethodDelete,
			LongRunningOperation: true,
			ExpectedStatusCodes: []int{
				200,
			},
		},
		{
			Name:                 "Get",
			Method:               http.MethodGet,
			LongRunningOperation: false,
			ExpectedStatusCodes: []int{
				200,
			},
		},
		{
			Name:                 "Create",
			Method:               http.MethodPut,
			LongRunningOperation: true,
			ExpectedStatusCodes: []int{
				200,
				201,
			},
		},
		{
			Name:                 "Update",
			Method:               http.MethodPatch,
			LongRunningOperation: true,
			ExpectedStatusCodes: []int{
				200,
				201,
			},
		},
	}
	//templater := templates.NewClientTemplater("example", "EventHubNamespace", "2018-01-01", nil, methods)
	templater := templates.NewModelsTemplater("example", "EventHubNamespace", methods)
	out, err := templater.Build()
	if err != nil {
		return err
	}
	//log.Printf("---\n%s\n---", *out)

	fmt, err := utils.GolangCodeFormatter{}.Format(*out)
	if err != nil {
		return err
	}

	log.Printf("---\n%s\n---", *fmt)
	return nil
}

func resourceIdMain() error {
	formatString := "/subscriptions/{subscriptionId}/resourceGroups/{name}/providers/Microsoft.EventHub/namespaces/{namespace}"
	segments := []string{
		"subscriptionId",
		"resourceGroup",
		"namespace",
	}
	template := templates.NewResourceIDTemplate("example", "EventHubNamespace", formatString, segments)
	output, err := template.Build()
	if err != nil {
		return err
	}

	fmt, err := utils.GolangCodeFormatter{}.Format(*output)
	if err != nil {
		return err
	}

	log.Printf("---\n%s\n---", *fmt)
	return nil
}
