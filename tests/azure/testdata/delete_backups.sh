#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

vmname    = $ARGS[1]
bkpprefix = $ARGS[2]

azure_login()

echo "deleting all backups with prefix: "+$bkpprefix
echo "for vm: "+$vmname

backups <= azure_vm_backup_list($vmname, $bkpprefix)

for backup in $backups {
	echo "deleting backup: "+$backup

	err <= azure_vm_backup_delete($backup)
	if $err != "" {
		echo "unable to delete backup: " + $backup
		echo "error: " + $err
		exit("1")
	}

	echo "deleted backup: "+$backup
}

echo "done"
