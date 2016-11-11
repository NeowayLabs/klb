package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
	"github.com/NeowayLabs/klb/tests/nash"
)

var (
	ResourceGroupName = fmt.Sprintf("klb-resgroup-test-%d",
		rand.Intn(1000))
)

func TestHandleResourceGroupLifeCycle(t *testing.T) {
	session := azure.NewSession(t)

	nash.Exec(t, "TestHandleResourceGroupLifeCycle", `
             azure login -q -u `+session.ClientID+` --service-principal --tenant `+session.TenantID+` -p `+session.ClientSecret+`

             import ../../azure/all
             azure_group_create("`+ResourceGroupName+`", "eastus")
        `)

	resources := azure.NewResources(t, session)
	defer resources.Delete(t, ResourceGroupName)
	resources.Check(t, ResourceGroupName)
}
