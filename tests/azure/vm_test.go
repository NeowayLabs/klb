package azure_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

func genVMName() string {
	return fmt.Sprintf("klbvmtests%d", rand.Intn(1000))
}

func testVMCreate(t *testing.T, f fixture.F) {

	vm := genVMName()
	osType := "Linux"
	vmSize := "Standard_DS4_v2"
	username := "core"
	osDisk := "test.vhd"
	imageUrn := "OpenLogic:CentOS:7.2:7.2.20161026"
	customData := ""
	keyFile := "./testdata/key.pub"

	availSet, vnet, subnet, nic, storAcc := createVMResources(f)

	f.Shell.Run(
		"./testdata/create_vm.sh",
		vm,
		f.ResGroupName,
		f.Location,
		osType,
		vmSize,
		username,
		availSet,
		vnet,
		subnet,
		nic,
		storAcc,
		osDisk,
		imageUrn,
		customData,
		keyFile,
	)
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm)
}

func createVMResources(f fixture.F) (string, string, string, string, string) {

	availSet := genAvailSetName()
	vnet := genVnetName()
	subnet := genSubnetName()
	nic := genNicName()
	storAcc := genStorageAccountName()

	nsg := genNsgName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	addrnic := "10.116.1.100"

	f.Shell.Run(
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		availSet,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/create_vnet.sh",
		vnet,
		f.ResGroupName,
		f.Location,
		vnetAddress,
	)

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		subnet,
		f.ResGroupName,
		vnet,
		subnetAddress,
		nsg,
	)

	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		nic,
		f.Location,
		vnet,
		subnet,
		addrnic,
	)

	f.Shell.Run(
		"./testdata/create_storage_account.sh",
		f.ResGroupName,
		storAcc,
		f.Location,
	)

	return availSet, vnet, subnet, nic, storAcc
}

func TestVM(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "VM_Create", 25*time.Minute, location, testVMCreate)
}
