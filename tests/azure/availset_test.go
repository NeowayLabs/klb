package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
	"github.com/NeowayLabs/klb/tests/lib/nash"

	"github.com/NeowayLabs/nash/sh"
)

func genAvailSetName() string {
	return fmt.Sprintf("klb-availset-tests-%d", rand.Intn(1000))
}

func testAvailSetCreation(t *testing.T, f fixture.Fixture) {
	availset := genAvailSetName()
	nash.Run(
		t,
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availset,
		f.Location,
	)
	availSets := azure.NewAvailSet(t, f.Session)
	availSets.AssertExists(t, availset, f.ResGroupName)
}

func testAvailSetDeletion(t *testing.T) {

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

	session := azure.NewSession(t)
	availSets := azure.NewAvailSet(t, session)
	resources := azure.NewResources(t, session)

	defer resources.Delete(t, resgroup)

	availSets.AssertExists(t, availset, resgroup)

	err = shell.Exec("TestAvailSetCreation", `
             azure_availset_delete($AvailSet, $ResourceGroup)
        `)

	availSets.AssertDeleted(t, availset, resgroup)
}
