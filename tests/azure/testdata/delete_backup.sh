#!/usr/bin/env nash

import klb/azure/login
import klb/azure/vm

backup = $ARGS[1]

azure_login()

echo "deleting backup: " + $backup

err <= azure_vm_backup_delete($backup)
if $err != "" {
        echo "unable to delete backup: " + $backup
        echo "error: " + $err
        exit("1")
}

echo "done"
