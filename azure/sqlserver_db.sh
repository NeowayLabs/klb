# SQL Servers DB related functions

# azure_sqlserverdbdb_new creates a new SQL Server DB on Azure.
# `name` is the name of the Azure SQL Database.
# `group` is name of resource group.
# `servername` is the name of the Azure SQL server.
fn azure_sqlserverdb_new(name, group, servername) {
	instance = (
		"--name"
		$name
		"--resource-group"
		$group
		"--server"
		$servername
	)

	return $instance
}

# azure_sqlserverdb_set_collation set the collation of the database.
# `instance` is the sqlserverdb instance.
# `collation` is the collation desired for the database (e.g.: utf8mb4_general_ci).
# Default: SQL_Latin1_General_CP1_CI_AS
fn azure_sqlserverdb_set_collation(instance, collation) {
	instance <= append($instance, "--collation")
	instance <= append($instance, $collation)

	return $instance
}

# azure_sqlserverdb_set_service_objective set the name of the configured service level objective of the database. 
# `instance` is the sqlserverdb server instance.
# `serviceobjective` is the performance tier desired for the server.
# Allowed values: Basic, Standard (S0,S1,S2,S3), Premium (P1,P2,P4,P6,P11,P15), PremiumRS (PRS1,PRS2,PRS4,PRS6)
fn azure_sqlserverdb_set_service_objective(instance, serviceobjective) {
	instance <= append($instance, "--service-objective")
	instance <= append($instance, $serviceobjective)

	return $instance
}

# azure_sqlserverdb_set_edition set the edition of the database.
# `instance` is the sqlserverdb instance.
# `edition` is the edition desired.
# Default: SQL_Latin1_General_CP1_CI_AS
fn azure_sqlserverdb_set_edition(instance, edition) {
	instance <= append($instance, "--edition")
	instance <= append($instance, $edition)

	return $instance
}


# azure_sqlserverdb_set_max_size sets the max storage size for the server.
# `instance` is the sqlserverdb server instance.
# `size` is the sqlserverdb server storage size.
# Allowed values: 100MB, 500MB, 1GB, 5GB, 10GB, 20GB, 30GB, 150GB, 200GB, 500GB
fn azure_sqlserverdb_set_max_size(instance, size) {
	instance <= append($instance, "--max-size")
	instance <= append($instance, $size)

	return $instance
}

# azure_sqlserverdb_set_elastic_pool set the elastic pool the database is in. Not supported for DataWarehouse edition.
# `instance` is the sqlserverdb instance.
# `elasticpool` is the edition desired.
# Default: SQL_Latin1_General_CP1_CI_AS
fn azure_sqlserverdb_set_elastic_pool(instance, elasticpool) {
	instance <= append($instance, "--elastic-pool")
	instance <= append($instance, $elasticpool)

	return $instance
}

# azure_sqlserverdb_set_sample_name sets the name of the sample schema to apply when creating this database. Not supported for DataWarehouse edition.
# `instance` is the sqlserverdb server instance.
# `samplename` is the name of the sample schema.
fn azure_sqlserverdb_set_sample_name(instance, samplename) {
	instance <= append($instance, "--sample-name")
	instance <= append($instance, $samplename)

	return $instance
}

# azure_sqlserverdb_set_tags set the resource tags.
# `instance` is the sqlserverdb server instance.
# `tags` is the name of the tag.
fn azure_sqlserverdb_set_tags(instance, tags) {
	instance <= append($instance, "--tags")
	instance <= append($instance, $tags)

	return $instance
}

# azure_sqlserverdb_server_create creates a new SQLServer Database on Azure.
# `instance` is the sqlserverdb server instance.
fn azure_sqlserverdb_server_create(instance) {
        az sql db create $instance
}
