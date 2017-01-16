package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genRouteTableName() string {
	return fmt.Sprintf("klb-route-tests-%d", rand.Intn(1000))
}

func testRouteTableCreate(t *testing.T, f fixture.F) {
	route := genRouteTableName()
	f.Shell.Run(
		"./testdata/create_route_table.sh",
		route,
		f.ResGroupName,
		f.Location,
	)
	routes := azure.NewRouteTable(f)
	routes.AssertExists(t, route)
}

func TestRouteTable(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "RouteTable_Create", timeout, location, testRouTableteCreate)
}
