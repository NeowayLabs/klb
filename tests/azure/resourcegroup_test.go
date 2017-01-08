package azure_test

import (
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

	session := azure.NewSession(t)
	resources := azure.NewResourceGroup(t, session)
	defer resources.Delete(t, resgroup)
	resources.AssertExists(t, resgroup)
}

func testResourceGroupDeletion(t *testing.T) {

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

	session := azure.NewSession(t)
	resources := azure.NewResourceGroup(t, session)

	resources.AssertExists(t, resgroup)

	err = shell.Exec("TestResourceGroupDeletion2", `
            azure_group_delete($ResourceGroup)
        `)

	if err != nil {
		t.Error(err)
	}

	resources.AssertDeleted(t, resgroup)
}
