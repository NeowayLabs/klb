package azure_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/nash"
	"github.com/NeowayLabs/nash/sh"
)

func genResourceGroupName() string {
	return fmt.Sprintf("klb-resgroup-tests-%d", rand.Intn(999999))
}

func testResourceGroupCreation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resgroup := genResourceGroupName()
	session := azure.NewSession(t)
	resources := azure.NewResourceGroup(ctx, t, session)
	defer resources.Delete(t, resgroup)

	shell := nash.Setup(t)
	shell.Setvar("ResourceGroup", sh.NewStrObj(resgroup))

	err := shell.Exec("TestResourceGroupCreation", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	resources.AssertExists(t, resgroup)
}

func testResourceGroupDeletion(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resgroup := genResourceGroupName()
	session := azure.NewSession(t)
	resources := azure.NewResourceGroup(ctx, t, session)

	shell := nash.Setup(t)
	shell.Setvar("ResourceGroup", sh.NewStrObj(resgroup))

	err := shell.Exec("TestResourceGroupDeletion", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	resources.AssertExists(t, resgroup)

	err = shell.Exec("TestResourceGroupDeletion2", `
            azure_group_delete($ResourceGroup)
        `)

	if err != nil {
		t.Error(err)
	}

	resources.AssertDeleted(t, resgroup)
}
