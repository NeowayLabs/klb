package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func testALBCreate(t *testing.T, f fixture.F) {
	const cidr = "10.120.0.0/16"
	const subnetaddr = "10.120.1.0/24"
	const lbname = "loadbalancer"
	const frontendip_name = "lbfrontendip"
	const lb_private_ip = "10.120.1.4"
	const poolname = "lbpool"

	f.Shell.Run(
		"./testdata/create_alb.sh",
		f.ResGroupName,
		f.Location,
		cidr,
		subnetaddr,
		lbname,
		frontendip_name,
		lb_private_ip,
		poolname,
	)

	loadbalancer := azure.NewLoadBalancers(f)
	loadbalancer.AssertExists(t, lbname, frontendip_name, lb_private_ip, poolname)
}

func TestLoadBalancer(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "LoadBalancer_Create", timeout, location, testALBCreate)
}
