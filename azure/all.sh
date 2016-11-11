IFS = ()

import availset
import group
import nic
import nsg
import route
import storage
import subnet
import vm
import vnet
import app
import provider
import public-ip

azure config mode arm

azure login -q -u $AZURE_CLIENT_ID --service-principal --tenant $AZURE_TENANT_ID -p $AZURE_CLIENT_SECRET
