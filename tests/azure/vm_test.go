package azure_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/NeowayLabs/klb/tests/lib/azure"
	"github.com/NeowayLabs/klb/tests/lib/azure/fixture"
)

type VMResources struct {
	availSet string
	vnet     string
	subnet   string
	nic      string
}

func TestVM(t *testing.T) {
	t.Parallel()
	vmtesttimeout := 45 * time.Minute
	fixture.Run(t, "VMStandardDisk", vmtesttimeout, location, testStandardDiskVM)
	fixture.Run(t, "VMPremiumDisk", vmtesttimeout, location, testPremiumDiskVM)
	fixture.Run(t, "VMPremiumDiskReadCache", vmtesttimeout, location, testVMPremiumDiskReadCache)
	fixture.Run(t, "VMPremiumDiskRWCache", vmtesttimeout, location, testVMPremiumDiskRWCache)
	fixture.Run(t, "VMSnapshotStandard", vmtesttimeout, location, testVMSnapshotStandard)
	fixture.Run(t, "VMSnapshotPremium", vmtesttimeout, location, testVMSnapshotPremium)
	fixture.Run(t, "VMPremiumDiskToStdSnapshot", vmtesttimeout, location, testVMPremiumDiskToStdSnapshot)
	fixture.Run(t, "VMDuplicatedAvSet", vmtesttimeout, location, testDuplicatedAvailabilitySet)
}

func genUniqName() string {
	return fixture.NewUniqueName("vm")
}

func createVM(
	t *testing.T,
	f fixture.F,
	availset string,
	nic string,
	vmSize string,
	sku string,
	caching string,
) string {
	vm := genUniqName()
	username := "core"
	osDisk := genUniqName()
	imageUrn := "OpenLogic:CentOS:7.2:7.2.20161026"
	keyFile := "./testdata/key.pub"

	f.Shell.Run(
		"./testdata/create_vm.sh",
		vm,
		f.ResGroupName,
		f.Location,
		vmSize,
		username,
		availset,
		nic,
		osDisk,
		imageUrn,
		keyFile,
		sku,
		caching,
	)

	f.Logger.Println("creating VM")
	vms := azure.NewVM(f)
	vms.AssertExists(t, vm, availset, vmSize, nic)
	f.Logger.Println("created VM with success, attaching a disk")
	return vm
}

func testVMCreation(
	t *testing.T,
	f fixture.F,
	vmSize string,
	sku string,
	caching string,
) {

	resources := createVMResources(t, f)
	vm := createVM(
		t,
		f,
		resources.availSet,
		resources.nic,
		vmSize,
		sku,
		caching,
	)

	vms := azure.NewVM(f)
	osdisk := vms.OsDisk(t, vm)

	if caching != "None" {
		// Why: OS Disks don't support None caching,
		// we ignore when the caching config is None
		if caching != osdisk.Caching {
			t.Fatalf(
				"expected osdisk caching to be[%s] but it is [%s]",
				caching,
				osdisk.Caching,
			)
		}
	}

	diskname := "createVMExtraDisk"
	size := 10

	attachNewDiskOnVM(t, f, vm, diskname, size, sku, caching)

	vms.AssertAttachedDataDisk(t, vm, diskname, size, sku, caching)
}

func testPremiumDiskVM(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Standard_DS4_v2", "Premium_LRS", "None")
}

func testVMPremiumDiskReadCache(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Standard_DS4_v2", "Premium_LRS", "ReadOnly")
}

func testVMPremiumDiskRWCache(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Standard_DS4_v2", "Premium_LRS", "ReadWrite")
}

func testStandardDiskVM(t *testing.T, f fixture.F) {
	testVMCreation(t, f, "Basic_A0", "Standard_LRS", "None")
}

type VMDisk struct {
	Name    string
	Sku     string
	Size    int
	Caching string
}

func validateDisk(t *testing.T, disk VMDisk) {
	if disk.Name == "" {
		t.Fatal("disk name cant be empty")
	}
	if disk.Sku == "" {
		t.Fatal("disk sku cant be empty")
	}
	if disk.Caching == "" {
		t.Fatal("disk caching cant be empty")
	}
}

func attachDisks(t *testing.T, f fixture.F, vmname string, disks []VMDisk) {
	vms := azure.NewVM(f)
	for _, disk := range disks {
		validateDisk(t, disk)
		attachNewDiskOnVM(t, f, vmname, disk.Name, disk.Size, disk.Sku, disk.Caching)
		vms.AssertAttachedDataDisk(t, vmname, disk.Name, disk.Size, disk.Sku, disk.Caching)
	}
}

