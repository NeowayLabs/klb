#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm
import klb/azure/snapshot

fn addsuffix(name) {
	# Providing true uniqueness with the limits on the names is pretty hard :-)
	s <= head -n1 /dev/urandom | md5sum | tr -dc A-Za-z0-9 | cut -b 1-10

	return $name+"-"+$s
}

resgroup = $ARGS[1]
vmname   = $ARGS[2]
sku      = $ARGS[3]
results  = $ARGS[4]

azure_login()

ids      <= azure_vm_get_datadisks_ids($vmname, $resgroup)

for id in $ids {
	snapshot_name <= addsuffix("snapshot")

	echo "creating snapshot: "+$snapshot_name+" from id: "+$id

	snapshotid <= azure_snapshot_create($snapshot_name, $resgroup, $id, $sku)

	echo "created snapshot id: "+$snapshotid

	echo $snapshotid | tee --append $results
}

echo "done"
