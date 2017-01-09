package azure_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/nash"
)

func genResourceGroupName() string {
	return fmt.Sprintf("klb-resgroup-tests-%d", rand.Intn(999999))
}

func testResourceGroupCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resgroup := genResourceGroupName()
	session := azure.NewSession(t)
	resources := azure.NewResourceGroup(ctx, t, session)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		resources := azure.NewResourceGroup(ctx, t, session)
		resources.Delete(t, resgroup)
	}()

	nash.Run(
		ctx,
		t,
		"./testdata/create_resource_group.sh",
		resgroup,
		location,
	)
	resources.AssertExists(t, resgroup)
}

func testResourceGroupDelete(t *testing.T) {
	//t.Parallel()

	//ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//defer cancel()

	//resgroup := genResourceGroupName()
	//session := azure.NewSession(t)
	//resources := azure.NewResourceGroup(ctx, t, session)

	//shell := nash.New(t)
	//shell.Setvar("ResourceGroup", resgroup)

	//err := shell.Exec("TestResourceGroupDeletion", `
	//import ../../azure/all
	//azure_group_create($ResourceGroup, "eastus")
	//`)

	//if err != nil {
	//t.Fatal(err)
	//}

	//resources.AssertExists(t, resgroup)

	//err = shell.Exec("TestResourceGroupDeletion2", `
	//azure_group_delete($ResourceGroup)
	//`)

	//if err != nil {
	//t.Error(err)
	//}

	//resources.AssertDeleted(t, resgroup)
}
