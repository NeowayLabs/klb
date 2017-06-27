# Public IP related functions

fn azure_public_ip_create(name, group, location, allocation) {
       (
               az network public-ip create --name $name --resource-group $group --location $location --allocation-method $allocation
       )
}

fn azure_public_ip_get_address(name, group) {
       public_ip_info    <= az network public-ip show --resource-group $group --name $name --output json
       public_ip_address <= echo $public_ip_info | jq -r ".ipAddress"

       if $public_ip_address == "" {
               return "", format("unable to find address on public_ip[%s] resgroup[%s], public_ip probably does not exist", $name, $group)
       }

       return $public_ip_address, ""
}

fn azure_public_ip_delete(name, group) {
       (
               az network public-ip delete --name $name --resource-group $group
       )
}
