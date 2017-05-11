#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

vm_name = $ARGS[1]
resgroup = $ARGS[2]
prefix = $ARGS[3]
location = $ARGS[4]
output = $ARGS[5]

azure_login()

print("vm name %q resgroup %q prefix %q location %q\n", $vm_name, $resgroup, $prefix, $location)
echo "creating backup"
backup, err <= azure_vm_backup_create($vm_name, $resgroup, $prefix, $location)
if $err != "" {
	echo $err
	exit("1")
}
echo "created backup: " + $backup
echo $backup > $output
