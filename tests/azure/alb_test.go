package azure_test

import (
	"strconv"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func TestLoadBalancer(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "LoadBalancer", timeout, location, testLoadBalancer)
}

func testLoadBalancer(t *testing.T, f fixture.F) {

	const nsg = "lbnsg"
	const vnet = "lbvnet"
	const vnetCIDR = "10.120.0.0/16"
	const subnet = "lbsubnet"
	const subnetCIDR = "10.120.1.0/24"
	const lbname = "loadbalancer"
	const frontendIPName = "lbfrontendip"
	const lbPrivateIP = "10.120.1.4"
	const poolname = "lbpool"

	createVNET(t, f, vnetDescription{name: vnet, vnetAddr: vnetCIDR})
	createNSG(t, f, nsg)
	createSubnet(t, f, vnet, subnet, subnetCIDR, nsg)
	createLoadBalancer(t, f, vnet, subnet, lbname, frontendIPName, lbPrivateIP, poolname)

	loadbalancer := azure.NewLoadBalancers(f)
	loadbalancer.AssertExists(t, lbname, frontendIPName, lbPrivateIP, poolname)

	const tcpprobePort int32 = 8080
	const httpprobePort int32 = 8081

	createLoadBalancerProbes(t, f, lbname, []azure.LoadBalancerProbe{
		{
			Name:     "tcpprobe",
			Protocol: "Tcp",
			Port:     tcpprobePort,
			Interval: 60,
			Count:    10,
		},
		{
			Name:     "httpprobe",
			Protocol: "Http",
			Port:     httpprobePort,
			Interval: 120,
			Count:    20,
			Path:     "/healthz",
		},
	})

	createLoadBalancerRules(t, f, lbname, frontendIPName, poolname, []azure.LoadBalancerRule{
		{
			Name:         "tcprule",
			ProbeName:    "tcpprobe",
			Protocol:     "Tcp",
			FrontendPort: tcpprobePort,
			BackendPort:  tcpprobePort,
		},
		{
			Name:         "httprule",
			ProbeName:    "httpprobe",
			Protocol:     "Tcp",
			FrontendPort: httpprobePort,
			BackendPort:  httpprobePort,
		},
	})

}

func createLoadBalancer(
	t *testing.T,
	f fixture.F,
	vnet string,
	subnet string,
	lbname string,
	frontendIPName string,
	privateIP string,
	poolname string,
) {
	f.Shell.Run(
		"./testdata/create_alb.sh",
		f.ResGroupName,
		f.Location,
		vnet,
		subnet,
		lbname,
		frontendIPName,
		privateIP,
		poolname,
	)
}

func createLoadBalancerProbes(
	t *testing.T,
	f fixture.F,
	lbname string,
	probes []azure.LoadBalancerProbe,
) {
	loadbalancer := azure.NewLoadBalancers(f)
	for _, p := range probes {
		args := []string{
			f.ResGroupName,
			p.Name,
			lbname,
			strconv.Itoa(int(p.Port)),
			p.Protocol,
			strconv.Itoa(int(p.Interval)),
			strconv.Itoa(int(p.Count)),
		}
		if p.Path != "" {
			args = append(args, p.Path)
		}
		f.Shell.Run("./testdata/add_alb_probe.sh", args...)
		loadbalancer.AssertProbeExists(t, lbname, p)
	}
}

func createLoadBalancerRules(
	t *testing.T,
	f fixture.F,
	lbname string,
	frontendIPName string,
	poolname string,
	rules []azure.LoadBalancerRule,
) {
	loadbalancer := azure.NewLoadBalancers(f)
	for _, r := range rules {
		args := []string{
			f.ResGroupName,
			r.Name,
			lbname,
			r.ProbeName,
			frontendIPName,
			poolname,
			r.Protocol,
			strconv.Itoa(int(r.FrontendPort)),
			strconv.Itoa(int(r.BackendPort)),
		}
		f.Shell.Run("./testdata/add_alb_rule.sh", args...)
		loadbalancer.AssertRuleExists(t, lbname, r)
	}
}

func getLBAddressPoolID(t *testing.T, f fixture.F, lbname string, poolname string) string {
	return execWithIPC(t, f, func(outfile string) {
		f.Shell.Run(
			"./testdata/alb_get_pool_id.sh",
			poolname,
			f.ResGroupName,
			lbname,
			outfile,
		)
	})
}
