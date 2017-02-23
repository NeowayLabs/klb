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

type VMResources struct {
	availSet string
	vnet     string
	subnet   string
	nic      string
	storAcc  string
}

func testVMCreate(t *testing.T, f fixture.F) {

	vm := genVMName()
	osType := "Linux"
	vmSize := "Standard_DS4_v2"
	username := "core"
	osDisk := "test.vhd"
	imageUrn := "OpenLogic:CentOS:7.2:7.2.20161026"
	keyFile := "./testdata/key.pub"

	resources := createVMResources(t, f)

	f.Shell.Run(
		"./testdata/create_vm.sh",
		vm,
		f.ResGroupName,
		f.Location,
		osType,
		vmSize,
		username,
		resources.availSet,
		resources.vnet,
		resources.subnet,
		resources.nic,
		resources.storAcc,
		osDisk,
		imageUrn,
		keyFile,
	)
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm, resources.availSet, vmSize, osType, resources.nic)
}

func createVMResources(t *testing.T, f fixture.F) VMResources {

	resources := VMResources{}
	resources.availSet = genAvailSetName()
	resources.vnet = genVnetName()
	resources.subnet = genSubnetName()
	resources.nic = genNicName()
	resources.storAcc = genStorageAccountName()

	nsg := genNsgName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	addrnic := "10.116.1.100"

	f.Shell.Run(
		"./testdata/create_avail_set.sh",
		f.ResGroupName,
		resources.availSet,
		f.Location,
	)

	createVNet(t, f, vnetDescription{
		name:     resources.vnet,
		vnetAddr: vnetAddress,
	})

	f.Shell.Run(
		"./testdata/create_nsg.sh",
		nsg,
		f.ResGroupName,
		f.Location,
	)

	f.Shell.Run(
		"./testdata/create_subnet.sh",
		resources.subnet,
		f.ResGroupName,
		resources.vnet,
		subnetAddress,
		nsg,
	)

	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		resources.nic,
		f.Location,
		resources.vnet,
		resources.subnet,
		addrnic,
	)

	f.Shell.Run(
		"./testdata/create_storage_account.sh",
		f.ResGroupName,
		resources.storAcc,
		f.Location,
	)

	return resources
}

func TestVM(t *testing.T) {
	t.Parallel()
	fixture.Run(t, "VM_Create", 25*time.Minute, location, testVMCreate)
}
