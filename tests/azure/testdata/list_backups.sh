#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

vmname = $ARGS[1]
prefix = $ARGS[2]
output = $ARGS[3]

azure_login()

backups <= azure_vm_backup_list($vmname, $prefix)

for backup in $backups {
	echo "got backup: "+$backup

	echo $backup | tee --append $output >[1=]
}
