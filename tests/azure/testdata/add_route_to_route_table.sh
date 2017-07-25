#!/usr/bin/env nash

import klb/azure/all

routetable = $ARGS[1]
route      = $ARGS[2]
resgroup   = $ARGS[3]
address    = $ARGS[4]
hoptype    = $ARGS[5]
hopaddress = $ARGS[6]

azure_login()


fn get_route_id() {
routeid <= azure_route_table_route_get_id($name, $resgroup, $routetable)

	return $routeid
}

route_id <= get_route_id()
if $route_id == "" {
azure_route_table_add_route($route, $resgroup, $routetable, $address, $hoptype, $hopaddress)
}

route_id <= get_route_id()
if $route_id == "" {
    exit("1")
}
