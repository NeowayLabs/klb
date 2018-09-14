#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/nic
import klb/azure/subnet
import klb/azure/vm
import klb/azure/availset
import klb/azure/storage
import klb/azure/disk
import klb/azure/vnet
import klb/azure/nsg
import klb/azure/route
import klb/azure/snapshot
import config.sh


fn addsuffix(name) {
	# Providing true uniqueness with the limits on the names is pretty hard :-)
	s <= head -n1 /dev/urandom | md5sum | tr -dc A-Za-z0-9 | cut -b 1-10

	return $name+"-"+$s
}

fn log(msg) {
	ts <= date "+%T"
	echo $ts + ":" + $msg
}

fn new_disk(group, location) {
    name <= addsuffix("disk")
    d <= azure_disk_new($name, $group, $location)
	d <= azure_disk_set_size($d, $disk_size)
    d <= azure_disk_set_sku($d, $sku)
    return azure_disk_create($d)
}

azure_login()

echo "creating resource groups"

azure_group_create($group, $location)
azure_group_create($other_group, $other_location)


disk1 <= new_disk($group, $location)
disk2 <= new_disk($group, $location)

snapshot1 <= addsuffix("snapshot1")
snapshot2 <= addsuffix("snapshot2")

snapshot1_id <= azure_snapshot_create($snapshot1, $group, $disk1, $sku)
snapshot2_id <= azure_snapshot_create($snapshot2, $group, $disk2, $sku)

log(format("creates snapshots: [%s] [%s]", $snapshot1_id, $snapshot2_id))

copied_snapshots_ids, err <= azure_snapshot_copy($other_group, ($snapshot1_id $snapshot2_id))

if $err != "" {
    log(format("error[%s] copying snapshots between regions", $err))
    exit("1")
}

log(format("copied snapshots: %s", $copied_snapshots_ids))

az snapshot show --ids $copied_snapshots_ids

echo
echo "finished with no errors lol"
