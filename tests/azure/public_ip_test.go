package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genPublicIpName() string {
	return fmt.Sprintf("klb-public-ip-tests-%d", rand.Intn(99999999))
}

func testPublicIpCreate(t *testing.T, f fixture.F) {
	publicIpName := genPublicIpName()

	f.Shell.Run(
		"./testdata/create_public_ip.sh",
		f.ResGroupName,
		publicIpName,
		f.Location,
	)
	publicIp := azure.NewPublicIp(f)

	publicIp.AssertExists(t, publicIpName)
}

func TestPublicIp(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "PublicIp_Create", timeout, location, testPublicIpCreate)
}
