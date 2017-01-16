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
	routeTable := genRouteTableName()
	route := routeTable + "route-test"
	address := "0.0.0.0/0"
	hoptype := "VirtualAppliance"
	hopaddress := "10.116.1.100"

	f.Shell.Run(
		"./testdata/create_route_table.sh",
		routeTable,
		route,
		f.ResGroupName,
		f.Location,
		address,
		"VirtualAppliance",
		"10.116.1.100",
	)
	routeTables := azure.NewRouteTable(f)
	routeTables.AssertExists(t, routeTable)

	routes := azure.NewRoute(f)
	routes.AssertExists(t, routeTable, route, address, hoptype, hopaddress)
}

func TestRouteTable(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "RouteTable_Create", timeout, location, testRouteTableCreate)
}
