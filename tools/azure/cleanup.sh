#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group

klbtests_prefix = "klb-*"

echo
echo "========================================================================="
echo "this script will attempt to cleanup all pending klb tests resource groups"

azure_login()

resgroups <= azure_group_get_names()

filtered = ()

for resgroup in $resgroups {
	_ <= echo $resgroup | -grep $klbtests_prefix

	if $status == "0" {
		filtered <= append($filtered, $resgroup)
	}
}

echo
echo "going to delete the following resource groups:"
echo "==============================================="
echo

for resgroup in $filtered {
	echo $resgroup
}

echo
echo "==============================================="
echo
echo "from credentials:"
echo "==============================================="
echo
echo "AZURE_SUBSCRIPTION_NAME: "+$AZURE_SUBSCRIPTION_NAME
echo "AZURE_TENANT_ID: "+$AZURE_TENANT_ID
echo "AZURE_CLIENT_ID: "+$AZURE_CLIENT_ID
echo "AZURE_SERVICE_PRINCIPAL"+$AZURE_SERVICE_PRINCIPAL
echo
echo "==============================================="
echo
echo "are you sure you want to go on ?(Y/N)"

res <= head -n1 /dev/stdin

if $res != "Y" {
	echo "cancelling execution"
	
	exit("0")
}

echo "deleting resource groups, may the odds be at your favor"

for resgroup in $filtered {
	echo "deleting: "+$resgroup

	azure_group_delete($resgroup)
}

echo "done"
