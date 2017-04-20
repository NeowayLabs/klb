package azure_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func testLoadBalancer(t *testing.T, f fixture.F) {
	const cidr = "10.120.0.0/16"
	const subnetaddr = "10.120.1.0/24"
	const lbname = "loadbalancer"
	const frontendIPName = "lbfrontendip"
	const lbPrivateIP = "10.120.1.4"
	const poolname = "lbpool"

	f.Shell.Run(
		"./testdata/create_alb.sh",
		f.ResGroupName,
		f.Location,
		cidr,
		subnetaddr,
		lbname,
		frontendIPName,
		lbPrivateIP,
		poolname,
	)

	loadbalancer := azure.NewLoadBalancers(f)
	loadbalancer.AssertExists(t, lbname, frontendIPName, lbPrivateIP, poolname)

	const tcpprobePort int32 = 8080
	const httpprobePort int32 = 8081

	probes := []azure.LoadBalancerProbe{
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
	}

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

	rules := []azure.LoadBalancerRule{
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
	}

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

func TestLoadBalancer(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "LoadBalancer", 20*time.Minute, location, testLoadBalancer)
}
