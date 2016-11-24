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
import disk

azure config mode arm

tenantID = $AZURE_TENANT_ID
clientID = $AZURE_CLIENT_ID
secretID = $AZURE_CLIENT_SECRET

azure login -q -u $clientID --service-principal --tenant $tenantID -p $secretID
