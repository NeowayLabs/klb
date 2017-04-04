Vnet with Private an Public Subnets
==

This scenario creates a Vnet with two subnets, one private and other private.
The Public Subnet has a route table with default route to "Internet" and Private subnet a route table with default route to Nat Virtual Appliance.

Virtual Appliance into Public Subnet have access granted from/to "Internet".

Virtual Appliance into Private Subnet have access to "Internet" through Nat
Appliance, and don't have access with origen from the "Internet", If you need connect using ssh per example, we need make a tunnel using Bastion Virtual Applience.

For expose your service running in Private Subnet you need a ALB (Azure Load Balance) or a reverse proxy (like Haproxy) into Public Subnet.

Access betwen Public Subnet and Private Subnet is granted by default.

In this scenario you will create a VNet named vnet-pub-priv with a reserved CIDR
block of 10.50.0.0./16.

Your VNet will contain the following subnets:

Public, using 10.50.1.0/24 as its CIDR block.
Private, using 10.50.2.0/24 as its CIDR block.

And three Virtual Appliances:
Nat in Public Subnet, using 10.50.1.100 as its address.
Bastion in Public Subnet, using 10.50.1.200 as its address.
App in Private Subnet, using address given by dhcp server.