func testVMSnapshot(
	t *testing.T,
	f fixture.F,
	vmSize string,
	vmSKU string,
	snapshotSKU string,
	disks []VMDisk,
) {
	resources := createVMResources(t, f)
	vm := createVM(t, f, resources.availSet, resources.nic, vmSize, vmSKU, "None")

	vms := azure.NewVM(f)

	attachDisks(t, f, vm, disks)

	idsraw := execWithIPC(t, f, func(outfile string) {
		f.Shell.Run(
			"./testdata/create_vm_snapshots.sh",
			f.ResGroupName,
			vm,
			snapshotSKU,
			outfile,
		)
	})

	f.Logger.Println("created snapshots, retrieving ids")
	ids := strings.Split(strings.Trim(string(idsraw), "\n"), "\n")
	f.Logger.Printf("parsed ids: %s", ids)

	if len(ids) != len(disks) {
		t.Fatalf("expected %d snapshots, got %d", len(disks), len(ids))
	}

	nic := genNicName()
	createVMNIC(f, nic, resources.vnet, resources.subnet)
	vmbackup := createVM(t, f, resources.availSet, nic, vmSize, vmSKU, "None")

	for _, id := range ids {
		attachSnapshotOnVM(t, f, vmbackup, id, vmSKU)
	}

	originaldisks := vms.DataDisks(t, vm)
	if len(originaldisks) != len(disks) {
		t.Fatalf("expected %d disks, got %d", len(disks), len(originaldisks))
	}

	backupdisks := vms.DataDisks(t, vmbackup)
	if len(originaldisks) != len(backupdisks) {
		t.Fatalf("expected original disks %q == %q backup disks", originaldisks, backupdisks)
	}

	for _, originaldisk := range originaldisks {
		size := originaldisk.SizeGB
		got := false
		for _, backupdisk := range backupdisks {
			if backupdisk.SizeGB == size {
				got = true
			}
		}
		if !got {
			t.Fatalf(
				"unable to find disk: %q on backup disks %q",
				originaldisk,
				backupdisks,
			)
		}
	}
}

func testVMSnapshotStandard(t *testing.T, f fixture.F) {
	sku := "Standard_LRS"
	testVMSnapshot(t, f, "Basic_A2", sku, sku,
		[]VMDisk{
			{Name: genUniqName(), Size: 20, Sku: sku, Caching: "None"},
			{Name: genUniqName(), Size: 30, Sku: sku, Caching: "None"},
		},
	)

}

func testVMSnapshotPremium(t *testing.T, f fixture.F) {
	sku := "Premium_LRS"
	testVMSnapshot(t, f, "Standard_DS4_v2", sku, sku,
		[]VMDisk{
			{Name: genUniqName(), Size: 50, Sku: sku, Caching: "None"},
			{Name: genUniqName(), Size: 150, Sku: sku, Caching: "None"},
		},
	)
}

func testVMPremiumDiskToStdSnapshot(t *testing.T, f fixture.F) {
	vmSKU := "Premium_LRS"
	snapshotSKU := "Standard_LRS"

	testVMSnapshot(t, f, "Standard_DS4_v2", vmSKU, snapshotSKU,
		[]VMDisk{{Name: genUniqName(), Size: 50, Sku: vmSKU, Caching: "None"}},
	)
}

func testDuplicatedAvailabilitySet(t *testing.T, f fixture.F) {
	name := "duplicatedAvSet"

	createAvSet := func() {
		f.Shell.Run(
			"./testdata/create_vm_avail_set.sh",
			f.ResGroupName,
			name,
			f.Location,
			"1",
			"1",
		)
	}

	createAvSet()
	availSets := azure.NewAvailSet(f)
	availSets.AssertExists(t, name)

	createAvSet()
	availSets.AssertExists(t, name)
}

func attachSnapshotOnVM(
	t *testing.T,
	f fixture.F,
	vmname string,
	snapshotid string,
	disksku string,
) string {
	diskname := genUniqName()
	f.Shell.Run(
		"./testdata/attach_snapshot.sh",
		f.ResGroupName,
		f.Location,
		vmname,
		diskname,
		disksku,
		snapshotid,
	)
	return diskname
}

func createVMResources(t *testing.T, f fixture.F) VMResources {

	resources := VMResources{}
	resources.availSet = genAvailSetName()
	resources.vnet = genVnetName()
	resources.subnet = genSubnetName()
	resources.nic = genNicName()

	nsg := genNsgName()
	vnetAddress := "10.116.0.0/16"
	subnetAddress := "10.116.1.0/24"
	updatedomain := "3"
	faultdomain := "3"

	f.Shell.Run(
		"./testdata/create_vm_avail_set.sh",
		f.ResGroupName,
		resources.availSet,
		f.Location,
		updatedomain,
		faultdomain,
	)

	createVNET(t, f, vnetDescription{
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

	createVMNIC(f, resources.nic, resources.vnet, resources.subnet)

	return resources
}

func createVMNIC(f fixture.F, nic string, vnet string, subnet string) {
	f.Shell.Run(
		"./testdata/create_nic.sh",
		f.ResGroupName,
		nic,
		f.Location,
		vnet,
		subnet,
	)
}

func attachNewDiskOnVM(
	t *testing.T,
	f fixture.F,
	vmname string,
	diskname string,
	diskSizeGB int,
	sku string,
	caching string,
) {
	f.Shell.Run(
		"./testdata/attach_new_disk.sh",
		f.ResGroupName,
		vmname,
		diskname,
		strconv.Itoa(diskSizeGB),
		sku,
		caching,
	)
}
