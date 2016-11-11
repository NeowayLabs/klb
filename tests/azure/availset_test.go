package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
	"github.com/NeowayLabs/klb/tests/nash"

	"github.com/NeowayLabs/nash/sh"
)

func genAvailSetName() string {
	return fmt.Sprintf("klb-availset-tests-%d", rand.Intn(1000))
}

func TestAvailSetCreation(t *testing.T) {
	session := azure.NewSession(t)

	shell := nash.Setup(t)

	resgroup := genResourceGroupName()
	availset := genAvailSetName()

	shell.Setvar("ResourceGroup", sh.NewStrObj(resgroup))
	shell.Setvar("AvailSet", sh.NewStrObj(availset))

	err := shell.Exec("TestAvailSetCreation", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
             azure_availset_create($AvailSet, $ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	availSets := azure.NewAvailSet(t, session)
	resources := azure.NewResources(t, session)

	defer resources.Delete(t, resgroup)
	availSets.Check(t, availset, resgroup)
}

func TestAvailSetDeletion(t *testing.T) {
	session := azure.NewSession(t)

	shell := nash.Setup(t)

	resgroup := genResourceGroupName()
	availset := genAvailSetName()

	shell.Setvar("ResourceGroup", sh.NewStrObj(resgroup))
	shell.Setvar("AvailSet", sh.NewStrObj(availset))

	err := shell.Exec("TestAvailSetCreation", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
             azure_availset_create($AvailSet, $ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	availSets := azure.NewAvailSet(t, session)
	resources := azure.NewResources(t, session)

	defer resources.Delete(t, resgroup)

	availSets.Check(t, availset, resgroup)

	err = shell.Exec("TestAvailSetCreation", `
             azure_availset_delete($AvailSet, $ResourceGroup)
        `)

	availSets.CheckDelete(t, availset, resgroup)
}
