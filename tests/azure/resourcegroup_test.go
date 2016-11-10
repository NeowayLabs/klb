package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
	"github.com/NeowayLabs/klb/tests/nash"
)

const (
	ResourceGroupName = "klb-resource-group-test"
)

func TestHandleResourceGroupLifeCycle(t *testing.T) {
	session := azure.NewSession(t)

	nash.Exec(t, "TestHandleResourceGroupLifeCycle", `
             azure login -u `+session.ClientID+` --service-principal --tenant `+session.TenantID+` -p `+session.ClientSecret+`

             import ../../azure/all
             azure_group_create("`+ResourceGroupName+`", "eastus")
        `)

	resources := azure.NewResources(t, session)
	defer resources.Delete(t, ResourceGroupName)
	resources.Check(t, ResourceGroupName)
}
