package azure_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVMName() string {
	return fmt.Sprintf("klb-vm-tests-%d", rand.Intn(1000))
}

func testVMCreate(t *testing.T, f fixture.F) {
	vm := genVMName()

	f.Shell.Run(
		"./testdata/create_vm.sh",
		vm,
		f.ResGroupName,
		f.Location,
		"Linux",
	)
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm)
}

func TestVM(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "VM_Create", timeout, location, testVMCreate)
}
