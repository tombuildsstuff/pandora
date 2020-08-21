package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tombuildsstuff/pandora/sdk"
	"github.com/tombuildsstuff/pandora/target"
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
	name := fmt.Sprintf("tom-pandora-%d", time.Now().Unix())
	input := target.CreateResourceGroupInput{
		Location: "West Europe",
		Tags: map[string]string{
			"hello": "world",
		},
	}

	if err := input.Validate(); err != nil {
		return err
	}

	auth := sdk.NewClientSecretAuthorizer(clientId, clientSecret, tenantId)
	client := target.NewResourceGroupsClient(subscriptionId, auth)

	log.Printf("Creating %q", name)
	if err := client.Create(ctx, name, input); err != nil {
		return fmt.Errorf("creating: %+v", err)
	}

	id := target.NewResourceGroupID(name)
	log.Printf("Created %q", id.ID(subscriptionId))

	log.Printf("Retrieving %q", name)
	group, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving: %+v", err)
	}
	log.Printf("Exists in %q..", group.ResourceGroup.Location)
	log.Printf("Value for the Tag 'hello': %q..", group.ResourceGroup.Tags["hello"])

	log.Printf("Updating tags..")
	updateInput := target.UpdateResourceGroupInput{
		Tags: &map[string]string{
			"hello": "pandora",
		},
	}
	if err := client.Update(ctx, id, updateInput); err != nil {
		return fmt.Errorf("updating: %+v", err)
	}

	log.Printf("Retrieving %q", name)
	group, err = client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving: %+v", err)
	}
	log.Printf("Exists in %q..", group.ResourceGroup.Location)
	log.Printf("Value for the Tag 'hello': %q..", group.ResourceGroup.Tags["hello"])

	log.Printf("Deleting %q", name)
	poller, err := client.Delete(ctx, id)
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
