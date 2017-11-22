#!/usr/bin/env nash

import klb/azure/login
import klb/azure/storage

resgroup      = $ARGS[1]
accountname   = $ARGS[2]
containername = $ARGS[3]

azure_login()

exit(azure_storage_container_exists_by_resgroup($containername, $accountname, $resgroup))
