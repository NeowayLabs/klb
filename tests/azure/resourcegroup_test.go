package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
	"github.com/NeowayLabs/klb/tests/nash"
	"github.com/NeowayLabs/nash/sh"
)

var (
	ResourceGroupName = fmt.Sprintf("klb-resgroup-test-%d",
		rand.Intn(1000))
)

func TestHandleResourceGroupLifeCycle(t *testing.T) {
	session := azure.NewSession(t)

	shell := nash.Setup(t)
	shell.Setvar("ResourceGroup", sh.NewStrObj(ResourceGroupName))

	err := shell.Exec("TestHandleResourceGroupLifeCycle", `
             import ../../azure/all
             azure_group_create($ResourceGroup, "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	resources := azure.NewResources(t, session)
	defer resources.Delete(t, ResourceGroupName)
	resources.Check(t, ResourceGroupName)
}
