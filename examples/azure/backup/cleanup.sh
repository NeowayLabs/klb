#!/usr/bin/env nash

import klb/azure/login
import klb/azure/group
import klb/azure/vm

# londing configs from config.sh
import config.sh

azure_login()
azure_group_delete($group)

backups <= azure_vm_backup_list($vm_name, $backup_prefix)

print("deleting backups: %q\n", $backups)

for backup in $backups {
	echo "deleting backup: " + $backup
	azure_vm_backup_delete($backup)
}
