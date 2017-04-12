package azure_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func diskname() string {
	return fmt.Sprintf("klbdisktests%d", rand.Intn(99999))
}

func testDiskCreate(t *testing.T, f fixture.F) {

	//name := diskname()
	//size := "10"
	//sku := "Standard_LRS"

	//f.Shell.Run(
	//"./testdata/create_managed_disk.sh",
	//f.ResGroupName,
	//f.Location,
	//name,
	//size,
	//sku,
	//)
	//disk := azure.NewDisk(f)
	//disk.AssertExists(t, name, size, sku)
}

func TestDisk(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "DiskCreate", 10*time.Minute, location, testDiskCreate)
}
