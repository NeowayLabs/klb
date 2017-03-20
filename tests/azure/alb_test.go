package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genALBName() string {
	return fmt.Sprintf("klb-alb-tests-%d", rand.Intn(1000))
}

func testALBCreate(t *testing.T, f fixture.F) {
	name := genALBName()
	f.Shell.Run(
		"./testdata/create_alb.sh",
		f.ResGroupName,
		name,
		f.Location,
	)
	loadbalancer := azure.NewLoadBalancer(f)
	loadbalancer.AssertExists(t, name)
}

//func testAvailSetDelete(t *testing.T, f fixture.F) {
//loadbalancer := genALBName()
//f.Shell.Run(
//"./testdata/create_avail_set.sh",
//f.ResGroupName,
//loadbalancer,
//f.Location,
//)

//availSets := azure.NewAvailSet(f)
//availSets.AssertExists(t, loadbalancer)

//f.Shell.Run(
//"./testdata/delete_avail_set.sh",
//f.ResGroupName,
//loadbalancer,
//)
//availSets.AssertDeleted(t, loadbalancer)
//}

func TestLoadBalancer(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "LoadBalancer_Create", timeout, location, testALBCreate)
}
