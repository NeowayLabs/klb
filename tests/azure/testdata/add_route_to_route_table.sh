#!/usr/bin/env nash

import klb/azure/all

routetable = $ARGS[1]
route      = $ARGS[2]
resgroup   = $ARGS[3]
address    = $ARGS[4]
hoptype    = $ARGS[5]
hopaddress = $ARGS[6]

azure_login()
azure_route_table_add_route($route, $resgroup, $routetable, $address, $hoptype, $hopaddress)
