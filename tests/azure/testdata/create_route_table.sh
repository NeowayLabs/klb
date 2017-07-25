#!/usr/bin/env nash

import klb/azure/all

routetable = $ARGS[1]
resgroup   = $ARGS[2]
location   = $ARGS[3]

azure_login()

fn get_route_table_id() {
	routetableid <= azure_route_table_get_id($routetable, $resgroup)

	return $routetableid
}

route_table_id <= get_route_table_id()

if $route_table_id != "" {
	azure_route_table_create($routetable, $resgroup, $location)
}

route_table_id <= get_route_table_id()

if $route_table_id == "" {
	exit("1")
}
