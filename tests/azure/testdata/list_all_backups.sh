#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

prefix = $ARGS[1]
output = $ARGS[2]

azure_login()

backups <= azure_vm_backup_list_all($prefix)
for backup in $backups {
	echo "got backup: " + $backup
	echo $backup | tee --append $output >[1=]
}
