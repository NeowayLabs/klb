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
		"testResourceGroupCreate",
		"./testdata/create_resource_group.sh",
		resgroup,
		location,
	)
	resources.AssertExists(t, resgroup)
}

func testResourceGroupDelete(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resgroup := genResourceGroupName()
	session := azure.NewSession(t)
	resources := azure.NewResourceGroup(ctx, t, session)

	nash.Run(
		ctx,
		t,
		"testResourceGroupDelete",
		"./testdata/create_resource_group.sh",
		resgroup,
		location,
	)
	resources.AssertExists(t, resgroup)

	nash.Run(
		ctx,
		t,
		"testResourceGroupDelete",
		"./testdata/delete_resource_group.sh",
		resgroup,
		location,
	)
	resources.AssertDeleted(t, resgroup)
}
