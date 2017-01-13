#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[0]
availset = $ARGS[1]

azure_storage_account_delete($availset, $resgroup)
