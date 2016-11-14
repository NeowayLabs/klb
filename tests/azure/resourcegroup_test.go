package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
	"github.com/NeowayLabs/klb/tests/nash"
	"github.com/NeowayLabs/nash/sh"
)

func genResourceGroupName() string {
	return fmt.Sprintf("klb-resgroup-tests-%d", rand.Intn(1000))
}

func testResourceGroupCreation(t *testing.T) {
	session := azure.NewSession(t)

	shell := nash.Setup(t)

	resgroup := genResourceGroupName()

	shell.Setvar("ResourceGroup", sh.NewStrObj(resgroup))

	err := shell.Exec("TestResourceGroupCreation", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	resources := azure.NewResources(t, session)
	defer resources.Delete(t, resgroup)
	resources.AssertExists(t, resgroup)
}

func testResourceGroupDeletion(t *testing.T) {
	session := azure.NewSession(t)

	shell := nash.Setup(t)

	resgroup := genResourceGroupName()

	shell.Setvar("ResourceGroup", sh.NewStrObj(resgroup))

	err := shell.Exec("TestResourceGroupDeletion", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	resources := azure.NewResources(t, session)

	resources.AssertExists(t, resgroup)

	err = shell.Exec("TestResourceGroupDeletion2", `
            azure_group_delete($ResourceGroup)
        `)

	if err != nil {
		t.Error(err)
	}

	resources.AssertDeleted(t, resgroup)
}
