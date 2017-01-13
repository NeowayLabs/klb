#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[1]
availset = $ARGS[2]

azure_availset_delete($availset, $resgroup)
