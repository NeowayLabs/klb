#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/vm

klbtests_prefix = "^klb-*"

echo
echo "========================================================================="
echo "this script will attempt to cleanup all pending klb tests resource groups"
echo "logging in"

azure_login()

echo "getting resource group names"

resgroups <= azure_group_get_names()

filtered = ()

for resgroup in $resgroups {
	_, err <= echo $resgroup | grep $klbtests_prefix

	if $err == "0" {
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
echo "AZURE_SUBSCRIPTION_ID: "+$AZURE_SUBSCRIPTION_ID
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

	out, status <= echo $resgroup | grep -- "-bkp-"
	if $status == "0" {
		echo "seems like a backup, deleting it"
		err <= azure_vm_backup_delete($resgroup)
		if $err != "" {
			echo "error deleting backup: " + $err
			exit ("1")
		}
	} else {
		echo "deleting resgroup: "+$resgroup
		azure_group_delete($resgroup)
	}
}

echo "done"
