package azure_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
	testlog "github.com/NeowayLabs/klb/tests/lib/log"
	"github.com/NeowayLabs/klb/tests/lib/nash"
)

func genResourceGroupName() string {
	return fmt.Sprintf("klb-resgroup-tests-%d", rand.Intn(999999))
}

func testResourceGroupCreate(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger, teardown := testlog.New(t, "testResourceGroupCreate")
	defer teardown()

	resgroup := genResourceGroupName()
	session := fixture.NewSession(t)
	resources := fixture.NewResourceGroup(ctx, t, session, logger)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		resources := fixture.NewResourceGroup(ctx, t, session, logger)
		resources.Delete(t, resgroup)
	}()

	shell := nash.New(ctx, t, logger)
	shell.Run(
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

	logger, teardown := testlog.New(t, "testResourceGroupDelete")
	defer teardown()

	resgroup := genResourceGroupName()
	session := fixture.NewSession(t)
	resources := fixture.NewResourceGroup(ctx, t, session, logger)

	shell := nash.New(ctx, t, logger)
	shell.Run(
		"./testdata/create_resource_group.sh",
		resgroup,
		location,
	)
	resources.AssertExists(t, resgroup)

	shell.Run(
		"./testdata/delete_resource_group.sh",
		resgroup,
		location,
	)
	resources.AssertDeleted(t, resgroup)
}

func TestResourceGroup(t *testing.T) {
	t.Parallel()
	t.Run("Create", testResourceGroupCreate)
	t.Run("Delete", testResourceGroupDelete)
}
