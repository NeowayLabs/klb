package azure_test

import (
	"os"
	"testing"

	"github.com/NeowayLabs/klb/tests/azure"
	"github.com/NeowayLabs/nash"
)

const (
	ResourceGroupName = "klb-resource-group-test"
)

func TestHandleResourceGroupLifeCycle(t *testing.T) {
	session := azure.NewSession(t)

	sh, err := nash.New()

	if err != nil {
		t.Fatal(err)
	}

	sh.SetStdout(os.Stdout)
	sh.SetStderr(os.Stderr)

	err = sh.Exec("TestHandleResourceGroupLifeCycle", `
             import ../../azure/all
             azure login -u `+session.ClientID+` --service-principal --tenant `+session.TenantID+` -p `+session.ClientSecret+`
             azure_group_create("`+ResourceGroupName+`", "eastus")
        `)

	if err != nil {
		t.Fatal(err)
	}

	resources := azure.NewResources(t, session)
	defer resources.Delete(t, ResourceGroupName)
	resources.Check(t, ResourceGroupName)
}
