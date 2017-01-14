package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genSubnetName() string {
	return fmt.Sprintf("klb-subnet-tests-%d", rand.Intn(1000))
}

func testSubnetCreate(t *testing.T, f fixture.F) {
	subnet := genSubnetName()
	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		f.Location,
		"10.116.0.0/16",
	)
	subnets := azure.NewSubnetsSubnet(f)
	subnets.AssertExists(t, subnet)
}

func TestSubnetSet(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "Subnet_Create", timeout, location, testSubnetCreate)
}
