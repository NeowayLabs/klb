#!/usr/bin/env nash

import klb/azure/all

routetable = $ARGS[1]
name       = $ARGS[2]
resgroup   = $ARGS[3]
address    = $ARGS[4]
hoptype    = $ARGS[5]

azure_login()

fn get_route_table_id() {
	routetableid, status <= azure_route_table_get_id($routetable, $resgroup)

	return $routetableid
}

route_id <= get_route_table_id()

if $route_id == "" {
	route <= azure_route_table_route_new($name, $resgroup, $routetable, $address, $hoptype)

	azure_route_table_route_create($route)
}

route_id <= get_route_table_id()

if $route_id == "" {
	exit("1")
}
