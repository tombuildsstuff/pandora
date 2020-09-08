package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tombuildsstuff/pandora/resource-manager/eventhubs/2018-01-01-preview/namespaces"
	"github.com/tombuildsstuff/pandora/resource-manager/resources/2018-05-01/resourcegroups"
	"github.com/tombuildsstuff/pandora/sdk"
)

func main() {
	if err := run(context.TODO()); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	tenantId := os.Getenv("ARM_TENANT_ID")
	rInt := time.Now().Unix()
	name := fmt.Sprintf("tom-pandora-%d", rInt)
	input := resourcegroups.CreateResourceGroupInput{
		Location: "West Europe",
		Tags: map[string]string{
			"hello": "world",
		},
	}

	if err := input.Validate(); err != nil {
		return err
	}

	auth := sdk.NewClientSecretAuthorizer(clientId, clientSecret, tenantId)
	groupsClient := resourcegroups.NewClient(subscriptionId, auth)
	namespacesClient := namespaces.NewNamespacesClient(subscriptionId, auth)

	id := resourcegroups.NewResourceGroupID(name)

	log.Printf("Creating %q", name)
	if err := groupsClient.Create(ctx, id, input); err != nil {
		return fmt.Errorf("creating: %+v", err)
	}

	log.Printf("Created %q", id.ID(subscriptionId))

	log.Printf("Retrieving %q", name)
	group, err := groupsClient.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving: %+v", err)
	}
	log.Printf("Exists in %q..", group.ResourceGroup.Location)
	log.Printf("Value for the Tag 'hello': %q..", group.ResourceGroup.Tags["hello"])

	log.Printf("Updating tags..")
	updateInput := resourcegroups.UpdateResourceGroupInput{
		Tags: &map[string]string{
			"hello": "pandora",
		},
	}
	if err := groupsClient.Update(ctx, id, updateInput); err != nil {
		return fmt.Errorf("updating: %+v", err)
	}

	log.Printf("Retrieving %q", name)
	group, err = groupsClient.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving: %+v", err)
	}
	log.Printf("Exists in %q..", group.ResourceGroup.Location)
	log.Printf("Value for the Tag 'hello': %q..", group.ResourceGroup.Tags["hello"])

	// add a nested item
	namespaceName := fmt.Sprintf("tomdev%d", rInt)
	namespaceId := namespaces.NewNamespacesId(id.Name, namespaceName)
	ptr := false
	createNamespaceInput := namespaces.CreateNamespaceInput{
		Location: input.Location,
		Sku: namespaces.Sku{
			Name: namespaces.Basic,
			Tier: namespaces.Basic,
		},
		Properties: namespaces.CreateNamespaceProperties{
			IsAutoInflateEnabled: &ptr,
			ZoneRedundant:        &ptr,
		},
		Tags: &map[string]string{},
	}
	log.Printf("Adding a EventHub Namespace %q", namespaceName)
	poller, err := namespacesClient.Create(ctx, namespaceId, createNamespaceInput)
	if err != nil {
		return fmt.Errorf("creating namespace: %+v", err)
	}
	log.Printf("Waiting for creation of %q", namespaceName)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation: %+v", err)
	}

	log.Printf("Retrieving Namespace %q", namespaceName)
	namespace, err := namespacesClient.Get(ctx, namespaceId)
	if err != nil {
		return fmt.Errorf("retrieving namespace: %+v", err)
	}

	log.Printf("ServiceBus Endpoint is at %q", namespace.GetNamespace.Properties.ServiceBusEndpoint)
	time.Sleep(10 * time.Second)

	log.Printf("Deleting EH namespace %q", namespaceName)
	poller, err = namespacesClient.Delete(ctx, namespaceId)
	if err != nil {
		return fmt.Errorf("deleting namespace: %+v", err)
	}
	log.Printf("Waiting for deletion of %q", namespaceName)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion: %+v", err)
	}
	log.Printf("Deleted %q", namespaceName)

	log.Printf("Deleting %q", name)
	poller, err = groupsClient.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting: %+v", err)
	}
	log.Printf("Waiting for deletion of %q", name)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion: %+v", err)
	}
	log.Printf("Done")

	return nil
}
