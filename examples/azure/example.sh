#!/usr/bin/env nash

import nashlib/all
import klb/azure/all

GROUP    = "klb-tests"
LOCATION = "eastus"

fn create_prod() {
	azure_group_create($GROUP, $LOCATION)
	azure_vnet_create("klb-tests-vnet", $GROUP, $LOCATION, "192.168.2.0/24")
}

fn delete_prod() {
	azure_group_delete($GROUP)
}

create_prod()
delete_prod()
