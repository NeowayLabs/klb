#!/usr/bin/env nash
import klb/azure/login
import klb/azure/group
import klb/azure/postgres

## Resource Group Settings
group    = "klb-examples-postgres"
location = "eastus2"

name     = "klb-pg-test"
username = "klbpgtestuser"
password = "klb2PgPass@secret"
units    = "100"
tier     = "Basic"
ssl      = "false"
size     = "51200"
version  = "9.5"

azure_login()

echo "creating new resource group"

azure_group_create($group, $location)

echo "creating Postgres Server"

server <= azure_postgres_new($name, $group, $location, $username, $password)
server <= azure_postgres_set_compute_units($server, $units)
server <= azure_postgres_set_performance_tier($server, $tier)
server <= azure_postgres_set_max_size($server, $size)
server <= azure_postgres_set_version($server, $version)

if $ssl == "false" {
        server <= azure_postgres_disable_ssl($server)
}

azure_postgres_server_create($server)

echo "created Postgres Server"

echo
echo "finished with no errors lol"
