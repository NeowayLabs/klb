#!/usr/bin/env nash

import ../../azure/all

routetable = $ARGS[1]
name       = $ARGS[2]
resgroup   = $ARGS[3]
address    = $ARGS[4]
hoptype    = $ARGS[5]
hopaddress = $ARGS[6]

azure_login()

route <= azure_route_table_route_new($name, $resgroup, $routetable, $address, $hoptype)
route <= azure_route_table_route_set_hop_address($route, $hopaddress)

azure_route_table_route_create($route)
