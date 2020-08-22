package templates

import "testing"

func TestEventHubNamespaceResourceID(t *testing.T) {
	expected := `package example

import "fmt"

type EventHubNamespaceID struct {
	ResourceGroup	string
	Namespace	string
}

func NewEventHubNamespaceID(resourceGroup string, namespace string) EventHubNamespaceID {
	return EventHubNamespaceID{
		ResourceGroup: resourceGroup,
		Namespace: namespace,
	}
}

func (id EventHubNamespaceID) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/{subscriptionId}/resourceGroups/{name}/providers/Microsoft.EventHub/namespaces/{namespace}", subscriptionId, id.ResourceGroup, id.Namespace)
}`

	formatString := "/subscriptions/{subscriptionId}/resourceGroups/{name}/providers/Microsoft.EventHub/namespaces/{namespace}"
	segments := []string{
		"resourceGroup",
		"namespace",
	}
	template := NewResourceIDTemplate("example", "EventHubNamespace", formatString, segments)
	actual, err := template.Build()
	if err != nil {
		t.Fatal(err)
	}

	if *actual != expected {
		t.Fatalf("Expected `%s` but got `%s`", expected, *actual)
	}
}

func TestResourceGroupResourceID(t *testing.T) {
	expected := `package dora

import "fmt"

type ResourceGroupID struct {
	Name	string
}

func NewResourceGroupID(name string) ResourceGroupID {
	return ResourceGroupID{
		Name: name,
	}
}

func (id ResourceGroupID) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/{subscriptionId}/resourceGroups/{name}", subscriptionId, id.Name)
}`

	formatString := "/subscriptions/{subscriptionId}/resourceGroups/{name}"
	segments := []string{
		"name",
	}
	template := NewResourceIDTemplate("dora", "ResourceGroup", formatString, segments)
	actual, err := template.Build()
	if err != nil {
		t.Fatal(err)
	}

	if *actual != expected {
		t.Fatalf("Expected `%s` but got `%s`", expected, *actual)
	}
}
