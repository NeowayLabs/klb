#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

backup_resgroup = $ARGS[1]

azure_login()

echo "deleting backup: " + $backup_resgroup
azure_vm_backup_delete($backup_resgroup)
echo "deleted backup: " + $backup_resgroup
