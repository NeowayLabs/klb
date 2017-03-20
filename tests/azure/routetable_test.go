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
		f.ResGroupName,
		f.Location,
	)
	routeTables := azure.NewRouteTable(f)
	routeTables.AssertExists(t, routeTable)

	f.Shell.Run(
		"./testdata/add_route_to_route_table.sh",
		routeTable,
		route,
		f.ResGroupName,
		address,
		hoptype,
		hopaddress,
	)

	routes := azure.NewRoute(f)
	routes.AssertVirtualApplianceRouteExists(t, routeTable, route, address, hoptype, hopaddress)
}

func testRouteTableAddInternetRoute(t *testing.T, f fixture.F) {
	routeTable := genRouteTableName()
	route := routeTable + "route-internet-test"
	address := "0.0.0.0/0"
	hoptype := "Internet"

	f.Shell.Run(
		"./testdata/create_route_table.sh",
		routeTable,
		f.ResGroupName,
		f.Location,
	)
	routeTables := azure.NewRouteTable(f)
	routeTables.AssertExists(t, routeTable)

	f.Shell.Run(
		"./testdata/add_internet_route_to_route_table.sh",
		routeTable,
		route,
		f.ResGroupName,
		address,
		hoptype,
	)

	routes := azure.NewRoute(f)
	routes.AssertRouteExists(t, routeTable, route, address, hoptype)
}

func testRouteTableAddVirtualApplianceRoute(t *testing.T, f fixture.F) {
	routeTable := genRouteTableName()
	route := routeTable + "route-virtual-appliance-test"
	address := "0.0.0.0/0"
	hoptype := "VirtualAppliance"
	hopaddress := "10.116.1.100"

	f.Shell.Run(
		"./testdata/create_route_table.sh",
		routeTable,
		f.ResGroupName,
		f.Location,
	)
	routeTables := azure.NewRouteTable(f)
	routeTables.AssertExists(t, routeTable)

	f.Shell.Run(
		"./testdata/add_virtual_appliance_route_to_route_table.sh",
		routeTable,
		route,
		f.ResGroupName,
		address,
		hoptype,
		hopaddress,
	)

	routes := azure.NewRoute(f)
	routes.AssertVirtualApplianceRouteExists(t, routeTable, route, address, hoptype, hopaddress)
}

func TestRouteTable(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "RouteTable_Create", timeout, location, testRouteTableCreate)
	fixture.Run(t, "RouteTable_AddInternetRoute", timeout, location, testRouteTableAddInternetRoute)
	fixture.Run(t, "RouteTable_AddVirtualApplianceRoute", timeout, location, testRouteTableAddVirtualApplianceRoute)
}
