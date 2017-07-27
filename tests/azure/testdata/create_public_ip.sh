#!/usr/bin/env nash

import klb/azure/public_ip
import klb/azure/login

resgroup   = $ARGS[1]
name       = $ARGS[2]
location   = $ARGS[3]
allocation = "Static"

if len($ARGS) == "5" {
	allocation = $ARGS[4]
}

azure_login()
azure_public_ip_create($name, $resgroup, $location, $allocation)

public_ip_address, err <= azure_public_ip_get_address($name, $resgroup)

if $err != "" {
	print("error: %s", $err)
	exit("1")
}
