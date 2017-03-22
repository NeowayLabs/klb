package azure_test

import (
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func testLoadBalancer(t *testing.T, f fixture.F) {
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

	probes := []azure.LoadBalancerProbe{
		{
			Name:     "tcp",
			Protocol: "tcp",
			Port:     "7777",
			Interval: "60",
			Count:    "10",
		},
		{
			Name:     "http",
			Protocol: "http",
			Port:     "7776",
			Interval: "120",
			Count:    "20",
			Path:     "/healthz",
		},
	}

	for _, p := range probes {
		args := []string{
			f.ResGroupName,
			p.Name,
			lbname,
			p.Port,
			p.Protocol,
			p.Interval,
			p.Count,
		}
		if p.Path != "" {
			args = append(args, p.Path)
		}
		f.Shell.Run("./testdata/add_alb_probe.sh", args...)
		loadbalancer.AssertProbeExists(t, lbname, p)
	}
}

func TestLoadBalancer(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "LoadBalancer", timeout, location, testLoadBalancer)
}
