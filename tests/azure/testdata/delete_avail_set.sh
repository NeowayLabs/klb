#!/usr/bin/env nash

import ../../azure/all

resgroup = $ARGS[0]
availset = $ARGS[1]

azure_availset_delete($availset, $resgroup)
